package file

import (
	"DragDrop-Files/internal/domain"
	"DragDrop-Files/internal/domain/entity"
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/minio/minio-go/v7"
)

type Save struct {
	minioClient *minio.Client
	cfg         *entity.MinioConfig
}

func NewSave(minioClient *minio.Client, cfg *entity.MinioConfig) *Save {
	return &Save{
		minioClient: minioClient,
		cfg:         cfg,
	}
}

func (s *Save) File(data []byte, filename string) (*minio.UploadInfo, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("пустой файл")
	}

	ctx := context.Background()

	_, err := s.minioClient.StatObject(ctx, s.cfg.MinioBucketName, filename, minio.StatObjectOptions{})
	if err == nil {
		return nil, domain.ErrFileDuplicate
	}
	if minio.ToErrorResponse(err).Code != "NoSuchKey" {
		return nil, fmt.Errorf("ошибка при проверке существования файла: %w", err)
	}

	uploadInfo, err := s.minioClient.PutObject(ctx, s.cfg.MinioBucketName, filename, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: http.DetectContentType(data),
	})
	if err != nil {
		return nil, err
	}

	return &uploadInfo, nil
}
