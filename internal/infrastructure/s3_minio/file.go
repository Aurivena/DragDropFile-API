package s3_minio

import (
	"DragDrop-Files/internal/domain/entity"

	"github.com/minio/minio-go/v7"
)

type MinioS3 struct {
	minioClient *minio.Client
	cfg         *entity.MinioConfig
}

func New(minioClient *minio.Client, cfg *entity.MinioConfig) *MinioS3 {
	return &MinioS3{
		minioClient: minioClient,
		cfg:         cfg,
	}
}
