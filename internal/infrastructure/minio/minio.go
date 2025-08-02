package minio

import (
	"DragDrop-Files/internal/domain/entity"
	minio2 "DragDrop-Files/internal/domain/interfaces/s3/minio"
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
	minio2.Get
	minio2.Save
	minio2.Delete
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func NewMinio(cfg *entity.MinioConfig, minioClient *minio.Client) *Minio {
	return &Minio{
		Get:    file.NewGet(minioClient, cfg),
		Save:   file.NewSave(minioClient, cfg),
		Delete: file.NewDelete(minioClient, cfg),
	}
}
