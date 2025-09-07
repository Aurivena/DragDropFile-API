package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/idgen"
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (a *File) checkFilesID(sessionID string) (string, []entity.File, error) {
	files, err := a.reader.FilesBySessionNotZip(sessionID)
	if err != nil {
		logrus.Error("failed to files by session")
		return "", nil, domain.InternalError
	}

	if len(files) == 0 {
		newID, err := idgen.GenerateID()
		if err != nil {
			logrus.Error(err)
			return "", nil, domain.InternalError
		}
		return newID, nil, nil
	}

	var filesBase64 []entity.File
	for _, file := range files {
		path := fmt.Sprintf("%s/%s", sessionID, file.Name)
		out, err := a.minioStorage.Reader.ByFilename(path)
		if err != nil {
			logrus.Errorf("failed to %s from Minio", path)
			return "", nil, domain.InternalError
		}

		if err = domain.CheckFiles(out, file, &filesBase64, path); err != nil {
			logrus.Error(err)
			return "", nil, domain.InternalError
		}
	}

	return files[0].FileID, filesBase64, nil
}

func (a *File) registerDownload(countDownload int, session string) error {
	c := countDownload - 1
	if err := a.updater.CountDownload(c, session); err != nil {
		return domain.InternalError
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
