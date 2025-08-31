package application

import (
	"DragDrop-Files/internal/application/file"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"

	"github.com/Aurivena/spond/v2/core"
)

type Application struct {
	File *file.File
}

func New(post *postgres.Repository, minioStorage *minio.Minio, spond *core.Spond) *Application {
	return &Application{
		File: file.New(post, minioStorage, spond),
	}
}
