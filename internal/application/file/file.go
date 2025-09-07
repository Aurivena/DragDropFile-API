package file

import (
	"DragDrop-Files/internal/application/ports/repository"
	"DragDrop-Files/internal/infrastructure/s3_minio"
)

type File struct {
	writer       repository.FileWriter
	reader       repository.FileReader
	updater      repository.FileUpdater
	deleted      repository.FileDelete
	minioStorage *s3_minio.S3Minio
}

func NewFile(writer repository.FileWriter, reader repository.FileReader, updater repository.FileUpdater, deleted repository.FileDelete, minioStorage *s3_minio.S3Minio) *File {
	return &File{
		writer:       writer,
		reader:       reader,
		updater:      updater,
		deleted:      deleted,
		minioStorage: minioStorage,
	}
}
