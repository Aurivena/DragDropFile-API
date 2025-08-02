package file

import (
	"DragDrop-Files/internal/domain/entity"
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"net/http"
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

func (s *Save) File(data []byte, sessionID, filename string) (*minio.UploadInfo, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("пустой файл")
	}

	ctx := context.Background()
	path := fmt.Sprintf("%s/%s", sessionID, filename)

	_, err := s.minioClient.StatObject(ctx, s.cfg.MinioBucketName, path, minio.StatObjectOptions{})
	if err == nil {
		return nil, fmt.Errorf("file duplicate")
	}
	if minio.ToErrorResponse(err).Code != "NoSuchKey" {
		return nil, fmt.Errorf("ошибка при проверке существования файла: %w", err)
	}

	contentType := http.DetectContentType(data)

	reader := bytes.NewReader(data)
	fileSize := int64(len(data))

	uploadInfo, err := s.minioClient.PutObject(ctx, s.cfg.MinioBucketName, path, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	return &uploadInfo, nil
}
