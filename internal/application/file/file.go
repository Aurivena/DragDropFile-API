package file

import (
	"DragDrop-Files/internal/application/ports/repository"
	"DragDrop-Files/internal/infrastructure/minio"
)

type File struct {
	writer       repository.FileWriter
	reader       repository.FileReader
	updater      repository.FileUpdater
	deleted      repository.FileDeleted
	minioStorage *minio.Minio
}

func NewFile(writer repository.FileWriter, reader repository.FileReader, updater repository.FileUpdater, deleted repository.FileDeleted, minioStorage *minio.Minio) *File {
	return &File{
		writer:       writer,
		reader:       reader,
		updater:      updater,
		deleted:      deleted,
		minioStorage: minioStorage,
	}
}
