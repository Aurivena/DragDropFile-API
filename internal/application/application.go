package application

import (
	"DragDrop-Files/internal/application/file"
	"DragDrop-Files/internal/application/interfaces"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"

	"github.com/Aurivena/spond/v2/core"
)

type Application struct {
	file.File
}

func New(post *postgres.Repository, minioStorage *minio.Minio, spond *core.Spond) *Application {
	return &Application{
		FileSave:   file.New(post, minioStorage, spond),
		FileGet:    file.New(post, minioStorage, spond),
		FileUpdate: file.New(post, minioStorage, spond),
	}
}
