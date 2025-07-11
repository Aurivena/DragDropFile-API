package domain

import (
	"DragDrop-Files/internal/persistence"
	"DragDrop-Files/models"
	"context"
	"github.com/minio/minio-go/v7"
)

type Minio interface {
	Delete(filename string) error
	GetByFilename(path string) (*models.GetFileOutput, error)
	DownloadMinio(data []byte, sessionID, name string) (*minio.UploadInfo, error)
}

type File interface {
	GetZipMetaBySession(sessionID string) (*models.FileOutput, error)
	Create(ctx context.Context, input models.FileSave) error
	DeleteFilesBySessionID(sessionID string) error
	ValidatePassword(input *models.FileGet) error
	ValidateDateDeleted(id string) error
	ValidateCountDownload(id string) error
	UpdateCountDownload(count int, id string) error
	UpdateDateDeleted(countDayToDeleted int, id string) error
	UpdateDescription(description, id string) error
	UpdatePassword(password, id string) error
	GetFilesBySession(sessionID string) ([]models.FileOutput, error)
	GetDataFile(id string) (*models.DataOutput, error)
	GetByID(id string) (*models.FileOutput, error)
	Delete(id int) error
	DeleteFiles(id string) error
}

type Domain struct {
	Minio
	File
}

func NewDomain(persistence *persistence.Persistence, cfg *models.ConfigService, minioClient *minio.Client) *Domain {
	return &Domain{
		Minio: NewMinioService(minioClient, cfg),
		File:  NewFileService(persistence),
	}
}
