package file

import (
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/pkg/archive"
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (a *File) downloadZipFile(ctx context.Context, id, sessionID, prefixZipFile string, files []entity.FileFFF) (*minio.UploadInfo, error) {
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

	zipUniqueName := fmt.Sprintf("%s.zip", uuidV7.String())
	meta, err := a.downloadFile(ctx, zipData, ".zip", zipUniqueName, sessionID, fileIDZip)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (a *File) downloadFile(ctx context.Context, data []byte, mimeType, filename, sessionID, id string) (*minio.UploadInfo, error) {
	input := entity.File{
		FileID:    id,
		Name:      filename,
		SessionID: sessionID,
		MimeType:  mimeType,
	}

	if err := a.postgresql.FileSave.Execute(ctx, input); err != nil {
		logrus.Error("failed to save g metadata")
		return nil, err
	}

	meta, err := a.minioStorage.Save.File(data, sessionID, filename)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return meta, nil
}
