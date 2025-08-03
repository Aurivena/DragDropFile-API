package application

import (
	"DragDrop-Files/internal/application/file"
	"DragDrop-Files/internal/application/interfaces"
	"DragDrop-Files/internal/domain/service"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
)

type Application struct {
	interfaces.FileSave
	interfaces.FileGet
	interfaces.FileUpdate
}

func New(post *postgres.Repository, mi *minio.Minio, srv *service.Service) *Application {
	return &Application{
		FileSave: file.New(post, mi, srv),
	}
}
