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

func (a *Action) Create(ctx context.Context, sessionID string, files []multipart.File, headers []*multipart.FileHeader) (*models.FilSaveOutput, answer.ErrorCode) {
	var allData []models.File
	for i := range files {
		fileData, err := getFileData(files[i], headers[i])
		if err != nil {
			return nil, answer.InternalServerError
		}
		allData = append(allData, *fileData)
	}

	id, filesBase64, err := a.checkFilesID(sessionID)
	if err != nil {
		return nil, answer.InternalServerError
	}

	id, err = a.setfileID(id)
	if err != nil {
		return nil, answer.InternalServerError
	}

	filename := fmt.Sprintf("%s%s", id, ".zip")

	filesBase64 = append(filesBase64, allData...)

	out, err := a.save(ctx, id, sessionID, filename, filesBase64)
	if err != nil {
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) setfileID(id string) (string, error) {
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

func (a *Action) save(ctx context.Context, id, sessionID, filename string, files []models.File) (*models.FilSaveOutput, error) {
	var input models.FileSave

	updatedFiles := make([]models.File, 0, len(files))

	for _, val := range files {
		d, err := domain.DecodeFile(val.FileBase64)
		if err != nil {
			return nil, err
		}

		uniqueName := fmt.Sprintf("%d_%s", time.Now().Second(), val.Filename)

		_, err = a.domains.Minio.DownloadMinio(d, sessionID, uniqueName)
		if err != nil {
			return nil, err
		}

		input = models.FileSave{
			Id:        id,
			Name:      uniqueName,
			SessionID: sessionID,
			MimeType:  domain.GetMimeType(val.FileBase64),
		}

		err = a.domains.File.Create(ctx, input)
		if err != nil {
			return nil, err
		}

		val.Filename = uniqueName
		updatedFiles = append(updatedFiles, val)
	}
	fileIDZip := fmt.Sprintf("%s%s", prefixZipFile, id)
	data, err := a.domains.File.ZipFiles(updatedFiles, fileIDZip)
	if err != nil {
		return nil, err
	}

	zipUniqueName := fmt.Sprintf("%d_%s", time.Now().Hour(), filename)
	m, err := a.domains.Minio.DownloadMinio(data, sessionID, zipUniqueName)
	if err != nil {
		return nil, err
	}

	input = models.FileSave{
		Id:        fileIDZip,
		Name:      zipUniqueName,
		SessionID: sessionID,
		MimeType:  ".zip",
	}

	err = a.domains.File.Create(ctx, input)
	if err != nil {
		return nil, err
	}

	out := models.FilSaveOutput{
		ID:    id,
		Size:  m.Size,
		Count: len(files),
	}
	return &out, nil
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
	var filesBase64 []models.File
	var file models.File

	files, err := a.domains.File.GetFilesBySession(sessionID)
	if err != nil {
		logrus.Error(err)
		return "", nil, err
	}

	if files == nil {
		return "", nil, nil
	}

	for _, val := range files {
		path := fmt.Sprintf("%s/%s", sessionID, val.Name)
		out, err := a.domains.Minio.GetByFilename(path)
		if err != nil {
			return "", nil, err
		}

		content, err := io.ReadAll(out.File)
		if err != nil {
			return "", nil, err
		}

		encoded := base64.StdEncoding.EncodeToString(content)

		fileBase64 := fmt.Sprintf("data:%s;base64,%s", val.MimeType, encoded)

		file.FileBase64 = fileBase64
		file.Filename = val.Name

		filesBase64 = append(filesBase64, file)

		err = a.domains.Minio.Delete(val.Name)
		if err != nil {
			return "", nil, err
		}
	}

	if err = a.domains.DeleteFilesBySessionID(sessionID); err != nil {
		return "", nil, err
	}

	return files[0].FileID, filesBase64, nil
}

func getFileData(file multipart.File, header *multipart.FileHeader) (*models.File, error) {
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	encoded := base64.StdEncoding.EncodeToString(fileBytes)
	mimeType := header.Header.Get("Content-Type")
	fileBase64 := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

	fileData := models.File{
		FileBase64: fileBase64,
		Filename:   header.Filename,
	}
	return &fileData, nil
}
