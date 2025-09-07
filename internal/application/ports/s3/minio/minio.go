package minio

import (
	"DragDrop-Files/internal/domain/entity"

	"github.com/minio/minio-go/v7"
)

type Delete interface {
	DelByFilename(filename string) error
}

type Reader interface {
	ByFilename(filename string) (*entity.GetFileOutput, error)
}

type Writer interface {
	File(data []byte, filename string) (*minio.UploadInfo, error)
}
