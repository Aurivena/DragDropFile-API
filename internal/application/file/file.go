package file

import (
	"DragDrop-Files/internal/domain/service"
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"
)

type File struct {
	post *postgres.Repository
	mi   *minio.Minio
	srv  *service.Service
}

func New(post *postgres.Repository, mi *minio.Minio, srv *service.Service) *File {
	return &File{
		post: post,
		mi:   mi,
		srv:  srv,
	}
}
