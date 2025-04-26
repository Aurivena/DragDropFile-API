package action

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/domain"
	"encoding/base64"
	"fmt"
	"github.com/Aurivena/answer"
	"github.com/sirupsen/logrus"
	"io"
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

func (a *Action) GetFile(id string, input *model.FileGetInput) (*model.GetFileOutput, answer.ErrorCode) {
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

	f := model.FileGet{
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

	out, err := a.domains.Minio.GetByFilename(filename)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) Create(input *model.FileSaveInput, sessionID string) (*model.FilSaveOutput, answer.ErrorCode) {
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
		for _, val := range input.FileBase64 {
			filesBase64 = append(filesBase64, val)
		}

		out, err := a.save(id, sessionID, filename, filesBase64)
		if err != nil {
			return nil, answer.InternalServerError
		}
		return out, answer.OK
	}

	out, err := a.save(id, sessionID, filename, input.FileBase64)
	if err != nil {
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *Action) save(id, sessionID, filename string, fileBase64 []string) (*model.FilSaveOutput, error) {
	var input model.FileSave

	for i, val := range fileBase64 {
		d, fn, ext, err := domain.DecodeFile(val)
		if err != nil {
			return nil, err
		}
		name := fmt.Sprintf("%s%-d%s", fn, i+1, ext)
		fileID := fmt.Sprintf("%s-%d", id, i+1)
		_, err = a.domains.Minio.DownloadMinio(d, name)
		if err != nil {
			return nil, err
		}

		input = model.FileSave{
			Id:         fileID,
			Name:       name,
			SessionID:  sessionID,
			DataBase64: domain.GetMimeType(val),
		}

		err = a.domains.File.Create(input)
		if err != nil {
			return nil, err
		}
	}

	data, err := a.domains.File.ZipFiles(fileBase64, id)
	if err != nil {
		return nil, err
	}

	m, err := a.domains.Minio.DownloadMinio(data, filename)
	if err != nil {
		return nil, err
	}

	input = model.FileSave{
		Id:         id,
		Name:       filename,
		SessionID:  sessionID,
		DataBase64: ".zip",
	}

	err = a.domains.File.Create(input)
	if err != nil {
		return nil, err
	}

	out := model.FilSaveOutput{
		ID:   id,
		Size: m.Size,

		Count: len(fileBase64),
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

func (a *Action) downloadFile(id, sessionID string, fileBase64 []string) error {
	for i, val := range fileBase64 {
		d, fn, ext, err := domain.DecodeFile(val)

		name := fmt.Sprintf("%s-%d%s", fn, i+1, ext)

		_, err = a.domains.Minio.DownloadMinio(d, name)
		if err != nil {
			return err
		}

		v := model.FileSave{
			Id:         id,
			Name:       name,
			SessionID:  sessionID,
			DataBase64: domain.GetMimeType(val),
		}

		err = a.domains.File.Create(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Action) checkFilesID(sessionID string) ([]string, error) {
	var filesBase64 []string

	filesID, err := a.domains.File.GetIdFileBySession(sessionID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	if filesID == nil {
		return nil, nil
	}

	for _, val := range filesID {
		filename, err := a.domains.File.GetNameByID(val)
		if err != nil {
			return nil, err
		}

		dataBase64, err := a.domains.File.GetDataBase64ByID(val)
		if err != nil {
			return nil, err
		}

		out, err := a.domains.Minio.GetByFilename(filename)
		if err != nil {
			return nil, err
		}

		content, err := io.ReadAll(out.File)
		if err != nil {
			return nil, err
		}

		encoded := base64.StdEncoding.EncodeToString(content)

		fileBase64 := fmt.Sprintf("data:%s;base64,%s", dataBase64, encoded)

		filesBase64 = append(filesBase64, fileBase64)

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
