package action

import (
	"DragDrop-Files/models"
	"DragDrop-Files/pkg/domain"
	"encoding/base64"
	"fmt"
	"github.com/Aurivena/answer"
	"github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
)

func (a *Action) UpdateCountDownload(count int, sessionID string) answer.ErrorCode {
	if err := a.domains.UpdateCountDownload(count, sessionID); err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}
func (a *Action) UpdateDateDeleted(countDayToDeleted int, sessionID string) answer.ErrorCode {
	if err := a.domains.UpdateDateDeleted(countDayToDeleted, sessionID); err != nil {
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

func (a *Action) GetFile(id string, input *models.FileGetInput) (*models.GetFileOutput, answer.ErrorCode) {
	filename, err := a.domains.File.GetNameByID(id)
	if err != nil {
		logrus.Error(err)
		return nil, answer.BadRequest
	}

	sessionID, err := a.domains.File.GetSessionByID(id)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	f := models.FileGet{
		SessionID: sessionID,
		Password:  input.Password,
	}

	err = a.domains.ValidatePassword(&f)
	if err != nil {
		logrus.Error(err)
		return nil, answer.BadRequest
	}

	err = a.domains.ValidateCountDownload(sessionID)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	err = a.domains.ValidateDateDeleted(sessionID)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	path := fmt.Sprintf("%s/%s", sessionID, filename)
	out, err := a.domains.Minio.GetByFilename(path)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) Create(sessionID string, file multipart.File, header *multipart.FileHeader) (*models.FilSaveOutput, answer.ErrorCode) {
	fileData, err := getFileData(file, header)
	if err != nil {
		return nil, answer.InternalServerError
	}

	id, err := domain.GenerateID()
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	filename := fmt.Sprintf("%s%s", id, ".zip")

	filesBase64, err := a.checkFilesID(sessionID)
	if err != nil {
		return nil, answer.InternalServerError
	}

	if filesBase64 != nil {
		filesBase64 = append(filesBase64, *fileData)

		out, err := a.save(id, sessionID, filename, filesBase64)
		if err != nil {
			return nil, answer.InternalServerError
		}
		return out, answer.OK
	}

	filesBase64 = append(filesBase64, *fileData)

	out, err := a.save(id, sessionID, filename, filesBase64)
	if err != nil {
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) save(id, sessionID, filename string, files []models.File) (*models.FilSaveOutput, error) {
	var input models.FileSave

	for i, val := range files {
		d, err := domain.DecodeFile(val.FileBase64)
		if err != nil {
			return nil, err
		}
		fileID := fmt.Sprintf("%s%d", id, i+1)
		_, err = a.domains.Minio.DownloadMinio(d, sessionID, val.Filename)
		if err != nil {
			return nil, err
		}

		input = models.FileSave{
			Id:        fileID,
			Name:      val.Filename,
			SessionID: sessionID,
			MimeType:  domain.GetMimeType(val.FileBase64),
		}

		err = a.domains.File.Create(input)
		if err != nil {
			return nil, err
		}
	}

	data, err := a.domains.File.ZipFiles(files, id)
	if err != nil {
		return nil, err
	}

	m, err := a.domains.Minio.DownloadMinio(data, sessionID, filename)
	if err != nil {
		return nil, err
	}

	input = models.FileSave{
		Id:        id,
		Name:      filename,
		SessionID: sessionID,
		MimeType:  ".zip",
	}

	err = a.domains.File.Create(input)
	if err != nil {
		return nil, err
	}

	out := models.FilSaveOutput{
		ID:   id,
		Size: m.Size,

		Count: len(files),
	}
	return &out, nil
}

func (a *Action) DeleteFile(id string) answer.ErrorCode {
	err := a.domains.Minio.Delete(id)
	if err != nil {
		logrus.Error(err)
		return answer.InternalServerError
	}

	return answer.NoContent
}

func (a *Action) downloadFile(id, sessionID string, files []models.File) error {
	for _, val := range files {
		d, err := domain.DecodeFile(val.FileBase64)

		_, err = a.domains.Minio.DownloadMinio(d, sessionID, val.Filename)
		if err != nil {
			return err
		}

		v := models.FileSave{
			Id:        id,
			Name:      val.Filename,
			SessionID: sessionID,
			MimeType:  domain.GetMimeType(val.FileBase64),
		}

		err = a.domains.File.Create(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Action) checkFilesID(sessionID string) ([]models.File, error) {
	var filesBase64 []models.File
	var file models.File

	files, err := a.domains.File.GetFileBySession(sessionID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if files == nil {
		return nil, nil
	}

	for _, val := range files {
		filename, err := a.domains.File.GetNameByID(val.Id)
		if err != nil {
			return nil, err
		}

		mimeType, err := a.domains.File.GetMimeTypeByID(val.Id)
		if err != nil {
			return nil, err
		}

		path := fmt.Sprintf("%s/%s", sessionID, filename)
		out, err := a.domains.Minio.GetByFilename(path)
		if err != nil {
			return nil, err
		}

		content, err := io.ReadAll(out.File)
		if err != nil {
			return nil, err
		}

		encoded := base64.StdEncoding.EncodeToString(content)

		fileBase64 := fmt.Sprintf("data:%s;base64,%s", mimeType, encoded)

		file.FileBase64 = fileBase64
		file.Filename = filename

		filesBase64 = append(filesBase64, file)

		err = a.domains.Minio.Delete(filename)
		if err != nil {
			return nil, err
		}
	}

	if err = a.domains.DeleteFilesBySessionID(sessionID); err != nil {
		return nil, err
	}

	return filesBase64, nil
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
