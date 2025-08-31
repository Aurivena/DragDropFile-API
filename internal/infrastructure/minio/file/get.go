package file

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
)

type Get struct {
	minioClient *minio.Client
	cfg         *entity.MinioConfig
}

func NewGet(minioClient *minio.Client, cfg *entity.MinioConfig) *Get {
	return &Get{
		minioClient: minioClient,
		cfg:         cfg,
	}
}

func (s *Get) ByFilename(path string) (*entity.GetFileOutput, error) {
	var out entity.GetFileOutput
	optsGet := minio.GetObjectOptions{}
	objectReader, err := s.minioClient.GetObject(context.Background(), s.cfg.MinioBucketName, path, optsGet)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			log.Printf("Объект '%s' не найден в MinIO (Bucket: '%s').", path, s.cfg.MinioBucketName)
			return nil, fmt.Errorf("файл с ID '%s' не найден в хранилище: %w", path, err)
		}

		log.Printf("Ошибка получения потока объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.MinioBucketName, path, err)
		return nil, fmt.Errorf("ошибка получения содержимого файла: %w", err)
	}

	log.Printf("Поток для объекта '%s' из бакета '%s' успешно получен.", path, s.cfg.MinioBucketName)
	out.File = objectReader
	out.Name = path
	return &out, nil
}
