package domain

import (
	"DragDrop-Files/models"
	"DragDrop-Files/pkg/persistence"
	"github.com/minio/minio-go/v7"
)

type Minio interface {
	Delete(filename string) error
	GetByFilename(path string) (*models.GetFileOutput, error)
	DownloadMinio(data []byte, sessionID, name string) (*minio.UploadInfo, error)
}

type File interface {
	GetZipMetaBySession(sessionID string) (*models.FileOutput, error)
	GetIdFileBySession(sessionID string) ([]string, error)
	Create(input models.FileSave) error
	ZipFiles(files []models.File, id string) ([]byte, error)
	Delete(id string) error
	DeleteFilesBySessionID(sessionID string) error
	ValidatePassword(input *models.FileGet) error
	ValidateDateDeleted(id string) error
	ValidateCountDownload(id string) error
	UpdateCountDownload(count int, id string) error
	UpdateDateDeleted(countDayToDeleted int, id string) error
	UpdatePassword(password, id string) error
	GetFilesBySession(sessionID string) ([]models.FileOutput, error)
	GetDataFile(id string) (*models.DataOutput, error)
	GetByID(id string) (*models.FileOutput, error)
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
