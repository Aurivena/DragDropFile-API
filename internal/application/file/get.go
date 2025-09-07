package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"errors"
	"fmt"

	_ "github.com/Aurivena/spond/v2/core"
	"github.com/Aurivena/spond/v2/envelope"
	"github.com/sirupsen/logrus"
)

func (a *File) Get(id, password string) (*entity.GetFileOutput, *envelope.AppError) {
	zipFileID := fmt.Sprintf("%s%s", domain.PrefixZipFile, id)
	zipFile, err := a.postgresql.FileGet.ByID(zipFileID)
	if err != nil {
		logrus.Error(err)
		return nil, a.NotFound()
	}

	file, err := a.postgresql.FileGet.ByID(id)
	if err != nil {
		return nil, a.NotFound()
	}

	if err = domain.ValidateFile(password, file); err != nil {
		if errors.Is(err, domain.ErrFileDeleted) {
			if err = a.minioStorage.Delete.ByFilename(id); err != nil {
				return nil, a.NotFound()
			}
			return nil, a.Gone()
		}
		if errors.Is(err, domain.ErrPasswordInvalid) {
			return nil, a.PasswordInvalid()
		}
		return nil, a.InternalServerError()
	}

	if errResp := a.registerDownload(file.CountDownload, file.SessionID); errResp != nil {
		return nil, errResp
	}

	path := fmt.Sprintf("%s/%s", zipFile.SessionID, zipFile.Name)
	out, err := a.minioStorage.Get.ByFilename(path)
	if err != nil {
		logrus.Error(err)
		return nil, a.NotFound()
	}

	return out, nil
}

func (a *File) Data(id string) (*entity.FileData, *envelope.AppError) {
	out, err := a.postgresql.FileGet.DataFile(id)
	if err != nil {
		return nil, a.NotFound()
	}
	return out, nil
}
