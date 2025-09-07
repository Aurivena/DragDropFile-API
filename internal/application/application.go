package application

import (
	"DragDrop-Files/internal/application/file"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
	"DragDrop-Files/internal/infrastructure/s3_minio"
)

type Application struct {
	File *file.File
}

func New(postgresql *postgres.Repository, minioStorage *s3_minio.S3Minio) *Application {
	return &Application{
		File: file.NewFile(postgresql.FileWriter, postgresql.FileReader, postgresql.FileUpdater, postgresql.FileDeleted, minioStorage),
	}
}
