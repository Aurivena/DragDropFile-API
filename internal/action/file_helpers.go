package action

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/models"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"sync"
)

func (a *Action) checkValidDownloadFile(data []byte, f *models.File, sessionID, id, prefix string, ctx context.Context) error {
	_, err := a.downloadFile(data, f.FileBase64, f.Filename, sessionID, id, ctx)
	if err != nil {
		if errors.As(err, &ErrorDuplicateFile) {
			f.Filename = fmt.Sprintf("%s-%s", prefix, f.Filename)
			_, err := a.downloadFile(data, f.FileBase64, f.Filename, sessionID, id, ctx)
			if err != nil {
				logrus.Error(err)
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func (a *Action) downloadZipFile(id, sessionID string, files []models.File, ctx context.Context) (*minio.UploadInfo, error) {
	fileIDZip := fmt.Sprintf("%s%s", prefixZipFile, id)
	zipData, err := a.domains.File.ZipFiles(files, fileIDZip)
	if err != nil {
		logrus.Error("failed to zip files")
		return nil, err
	}

	zipUniqueName := fmt.Sprintf("%s.zip", uuid.NewString())
	meta, err := a.downloadFile(zipData, ".zip", zipUniqueName, sessionID, fileIDZip, ctx)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (a *Action) downloadFile(data []byte, mimeType, filename, sessionID, id string, ctx context.Context) (*minio.UploadInfo, error) {
	meta, err := a.domains.Minio.DownloadMinio(data, sessionID, filename)
	if err != nil {
		logrus.Error("failed to upload to Minio")
		return nil, err
	}

	input := models.FileSave{
		Id:        id,
		Name:      filename,
		SessionID: sessionID,
		MimeType:  mimeType,
	}

	if err = a.domains.File.Create(ctx, input); err != nil {
		logrus.Error("failed to save file metadata")
		return nil, err
	}

	return meta, nil
}

func getFileData(file multipart.File, header *multipart.FileHeader) (*models.File, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Error("failed to read file")
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	mimeType := header.Header.Get("Content-Type")
	fileBase64 := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	return &models.File{
		FileBase64: fileBase64,
		Filename:   header.Filename,
	}, nil
}

func (a *Action) saveFilesToStorage(ctx context.Context, id, sessionID string, files []models.File) (*models.FileSaveOutput, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	var (
		wg             sync.WaitGroup
		mu             sync.Mutex
		processedFiles []models.File
		errors         []error
	)

	prefix, err := domain.GenerateID(lenCodeForPrefix)
	if err != nil {
		logrus.Errorf("failed to generate prefix: %v", err)
		return nil, fmt.Errorf("failed to generate prefix: %w", err)
	}

	for _, file := range files {
		wg.Add(1)
		go func(f models.File) {
			defer wg.Done()

			data, err := domain.DecodeFile(f.FileBase64)
			if err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to decode file %s: %w", f.Filename, err))
				mu.Unlock()
				return
			}

			if err = a.checkValidDownloadFile(data, &f, sessionID, id, prefix, ctx); err != nil {
				mu.Lock()
				errors = append(errors, fmt.Errorf("failed to validate file %s: %w", f.Filename, err))
				mu.Unlock()
				return
			}

			mu.Lock()
			processedFiles = append(processedFiles, f)
			mu.Unlock()
		}(file)
	}
	wg.Wait()

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to process %d file(s): %v", len(errors), errors)
	}

	if len(processedFiles) == 0 {
		return nil, fmt.Errorf("no files were processed successfully")
	}

	meta, err := a.downloadZipFile(id, sessionID, processedFiles, ctx)
	if err != nil {
		logrus.Errorf("failed to create zip file: %v", err)
		return nil, fmt.Errorf("failed to create zip file: %w", err)
	}

	return &models.FileSaveOutput{
		ID:    id,
		Size:  meta.Size,
		Count: len(processedFiles),
	}, nil
}
func (a *Action) checkFilesID(sessionID string) (string, []models.File, error) {
	files, err := a.domains.File.GetFilesBySession(sessionID)
	if err != nil {
		logrus.Error("failed to get files by session")
		return "", nil, err
	}

	if len(files) == 0 {
		return "", nil, nil
	}

	var filesBase64 []models.File
	for _, val := range files {
		path := fmt.Sprintf("%s/%s", sessionID, val.Name)
		out, err := a.domains.Minio.GetByFilename(path)
		if err != nil {
			logrus.Errorf("failed to get file %s from Minio", path)
			return "", nil, err
		}

		content, err := io.ReadAll(out.File)
		if err != nil {
			logrus.Errorf("failed to read file %s", path)

			_ = out.File.Close()
			return "", nil, err
		}
		_ = out.File.Close()

		encoded := base64.StdEncoding.EncodeToString(content)
		fileBase64 := fmt.Sprintf("data:%s;base64,%s", val.MimeType, encoded)

		file := models.File{
			FileBase64: fileBase64,
			Filename:   val.Name,
		}
		filesBase64 = append(filesBase64, file)

		if err = a.domains.Minio.Delete(val.Name); err != nil {
			logrus.Errorf("failed to delete file %s from Minio", val.Name)
			return "", nil, err
		}
	}

	if err = a.domains.DeleteFilesBySessionID(sessionID); err != nil {
		logrus.Error("failed to delete files by session ID")
		return "", nil, err
	}

	return files[0].FileID, filesBase64, nil
}
