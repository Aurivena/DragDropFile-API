package s3_minio

import (
	minioS3 "DragDrop-Files/internal/application/ports/s3/minio"
	"DragDrop-Files/internal/domain/entity"

	"github.com/jmoiron/sqlx"
	"github.com/minio/minio-go/v7"
)

type S3Minio struct {
	minioS3.Reader
	minioS3.Writer
	minioS3.Delete
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func NewS3(client *minio.Client, cfg *entity.MinioConfig) *S3Minio {
	return &S3Minio{
		Delete: New(client, cfg),
		Writer: New(client, cfg),
		Reader: New(client, cfg),
	}
}
