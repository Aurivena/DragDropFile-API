package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/archive"
	"DragDrop-Files/pkg/idgen"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (a *File) downloadZipFile(id, sessionID string, files []entity.File) (*entity.FileSaveOutput, error) {
	fileIDZip := fmt.Sprintf("%s%s", domain.PrefixZipFile, id)
	zipData, err := archive.ZipFiles(files, fileIDZip)
	if err != nil {
		logrus.Error("failed to zip files")
		return nil, domain.InternalError
	}

	generatedID, err := idgen.GenerateID()
	if err != nil {
		logrus.Error("failed to generate id")
		return nil, domain.InternalError
	}

	if err = a.deleted.FileID(fileIDZip); err != nil {
		logrus.Error("failed to delete id file")
		return nil, domain.InternalError
	}

	zipFile := entity.File{
		FileID:    fileIDZip,
		Name:      fmt.Sprintf("%s%s", generatedID, domain.MimeTypeZip),
		SessionID: sessionID,
		MimeType:  domain.MimeTypeZip,
	}

	meta, err := a.downloadFile(zipData, zipFile)
	if err != nil {
		return nil, domain.InternalError
	}
	out := entity.FileSaveOutput{
		ID:    id,
		Size:  meta.Size,
		Count: len(files),
	}
	return &out, nil
}

func (a *File) downloadFile(data []byte, file entity.File) (*minio.UploadInfo, error) {
	id, err := a.writer.Execute(file)
	if err != nil {
		logrus.Error("failed to save metadata")
		return nil, err
	}

	file.ID = id

	if err = a.writer.ExecuteSession(file); err != nil {
		logrus.Error("failed to save session")
		return nil, err
	}

	week := time.Now().Add(time.Duration(24*time.Hour) * 7)
	file.TimeDeleted = &week
	if err = a.writer.ExecuteParameters(file); err != nil {
		logrus.Error("failed to save parameters")
		return nil, err
	}

	meta, err := a.minioStorage.Save.File(data, fmt.Sprintf("%s/%s", file.SessionID, file.Name))
	if err != nil {
		return nil, err
	}

	return meta, nil
}
