package domain

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/persistence"
	"github.com/minio/minio-go/v7"
)

type Minio interface {
	Delete(filename string) error
	GetByFilename(filename string) (*model.GetFileOutput, error)
	DownloadMinio(data []byte, name string) (*minio.UploadInfo, error)
}

type File interface {
	GetZipMetaBySession(sessionID string) (*model.File, error)
	GetIdFileBySession(sessionID string) ([]string, error)
	Create(input model.FileSave) error
	ZipFiles(filesBase64 []string, id string) ([]byte, error)
	GetNameByID(id string) (string, error)
	Delete(id string) error
	GetDataBase64ByID(id string) (string, error)
	DeleteFilesBySessionID(sessionID string) error
	ValidatePassword(input *model.FileGet) error
	ValidateDateDeleted(sessionID string) error
	ValidateCountDownload(sessionID string) error
	GetSessionByID(id string) (string, error)
	UpdateCountDownload(count int, sessionID string) error
	UpdateDateDeleted(countDayToDeleted int, sessionID string) error
	UpdatePassword(password, sessionID string) error
}

type Domain struct {
	Minio
	File
}

func NewDomain(persistence *persistence.Persistence, cfg *model.ConfigService, minioClient *minio.Client) *Domain {
	return &Domain{
		Minio: NewMinioService(minioClient, cfg),
		File:  NewFileService(persistence),
	}
}
