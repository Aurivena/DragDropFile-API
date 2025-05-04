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
	GetNameByID(id string) (string, error)
	Delete(id string) error
	GetMimeTypeByID(id string) (string, error)
	DeleteFilesBySessionID(sessionID string) error
	ValidatePassword(input *models.FileGet) error
	ValidateDateDeleted(sessionID string) error
	ValidateCountDownload(sessionID string) error
	GetSessionByID(id string) (string, error)
	UpdateCountDownload(count int, sessionID string) error
	UpdateDateDeleted(countDayToDeleted int, sessionID string) error
	UpdatePassword(password, sessionID string) error
	GetFileBySession(sessionID string) ([]models.FileOutput, error)
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
