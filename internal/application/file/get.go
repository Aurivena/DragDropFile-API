package file

import (
	"DragDrop-Files/internal/domain/entity"
	"errors"
	"fmt"
	"github.com/Aurivena/answer"
	"github.com/sirupsen/logrus"
)

var errFileDeleted = errors.New("file deleted")

func (a *File) GetFile(id, password string) (*entity.GetFileOutput, answer.ErrorCode) {
	zipFileID := fmt.Sprintf("%s%s", prefixZipFile, id)
	file, err := a.repo.Get.ByID(zipFileID)
	if err != nil {
		logrus.Error(err)
		return nil, answer.BadRequest
	}

	if errCode := a.service.Validate.File(zipFileID, id, password); err != nil {
		if errors.Is(errCode, errFileDeleted) {
			return nil, Gone
		}
		return nil, answer.InternalServerError
	}

	path := fmt.Sprintf("%s/%s", file.Session, file.Name)
	out, err := a.minio.Get.ByFilename(path)
	if err != nil {
		logrus.Error(err)
		return nil, answer.InternalServerError
	}

	return out, answer.OK
}

func (a *File) GetDataFile(id string) (*entity.DataOutput, answer.ErrorCode) {
	id = fmt.Sprintf("%s%s", prefixZipFile, id)
	out, err := a.repo.Get.DataFile(id)
	if err != nil {
		if errors.Is(err, ErrDuplicateFile) {
			if err = a.minio.Delete.File(id); err != nil {
				return nil, answer.InternalServerError
			}
			return nil, Gone
		}
		return nil, answer.InternalServerError
	}
	return out, answer.OK
}
