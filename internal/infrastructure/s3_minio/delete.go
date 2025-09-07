package s3_minio

import (
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (s *MinioS3) DelByFilename(filename string) error {
	opts := minio.RemoveObjectOptions{}
	err := s.minioClient.RemoveObject(context.Background(), s.cfg.MinioBucketName, filename, opts)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			logrus.Errorf("Объект '%s' уже удален или не существовал в MinIO (Bucket: '%s'). Продолжаем удаление метаданных.", filename, s.cfg.MinioBucketName)
			return fmt.Errorf("объект '%s' уже удален или не существовал в MinIO (Bucket: '%s'). Продолжаем удаление метаданных", filename, s.cfg.MinioBucketName)
		} else {
			logrus.Errorf("Ошибка удаления объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.MinioBucketName, filename, err)
			return fmt.Errorf("ошибка удаления файла из хранилища: %w", err)
		}
	}
	logrus.Infof("Объект '%s' успешно удален из MinIO (Bucket: '%s').", filename, s.cfg.MinioBucketName)

	return nil
}
