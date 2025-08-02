package file

import (
	"DragDrop-Files/internal/application/interfaces/file"
	"DragDrop-Files/internal/domain/interfaces/repository"
	"DragDrop-Files/internal/domain/interfaces/service"
	"DragDrop-Files/internal/infrastructure/minio"
)

type File struct {
	g       file.Get
	repo    repository.Repository
	minio   minio.Minio
	service service.Service
}
