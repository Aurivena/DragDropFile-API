package action

import (
	"DragDrop-Files/models"
	"DragDrop-Files/pkg/domain"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Aurivena/answer"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"time"
)

const (
	Gone          = 410
	prefixZipFile = "dg-"
)

var (
	ErrorFileDeleted = errors.New("file deleted")
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
	if err := a.domains.UpdateDateDeleted(countDayToDeleted, files.FileID); err != nil {
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

	f := models.FileGet{
		ID:       zipFileID,
		Password: password,
	}

	err = a.domains.ValidatePassword(&f)
	if err != nil {
		logrus.Error(err)
		return nil, answer.Unauthorized
	}

	err = a.domains.ValidateCountDownload(id)
	if err != nil {
		logrus.Error(err)
		if errors.As(err, &ErrorFileDeleted) {
			return nil, Gone
		}
		return nil, answer.InternalServerError
	}

	err = a.domains.ValidateDateDeleted(id)
	if err != nil {
		logrus.Error(err)
		if errors.As(err, &ErrorFileDeleted) {
			return nil, Gone
		}
		return nil, answer.InternalServerError
	}

	path := fmt.Sprintf("%s/%s", file.Session, file.Name)
	out, err := a.domains.Minio.GetByFilename(path)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) SaveFiles(ctx context.Context, sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*models.FileSaveOutput, answer.ErrorCode) {
	if sessionID == "" || len(files) == 0 || len(files) != len(headers) {
		return nil, answer.BadRequest
	}

	var allData []models.File
	for i, file := range files {
		fileData, err := getFileData(file, headers[i])
		if err != nil {
			logrus.WithError(err).Errorf("failed to process file %d", i)
			return nil, answer.BadRequest
		}
		allData = append(allData, *fileData)
	}

	id, existingFiles, err := a.checkFilesID(sessionID)
	if err != nil {
		logrus.WithError(err).Error("failed to check files ID")
		return nil, answer.InternalServerError
	}

	// Фильтруем дубликаты по имени файла
	newFiles := filterDuplicates(allData, existingFiles)
	if len(newFiles) == 0 {
		return nil, answer.BadRequest
	}

	id, err = a.setFileID(id)
	if err != nil {
		logrus.WithError(err).Error("failed to set file ID")
		return nil, answer.InternalServerError
	}

	filename := fmt.Sprintf("%s%s", id, ".zip")
	combinedFiles := append(existingFiles, newFiles...)

	out, err := a.save(ctx, id, sessionID, filename, combinedFiles)
	if err != nil {
		logrus.WithError(err).Error("failed to save files")
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) setFileID(id string) (string, error) {
	if id != "" {
		return id, nil
	}

	newID, err := domain.GenerateID()
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return newID, nil
}

func (a *Action) save(ctx context.Context, id, sessionID, filename string, files []models.File) (*models.FileSaveOutput, error) {
	var input models.FileSave
	updatedFiles := make([]models.File, 0, len(files))

	for _, file := range files {
		data, err := domain.DecodeFile(file.FileBase64)
		if err != nil {
			logrus.WithError(err).Error("failed to decode file")
			return nil, err
		}

		// Используем наносекунды для уникальности имени
		uniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)

		_, err = a.domains.Minio.DownloadMinio(data, sessionID, uniqueName)
		if err != nil {
			logrus.WithError(err).Error("failed to upload to Minio")
			return nil, err
		}

		input = models.FileSave{
			Id:        id,
			Name:      uniqueName,
			SessionID: sessionID,
			MimeType:  domain.GetMimeType(file.FileBase64),
		}

		if err := a.domains.File.Create(ctx, input); err != nil {
			logrus.WithError(err).Error("failed to save file metadata")
			return nil, err
		}

		file.Filename = uniqueName
		updatedFiles = append(updatedFiles, file)
	}

	fileIDZip := fmt.Sprintf("%s%s", prefixZipFile, id)
	zipData, err := a.domains.File.ZipFiles(updatedFiles, fileIDZip)
	if err != nil {
		logrus.WithError(err).Error("failed to zip files")
		return nil, err
	}

	zipUniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	meta, err := a.domains.Minio.DownloadMinio(zipData, sessionID, zipUniqueName)
	if err != nil {
		logrus.WithError(err).Error("failed to upload ZIP to Minio")
		return nil, err
	}

	input = models.FileSave{
		Id:        fileIDZip,
		Name:      zipUniqueName,
		SessionID: sessionID,
		MimeType:  "zip",
	}

	if err := a.domains.File.Create(ctx, input); err != nil {
		logrus.WithError(err).Error("failed to save ZIP metadata")
		return nil, err
	}

	return &models.FileSaveOutput{
		ID:    id,
		Size:  meta.Size,
		Count: len(files),
	}, nil
}

func (a *Action) GetDataFile(id string) (*models.DataOutput, answer.ErrorCode) {
	id = fmt.Sprintf("%s%s", prefixZipFile, id)
	out, err := a.domains.GetDataFile(id)
	if err != nil {
		if errors.As(err, &ErrorFileDeleted) {
			return nil, Gone
		}
		return nil, answer.InternalServerError
	}
	return out, answer.OK
}

func (a *Action) checkFilesID(sessionID string) (string, []models.File, error) {
	files, err := a.domains.File.GetFilesBySession(sessionID)
	if err != nil {
		logrus.WithError(err).Error("failed to get files by session")
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
			logrus.WithError(err).Errorf("failed to get file %s from Minio", path)
			return "", nil, err
		}

		content, err := io.ReadAll(out.File)
		if err != nil {
			logrus.WithError(err).Errorf("failed to read file %s", path)
			_ = out.File.Close() // Закрываем файл
			return "", nil, err
		}
		_ = out.File.Close() // Закрываем файл

		encoded := base64.StdEncoding.EncodeToString(content)
		fileBase64 := fmt.Sprintf("data:%s;base64,%s", val.MimeType, encoded)

		file := models.File{
			FileBase64: fileBase64,
			Filename:   val.Name,
		}
		filesBase64 = append(filesBase64, file)

		if err := a.domains.Minio.Delete(val.Name); err != nil {
			logrus.WithError(err).Errorf("failed to delete file %s from Minio", val.Name)
			return "", nil, err
		}
	}

	if err := a.domains.DeleteFilesBySessionID(sessionID); err != nil {
		logrus.WithError(err).Error("failed to delete files by session ID")
		return "", nil, err
	}

	return files[0].FileID, filesBase64, nil
}
func getFileData(file multipart.File, header *multipart.FileHeader) (*models.File, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.WithError(err).Error("failed to read file")
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

func filterDuplicates(newFiles, existingFiles []models.File) []models.File {
	existing := make(map[string]struct{}, len(existingFiles))
	for _, f := range existingFiles {
		existing[f.Filename] = struct{}{}
	}

	var uniqueFiles []models.File
	for _, f := range newFiles {
		if _, exists := existing[f.Filename]; !exists {
			uniqueFiles = append(uniqueFiles, f)
			existing[f.Filename] = struct{}{}
		}
	}
	return uniqueFiles
}
