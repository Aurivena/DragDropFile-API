package domain

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/persistence"

	"github.com/minio/minio-go/v7"
)

type Minio interface {
	Save(input *model.FileSave) (string, error)
	Delete(id string) error
	GetByID(id string) (*model.GetFileOutput, error)
}

type Domain struct {
	Minio
}

func NewDomain(persistence *persistence.Persistence, cfg *model.ConfigService, minioClient *minio.Client) *Domain {
	return &Domain{
		Minio: NewMinioService(minioClient, persistence, cfg),
	}
}
