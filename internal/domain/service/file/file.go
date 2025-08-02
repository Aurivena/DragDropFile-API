package file

import (
	"DragDrop-Files/internal/domain/interfaces/repository"
	"DragDrop-Files/internal/infrastructure/minio"
)

type File struct {
	repo  repository.Repository
	minio minio.Minio
}
