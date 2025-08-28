package s3

import (
	"DragDrop-Files/internal/application/ports/s3/minio"
)

type S3 struct {
	minio.Save
	minio.Get
	minio.Delete
}
