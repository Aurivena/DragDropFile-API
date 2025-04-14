package domain

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/persistence"

	"github.com/minio/minio-go/v7"
)

type Minio interface {
	Save(input *model.FileSave) (string, error)
	Delete(id string) error
	Get(id string) (*model.File, error)
}

type Domain struct {
	Minio
}

func NewDomain(persistence *persistence.Persistence, minioClient *minio.Client) *Domain {
	return &Domain{
		Minio: NewMinioService(minioClient, persistence),
	}
}
