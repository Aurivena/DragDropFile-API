package initialization

import (
	"DragDrop-Files/internal/domain/entity"

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
