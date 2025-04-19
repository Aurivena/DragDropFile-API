package domain

import (
	"DragDrop-Files/model"
	"DragDrop-Files/pkg/persistence"
	"archive/zip"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/minio/minio-go/v7"
	"io"
	"log"
	"math/big"
	"mime"
	"net/http"
	"strings"
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
		return "", err
	}
	var name string
	var ext string
	var data []byte

	if len(file.FileBase64) > 1 {
		data, err = zipFiles(file.FileBase64, id)
		if err != nil {
			return "", err
		}

		name = fmt.Sprintf("%s.zip", id)
	} else {
		data, ext, err = decodeFile(file.FileBase64[0])
		if err != nil {
			return "", err
		}
		name = fmt.Sprintf("%s%s", id, ext)
	}

	fileSize := int64(len(data))
	var contentType string
	if fileSize > 0 {
		contentType = http.DetectContentType(data)
	} else {
		contentType = "application/octet-stream"
	}

	file.Name = name
	answer, err := s.pers.File.Save(id, file)
	if err != nil {
		return "", err
	}

	if !answer {
		return "", err
	}

	reader := bytes.NewReader(data)

	uploadInfo, err := s.minioClient.PutObject(context.Background(), s.cfg.Minio.MinioBucketName, file.Name, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", err
	}

	log.Printf("Файл успешно загружен в MinIO: ID=%s, ETag=%s, Size=%d\n", id, uploadInfo.ETag, uploadInfo.Size)

	return id, nil
}

func (s *MinioService) Delete(id string) error {
	file, err := s.pers.Get(id)
	if err != nil {
		return err
	}
	opts := minio.RemoveObjectOptions{}
	err = s.minioClient.RemoveObject(context.Background(), s.cfg.Minio.MinioBucketName, file.Name, opts)
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

func (s *MinioService) GetByID(id string) (*model.GetFileOutput, error) {
	var out model.GetFileOutput
	file, err := s.pers.Get(id)
	if err != nil {
		return nil, err
	}
	optsGet := minio.GetObjectOptions{}
	objectReader, err := s.minioClient.GetObject(context.Background(), s.cfg.Minio.MinioBucketName, file.Name, optsGet)
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
	out.File = objectReader
	out.Name = file.Name
	return &out, nil
}

func generateID() (string, error) {
	lenCode := 12
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, lenCode)
	newInt := big.NewInt(int64(len(letters)))
	for i := range code {
		num, err := rand.Int(rand.Reader, newInt)
		if err != nil {
			log.Printf("не удалось сгенерировать часть ID: %w", err)
			return "", fmt.Errorf("не удалось сгенерировать часть ID: %w", err)
		}
		code[i] = letters[num.Int64()]
	}

	return string(code), nil
}

func decodeFile(fileBase64 string) ([]byte, string, error) {
	mimeType := ""
	base64Data := fileBase64

	if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
		parts := strings.SplitN(fileBase64, ";base64,", 2)
		if len(parts) == 2 {
			base64Data = parts[1]
			mimePart := parts[0]
			if strings.HasPrefix(mimePart, "data:") {
				mimeType = mimePart[len("data:"):]
			}
		}
	}

	decodedData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		log.Printf("Ошибка декодирования Base64 для строки '%s': %v", fileBase64[:min(len(fileBase64), 50)], err) // Добавляем контекст в лог
		return nil, "", fmt.Errorf("некорректные Base64 данные: %w", err)
	}

	exts, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(exts) == 0 {
		log.Printf("Не удалось определить расширение по MIME-типу '%s': %v", mimeType, err)
		return decodedData, "", nil
	}

	return decodedData, exts[0], nil
}

func zipFiles(filesBase64 []string, filename string) ([]byte, error) {
	var buff bytes.Buffer
	zipW := zip.NewWriter(&buff)

	for i, base64Content := range filesBase64 {
		fileBytes, ext, err := decodeFile(base64Content)
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при обработке файла %d: %w", i, err)
		}

		if len(fileBytes) == 0 {
			log.Printf("[zipFiles] Пустой файл #%d. Пропускаем.", i+1)
			continue
		}
		headerName := fmt.Sprintf("%s-%d%s", filename, i+1, ext)
		header := &zip.FileHeader{
			Name:   headerName,
			Method: zip.Deflate,
		}

		fileInZip, err := zipW.CreateHeader(header)
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при создании файла %s в zip-архиве: %w", header.Name, err)
		}

		_, err = io.Copy(fileInZip, bytes.NewReader(fileBytes))
		if err != nil {
			_ = zipW.Close()
			return nil, fmt.Errorf("ошибка при записи содержимого файла %s в zip-архив: %w", header.Name, err)
		}

	}

	err := zipW.Close()
	if err != nil {
		return nil, fmt.Errorf("ошибка при закрытии zip-архива: %w", err)
	}

	return buff.Bytes(), nil
}
