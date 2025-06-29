package domain

import (
	"DragDrop-Files/models"
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"log"
	"net/http"
)

type MinioService struct {
	minioClient *minio.Client
	cfg         *models.ConfigService
}

func NewMinioService(minioClient *minio.Client, cfg *models.ConfigService) *MinioService {
	return &MinioService{minioClient: minioClient, cfg: cfg}
}

func (s *MinioService) Delete(filename string) error {
	opts := minio.RemoveObjectOptions{}
	err := s.minioClient.RemoveObject(context.Background(), s.cfg.Minio.MinioBucketName, filename, opts)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			log.Printf("Объект '%s' уже удален или не существовал в MinIO (Bucket: '%s'). Продолжаем удаление метаданных.", filename, s.cfg.Minio.MinioBucketName)
			return fmt.Errorf("объект '%s' уже удален или не существовал в MinIO (Bucket: '%s'). Продолжаем удаление метаданных", filename, s.cfg.Minio.MinioBucketName)
		} else {
			log.Printf("Ошибка удаления объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.Minio.MinioBucketName, filename, err)
			return fmt.Errorf("ошибка удаления файла из хранилища: %w", err)
		}
	} else {
		log.Printf("Объект '%s' успешно удален из MinIO (Bucket: '%s').", filename, s.cfg.Minio.MinioBucketName)
	}
	return nil
}

func (s *MinioService) GetByFilename(path string) (*models.GetFileOutput, error) {
	var out models.GetFileOutput
	optsGet := minio.GetObjectOptions{}
	objectReader, err := s.minioClient.GetObject(context.Background(), s.cfg.Minio.MinioBucketName, path, optsGet)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			log.Printf("Объект '%s' не найден в MinIO (Bucket: '%s').", path, s.cfg.Minio.MinioBucketName)
			return nil, fmt.Errorf("файл с ID '%s' не найден в хранилище: %w", path, err)
		}

		log.Printf("Ошибка получения потока объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.Minio.MinioBucketName, path, err)
		return nil, fmt.Errorf("ошибка получения содержимого файла: %w", err)
	}

	log.Printf("Поток для объекта '%s' из бакета '%s' успешно получен.", path, s.cfg.Minio.MinioBucketName)
	out.File = objectReader
	out.Name = path
	return &out, nil
}

func (s *MinioService) DownloadMinio(data []byte, sessionID, filename string) (*minio.UploadInfo, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("пустой файл")
	}

	ctx := context.Background()
	path := fmt.Sprintf("%s/%s", sessionID, filename)

	_, err := s.minioClient.StatObject(ctx, s.cfg.Minio.MinioBucketName, path, minio.StatObjectOptions{})
	if err == nil {
		return nil, fmt.Errorf("file duplicate")
	}
	if minio.ToErrorResponse(err).Code != "NoSuchKey" {
		return nil, fmt.Errorf("ошибка при проверке существования файла: %w", err)
	}

	contentType := http.DetectContentType(data)

	reader := bytes.NewReader(data)
	fileSize := int64(len(data))

	uploadInfo, err := s.minioClient.PutObject(ctx, s.cfg.Minio.MinioBucketName, path, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return nil, err
	}

	return &uploadInfo, nil
}
