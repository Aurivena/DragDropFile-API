package file

import (
	"DragDrop-Files/internal/infrastructure/minio"
	"DragDrop-Files/internal/infrastructure/repository/postgres"

	"github.com/Aurivena/spond/v2/core"
)

type File struct {
	postgresql   *postgres.Repository
	minioStorage *minio.Minio
	spond        *core.Spond
}

func New(postgresql *postgres.Repository, minioStorage *minio.Minio, spond *core.Spond) *File {
	return &File{
		postgresql:   postgresql,
		minioStorage: minioStorage,
		spond:        spond,
	}
}
