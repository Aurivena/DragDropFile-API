package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/archive"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (a *File) downloadZipFile(id, sessionID, prefixZipFile string, files []entity.File) (*minio.UploadInfo, error) {
	fileIDZip := fmt.Sprintf("%s%s", prefixZipFile, id)
	zipData, err := archive.ZipFiles(files, fileIDZip)
	if err != nil {
		logrus.Error("failed to zip files")
		return nil, err
	}

	uuidV7, err := uuid.NewV7()
	if err != nil {
		logrus.Error("failed to generate uid")
		return nil, err
	}

	zipUniqueName := fmt.Sprintf("%s%s", uuidV7.String(), domain.MimeTypeZip)
	zipFile := entity.File{
		FileID:    id,
		Name:      zipUniqueName,
		SessionID: sessionID,
		MimeType:  domain.MimeTypeZip,
	}

	meta, err := a.downloadFile(zipData, zipFile)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (a *File) downloadFile(data []byte, file entity.File) (*minio.UploadInfo, error) {
	currentTime := time.Now().Format(time.RFC3339)

	id, err := a.postgresql.FileSave.Execute(file)
	if err != nil {
		logrus.Error("failed to save metadata")
		return nil, err
	}

	file.ID = id

	if err = a.postgresql.FileSave.ExecuteSession(file.SessionID, id); err != nil {
		logrus.Error("failed to save session")
		return nil, err
	}

	if err = a.postgresql.FileSave.ExecuteParameters(file, currentTime); err != nil {
		logrus.Error("failed to save parameters")
		return nil, err
	}

	meta, err := a.minioStorage.Save.File(data, file.SessionID, file.Name)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return meta, nil
}
