package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"errors"
	"fmt"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/sirupsen/logrus"
)

func (a *File) checkFilesID(sessionID string) (string, []entity.File, error) {
	files, err := a.postgresql.FileGet.FilesBySessionNotZip(sessionID)
	if err != nil {
		logrus.Error("failed to files by session")
		return "", nil, err
	}

	if len(files) == 0 {
		return "", nil, nil
	}

	var filesBase64 []entity.File
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", sessionID, file.Name)
		out, err := a.minioStorage.Get.ByFilename(path)
		if err != nil {
			logrus.Errorf("failed to g g %s from Minio", path)
			return "", nil, err
		}

		if err = domain.CheckFiles(out, file, &filesBase64, path); err != nil {
			logrus.Error(err)
			return "", nil, err
		}

		if err = a.minioStorage.Delete.File(file.Name); err != nil {
			logrus.Errorf("failed to delete g %s from Minio", file.Name)
			return "", nil, err
		}
	}

	if err = a.postgresql.FileDelete.FilesBySessionID(sessionID); err != nil {
		logrus.Error("failed to delete files by session ID")
		return "", nil, err
	}

	return files[0].FileID, filesBase64, nil
}

func (a *File) registerDownload(countDownload int, session string) *envelope.AppError {
	c := countDownload - 1
	err := a.postgresql.FileUpdate.CountDownload(c, session)
	if err != nil {
		return a.InternalServerError()
	}
	return nil
}

func (a *File) validDownloadFile(data []byte, file entity.File, prefix string) bool {
	_, err := a.downloadFile(data, file)
	if errors.Is(err, domain.ErrFileDeleted) {
		file.Name = fmt.Sprintf("dublicate-%s-%s", prefix, file.Name)
		_, err = a.downloadFile(data, file)
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
