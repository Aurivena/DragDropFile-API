package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"errors"
	"fmt"

	_ "github.com/Aurivena/spond/v2/core"
	"github.com/sirupsen/logrus"
)

func (a *File) Get(id, password string) (*entity.GetFileOutput, error) {
	zipFileID := fmt.Sprintf("%s%s", domain.PrefixZipFile, id)
	zipFile, err := a.reader.ByID(zipFileID)
	if err != nil {
		logrus.Error(err)
		return nil, domain.NotFoundError
	}

	file, err := a.reader.ByID(id)
	if err != nil {
		return nil, domain.NotFoundError
	}

	if err = domain.ValidateFile(password, file); err != nil {
		if errors.Is(err, domain.ErrFileDeleted) {
			if err = a.minioStorage.Delete.ByFilename(id); err != nil {
				return nil, domain.NotFoundError
			}
			return nil, domain.GoneError
		}
		if errors.Is(err, domain.PasswordInvalidError) {
			return nil, domain.PasswordInvalidError
		}
		return nil, domain.InternalError
	}

	path := fmt.Sprintf("%s/%s", zipFile.SessionID, zipFile.Name)
	out, err := a.minioStorage.Get.ByFilename(path)
	if err != nil {
		logrus.Error(err)
		return nil, domain.NotFoundError
	}

	return out, nil
}

func (a *File) Register(fileID string) error {
	file, err := a.reader.ByID(fileID)
	if err != nil {
		return domain.NotFoundError
	}

	if errResp := a.registerDownload(file.CountDownload, file.SessionID); errResp != nil {
		return errResp
	}

	return nil
}

func (a *File) Data(id string) (*entity.FileData, error) {
	out, err := a.reader.DataFile(id)
	if err != nil {
		return nil, domain.NotFoundError
	}
	return out, nil
}
