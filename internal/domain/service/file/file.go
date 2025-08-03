package file

import (
	"DragDrop-Files/internal/domain/interfaces/repository"
	"DragDrop-Files/internal/infrastructure/minio"
)

type File struct {
	repo  repository.File
	minio minio.Minio
}

func New(repo repository.File, minio minio.Minio) *File {
	return &File{
		repo:  repo,
		minio: minio,
	}
}
