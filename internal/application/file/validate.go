package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/idgen"
	"errors"
	"fmt"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/sirupsen/logrus"
)

func (a *File) checkFilesID(sessionID string) (string, []entity.File, *envelope.AppError) {
	files, err := a.postgresql.FileGet.FilesBySessionNotZip(sessionID)
	if err != nil {
		logrus.Error("failed to files by session")
		return "", nil, a.InternalServerError()
	}

	if len(files) == 0 {
		newID, err := idgen.GenerateID()
		if err != nil {
			logrus.Error(err)
			return "", nil, a.InternalServerError()
		}
		return newID, nil, nil
	}

	var filesBase64 []entity.File
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", sessionID, file.Name)
		out, err := a.minioStorage.Get.ByFilename(path)
		if err != nil {
			logrus.Errorf("failed to %s from Minio", path)
			return "", nil, a.InternalServerError()
		}

		if err = domain.CheckFiles(out, file, &filesBase64, path); err != nil {
			logrus.Error(err)
			return "", nil, a.InternalServerError()
		}
	}

	return files[0].FileID, filesBase64, nil
}

func (a *File) registerDownload(countDownload int, session string) *envelope.AppError {
	c := countDownload - 1
	if err := a.postgresql.FileUpdate.CountDownload(c, session); err != nil {
		return a.InternalServerError()
	}
	return nil
}

func (a *File) validDownloadFile(data []byte, file *entity.File) bool {
	_, err := a.downloadFile(data, *file)
	if errors.Is(err, domain.ErrFileDuplicate) {
		file.Name = fmt.Sprintf("dublicate-%s-%s", file.Prefix, file.Name)
		_, err = a.downloadFile(data, *file)
		if err != nil {
			return false
		}
	}
	if err != nil {
		logrus.Error(err)
		return false
	}
	return true
}
