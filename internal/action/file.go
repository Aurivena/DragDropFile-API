package action

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/models"
	"context"
	"errors"
	"fmt"
	"github.com/Aurivena/answer"
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
	ErrorFileDeleted   = errors.New("file deleted")
	ErrorDuplicateFile = errors.New("file duplicate")
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
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}

	var newFiles []models.File
	for i, file := range files {
		wg.Add(1)
		go func(f multipart.File, headers []*multipart.FileHeader, index int) {
			defer wg.Done()
			fileData, err := getFileData(f, headers[index])
			if err != nil {
				logrus.Errorf("failed to process file %d", i)
				return
			}
			mu.Lock()
			newFiles = append(newFiles, *fileData)
			mu.Unlock()
		}(file, headers, i)
	}

	wg.Wait()

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
		if errors.As(err, &ErrorFileDeleted) {
			return nil, Gone
		}
		return nil, answer.InternalServerError
	}
	return out, answer.OK
}
