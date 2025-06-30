package action

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/models"
	"DragDrop-Files/pkg"
	"context"
	"errors"
	"fmt"
	"github.com/Aurivena/answer"
	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"sync"
)

const (
	Gone             = 410
	prefixZipFile    = "dg-"
	lenCodeForID     = 12
	lenCodeForPrefix = 3
)

var (
	ErrFileDeleted   = errors.New("file deleted")
	ErrDuplicateFile = errors.New("file duplicate")
)

func (a *Action) UpdateCountDownload(count int, sessionID string) answer.ErrorCode {
	if err := a.domains.UpdateCountDownload(count, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
func (a *Action) UpdateDateDeleted(countDayToDeleted int, sessionID string) answer.ErrorCode {
	files, err := a.domains.File.GetZipMetaBySession(sessionID)
	if err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}
	if err = a.domains.UpdateDateDeleted(countDayToDeleted, files.FileID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
func (a *Action) UpdatePassword(password, sessionID string) answer.ErrorCode {
	if err := a.domains.UpdatePassword(password, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}

func (a *Action) UpdateDescription(description, sessionID string) answer.ErrorCode {
	if err := a.domains.UpdateDescription(description, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}

func (a *Action) GetFile(id, password string) (*models.GetFileOutput, answer.ErrorCode) {
	zipFileID := fmt.Sprintf("%s%s", prefixZipFile, id)
	file, err := a.domains.File.GetByID(zipFileID)
	if err != nil {
		logrus.Error(err)
		return nil, answer.BadRequest
	}

	errCode := a.validate(zipFileID, id, password)
	if errCode != answer.NoContent {
		return nil, errCode
	}

	path := fmt.Sprintf("%s/%s", file.Session, file.Name)
	out, err := a.domains.Minio.GetByFilename(path)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) validate(zipFileID, id, password string) answer.ErrorCode {

	f := models.FileGet{
		ID:       zipFileID,
		Password: password,
	}

	err := a.domains.ValidatePassword(&f)
	if err != nil {
		logrus.Error(err)
		return answer.Unauthorized
	}

	err = a.domains.ValidateCountDownload(id)
	if err != nil {
		logrus.Error(err)
		return a.handleIfDeleted(id, err)
	}

	err = a.domains.ValidateDateDeleted(id)
	if err != nil {
		logrus.Error(err)
		return a.handleIfDeleted(id, err)
	}
	return answer.NoContent
}

func (a *Action) SaveFiles(ctx context.Context, sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*models.FileSaveOutput, answer.ErrorCode) {
	if sessionID == "" || len(files) == 0 || len(files) != len(headers) {
		return nil, answer.BadRequest
	}

	newFiles := pkg.InfoNewFile(files, headers)

	id, existingFiles, err := a.checkFilesID(sessionID)
	if err != nil {
		logrus.Error("failed to check files ID")
		return nil, answer.InternalServerError
	}

	id, err = a.setFileID(id)
	if err != nil {
		logrus.Error("failed to set file ID")
		return nil, answer.InternalServerError
	}

	out, err := a.saveFilesToStorage(ctx, id, sessionID, newFiles, existingFiles)
	if err != nil {
		logrus.Error("failed to save files")
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) setFileID(id string) (string, error) {
	if id != "" {
		return id, nil
	}

	newID, err := domain.GenerateID(lenCodeForID)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newID, nil
}

func (a *Action) GetDataFile(id string) (*models.DataOutput, answer.ErrorCode) {
	id = fmt.Sprintf("%s%s", prefixZipFile, id)
	out, err := a.domains.GetDataFile(id)
	if err != nil {
		if errors.Is(err, ErrDuplicateFile) {
			if err = a.domains.Minio.Delete(id); err != nil {
				return nil, answer.InternalServerError
			}
			return nil, Gone
		}
		return nil, answer.InternalServerError
	}
	return out, answer.OK
}

func (a *Action) checkValidDownloadFile(data []byte, f *models.File, sessionID, id, prefix string, ctx context.Context) error {
	_, err := a.downloadFile(data, pkg.GetMimeType(f.FileBase64), f.Filename, sessionID, id, ctx)
	if errors.Is(err, ErrDuplicateFile) {
		f.Filename = fmt.Sprintf("dublicate-%s-%s", prefix, f.Filename)
		_, err = a.downloadFile(data, pkg.GetMimeType(f.FileBase64), f.Filename, sessionID, id, ctx)
		if err != nil {
			logrus.Error(err)
			return err
		}
	}

	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}

func (a *Action) downloadZipFile(id, sessionID string, files []models.File, ctx context.Context) (*minio.UploadInfo, error) {
	fileIDZip := fmt.Sprintf("%s%s", prefixZipFile, id)
	zipData, err := pkg.ZipFiles(files, fileIDZip)
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
		logrus.Error(err)
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

func (a *Action) saveFilesToStorage(ctx context.Context, id, sessionID string, newFiles, oldFiles []models.File) (*models.FileSaveOutput, error) {
	var (
		wg             sync.WaitGroup
		mu             sync.Mutex
		processedFiles []models.File
	)

	prefix, err := domain.GenerateID(lenCodeForPrefix)
	if err != nil {
		logrus.Errorf("failed to generate prefix: %v", err)
		return nil, fmt.Errorf("failed to generate prefix: %w", err)
	}

	for _, file := range newFiles {
		wg.Add(1)
		go func(f models.File) {
			defer wg.Done()

			data, err := pkg.DecodeFile(f.FileBase64)
			if err != nil {
				return
			}

			if err = a.checkValidDownloadFile(data, &f, sessionID, id, prefix, ctx); err != nil {
				return
			}

			mu.Lock()
			processedFiles = append(processedFiles, f)
			mu.Unlock()
		}(file)
	}
	wg.Wait()

	processedFiles = append(processedFiles, oldFiles...)

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
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", sessionID, file.Name)
		out, err := a.domains.Minio.GetByFilename(path)
		if err != nil {
			logrus.Errorf("failed to get file %s from Minio", path)
			return "", nil, err
		}

		if err = pkg.CheckFiles(out, file, &filesBase64, path); err != nil {
			logrus.Error(err)
			return "", nil, err
		}

		if err = a.domains.Minio.Delete(file.Name); err != nil {
			logrus.Errorf("failed to delete file %s from Minio", file.Name)
			return "", nil, err
		}
	}

	if err = a.domains.DeleteFilesBySessionID(sessionID); err != nil {
		logrus.Error("failed to delete files by session ID")
		return "", nil, err
	}

	return files[0].FileID, filesBase64, nil
}

func (a *Action) handleIfDeleted(id string, err error) answer.ErrorCode {
	if errors.Is(err, ErrFileDeleted) {
		if err = a.domains.Minio.Delete(id); err != nil {
			logrus.Error(err)
			return answer.InternalServerError
		}
		return Gone
	}
	return answer.InternalServerError
}
