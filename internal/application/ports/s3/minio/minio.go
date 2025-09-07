package minio

import (
	"DragDrop-Files/internal/domain/entity"

	"github.com/minio/minio-go/v7"
)

type Delete interface {
	ByFilename(filename string) error
}

type Get interface {
	ByFilename(filename string) (*entity.GetFileOutput, error)
}

type Save interface {
	File(data []byte, filename string) (*minio.UploadInfo, error)
}
