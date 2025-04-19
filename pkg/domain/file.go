package domain

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/persistence"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
)

type MinioService struct {
	minioClient *minio.Client
	pers        *persistence.Persistence
	cfg         *model.ConfigService
}

func NewMinioService(minioClient *minio.Client, pers *persistence.Persistence, cfg *model.ConfigService) *MinioService {
	return &MinioService{minioClient: minioClient, pers: pers, cfg: cfg}
}

func (s *MinioService) Save(file *model.FileSave) (string, error) {
	id, err := generateID()
	if err != nil {
		return "", nil
	}

	base64Data := file.FileBase64
	if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
		base64Data = base64Data[idx+len(";base64,"):]
	}

	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Printf("Ошибка декодирования Base64 для ID %s: %v", id, err)
		return "", fmt.Errorf("некорректные Base64 данные: %w", err)
	}

	fileSize := int64(len(decodedData))
	var contentType string
	if fileSize > 0 {
		contentType = http.DetectContentType(decodedData)
	} else {
		contentType = "application/octet-stream"
	}

	answer, err := s.pers.File.Save(id, file)
	if err != nil {
		return "", err
	}

	if !answer {
		return "", err
	}

	reader := bytes.NewReader(decodedData)

	uploadInfo, err := s.minioClient.PutObject(context.Background(), s.cfg.Minio.MinioBucketName, id, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	log.Printf("Файл успешно загружен в MinIO: ID=%s, ETag=%s, Size=%d\n", id, uploadInfo.ETag, uploadInfo.Size)

	return id, nil
}

func (s *MinioService) Delete(id string) error {
	opts := minio.RemoveObjectOptions{}
	err := s.minioClient.RemoveObject(context.Background(), s.cfg.Minio.MinioBucketName, id, opts)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			log.Printf("Объект '%s' уже удален или не существовал в MinIO (Bucket: '%s'). Продолжаем удаление метаданных.", id, s.cfg.Minio.MinioBucketName)
		} else {
			log.Printf("Ошибка удаления объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.Minio.MinioBucketName, id, err)
			return fmt.Errorf("ошибка удаления файла из хранилища: %w", err)
		}
	} else {
		log.Printf("Объект '%s' успешно удален из MinIO (Bucket: '%s').", id, s.cfg.Minio.MinioBucketName)
	}
	return s.pers.Delete(id)
}

func (s *MinioService) GetByID(id string) (*minio.Object, error) {
	optsGet := minio.GetObjectOptions{}
	objectReader, err := s.minioClient.GetObject(context.Background(), s.cfg.Minio.MinioBucketName, id, optsGet)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			log.Printf("Объект '%s' не найден в MinIO (Bucket: '%s').", id, s.cfg.Minio.MinioBucketName)
			return nil, fmt.Errorf("файл с ID '%s' не найден в хранилище: %w", id, err)
		}

		log.Printf("Ошибка получения потока объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.Minio.MinioBucketName, id, err)
		return nil, fmt.Errorf("ошибка получения содержимого файла: %w", err)
	}

	log.Printf("Поток для объекта '%s' из бакета '%s' успешно получен.", id, s.cfg.Minio.MinioBucketName)
	return objectReader, nil
}

func generateID() (string, error) {
	lenCode := 12
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, lenCode)
	max := big.NewInt(int64(len(letters)))
	for i := range code {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			log.Printf("не удалось сгенерировать часть ID: %w", err)
			return "", fmt.Errorf("не удалось сгенерировать часть ID: %w", err)
		}
		code[i] = letters[num.Int64()]
	}

	return string(code), nil
}
