package application

import (
	"DragDrop-Files/internal/application/file"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
)

type Application struct {
	File *file.File
}

func New(postgresql *postgres.Repository, minioStorage *minio.Minio) *Application {
	return &Application{
		File: file.NewFile(postgresql.FileWriter, postgresql.FileReader, postgresql.FileUpdater, postgresql.FileDeleted, minioStorage),
	}
}
