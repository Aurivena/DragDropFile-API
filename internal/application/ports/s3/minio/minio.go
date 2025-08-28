package minio

import (
	"DragDrop-Files/internal/domain/entity"
	"github.com/minio/minio-go/v7"
)

type Delete interface {
	File(filename string) error
}

type Get interface {
	ByFilename(path string) (*entity.GetFileOutput, error)
}

type Save interface {
	File(data []byte, sessionID, name string) (*minio.UploadInfo, error)
}
