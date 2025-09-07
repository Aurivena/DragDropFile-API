package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/archive"
	"DragDrop-Files/pkg/idgen"
	"fmt"
	"time"

	"github.com/Aurivena/spond/v2/envelope"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (a *File) downloadZipFile(id, sessionID string, files []entity.File) (*entity.FileSaveOutput, *envelope.AppError) {
	fileIDZip := fmt.Sprintf("%s%s", domain.PrefixZipFile, id)
	zipData, err := archive.ZipFiles(files, fileIDZip)
	if err != nil {
		logrus.Error("failed to zip files")
		return nil, a.InternalServerError()
	}

	generatedID, err := idgen.GenerateID()
	if err != nil {
		logrus.Error("failed to generate id")
		return nil, a.InternalServerError()
	}

	if err = a.postgresql.FileDelete.FileID(fileIDZip); err != nil {
		logrus.Error("failed to delete id file")
		return nil, a.InternalServerError()
	}

	zipFile := entity.File{
		FileID:    fileIDZip,
		Name:      fmt.Sprintf("%s%s", generatedID, domain.MimeTypeZip),
		SessionID: sessionID,
		MimeType:  domain.MimeTypeZip,
	}

	meta, err := a.downloadFile(zipData, zipFile)
	if err != nil {
		return nil, a.InternalServerError()
	}
	out := entity.FileSaveOutput{
		ID:    id,
		Size:  meta.Size,
		Count: len(files),
	}
	return &out, nil
}

func (a *File) downloadFile(data []byte, file entity.File) (*minio.UploadInfo, error) {
	id, err := a.postgresql.FileSave.Execute(file)
	if err != nil {
		logrus.Error("failed to save metadata")
		return nil, err
	}

	file.ID = id

	if err = a.postgresql.FileSave.ExecuteSession(file); err != nil {
		logrus.Error("failed to save session")
		return nil, err
	}

	week := time.Now().Add(time.Duration(24*time.Hour) * 7)
	file.TimeDeleted = &week
	if err = a.postgresql.FileSave.ExecuteParameters(file); err != nil {
		logrus.Error("failed to save parameters")
		return nil, err
	}

	meta, err := a.minioStorage.Save.File(data, fmt.Sprintf("%s/%s", file.SessionID, file.Name))
	if err != nil {
		return nil, err
	}

	return meta, nil
}
