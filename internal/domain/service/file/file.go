package file

import (
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
)

type File struct {
	repo  *postgres.Repository
	minio *minio.Minio
}

func New(repo *postgres.Repository, minio *minio.Minio) *File {
	return &File{
		repo:  repo,
		minio: minio,
	}
}
