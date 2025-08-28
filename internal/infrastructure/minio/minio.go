package minio

import (
	minioS3 "DragDrop-Files/internal/application/ports/s3/minio"
	"DragDrop-Files/internal/domain/entity"
	"DragDrop-Files/internal/infrastructure/minio/file"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewMinioStorage(cfg entity.MinioConfig) *minio.Client {
	client, err := minio.New(cfg.Endpoint, &minio.Options{Creds: credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.SSL})
	if err != nil {
		return nil
	}

	return client
}

type Minio struct {
	minioS3.Get
	minioS3.Save
	minioS3.Delete
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func New(cfg *entity.MinioConfig, minioClient *minio.Client) *Minio {
	return &Minio{
		Get:    file.NewGet(minioClient, cfg),
		Save:   file.NewSave(minioClient, cfg),
		Delete: file.NewDelete(minioClient, cfg),
	}
}
