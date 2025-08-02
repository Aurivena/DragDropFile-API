package file

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
)

type Delete struct {
	minioClient *minio.Client
	cfg         *entity.MinioConfig
}

func NewDelete(minioClient *minio.Client, cfg *entity.MinioConfig) *Delete {
	return &Delete{
		minioClient: minioClient,
		cfg:         cfg,
	}
}

func (s *Delete) File(filename string) error {
	opts := minio.RemoveObjectOptions{}
	err := s.minioClient.RemoveObject(context.Background(), s.cfg.MinioBucketName, filename, opts)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			log.Printf("Объект '%s' уже удален или не существовал в MinIO (Bucket: '%s'). Продолжаем удаление метаданных.", filename, s.cfg.MinioBucketName)
			return fmt.Errorf("объект '%s' уже удален или не существовал в MinIO (Bucket: '%s'). Продолжаем удаление метаданных", filename, s.cfg.MinioBucketName)
		} else {
			log.Printf("Ошибка удаления объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.MinioBucketName, filename, err)
			return fmt.Errorf("ошибка удаления файла из хранилища: %w", err)
		}
	} else {
		log.Printf("Объект '%s' успешно удален из MinIO (Bucket: '%s').", filename, s.cfg.MinioBucketName)
	}
	return nil
}
