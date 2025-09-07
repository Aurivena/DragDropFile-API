package s3_minio

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

func (s *MinioS3) ByFilename(filename string) (*entity.GetFileOutput, error) {
	var out entity.GetFileOutput
	optsGet := minio.GetObjectOptions{}
	objectReader, err := s.minioClient.GetObject(context.Background(), s.cfg.MinioBucketName, filename, optsGet)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			logrus.Errorf("Объект '%s' не найден в MinIO (Bucket: '%s').", filename, s.cfg.MinioBucketName)
			return nil, fmt.Errorf("файл с ID '%s' не найден в хранилище: %w", filename, err)
		}

		logrus.Errorf("Ошибка получения потока объекта из MinIO: Bucket='%s', Object='%s', Err: %v", s.cfg.MinioBucketName, filename, err)
		return nil, fmt.Errorf("ошибка получения содержимого файла: %w", err)
	}

	logrus.Infof("Поток для объекта '%s' из бакета '%s' успешно получен.", filename, s.cfg.MinioBucketName)
	d, _ := objectReader.Stat()
	logrus.Infof("filename %d", d.Size)
	logrus.Infof("filename %s", d.Key)
	logrus.Infof("filename %v", d.Internal)
	out.File = objectReader
	out.Name = filename
	return &out, nil
}
