package s3

import "DragDrop-Files/internal/domain/interfaces/s3/minio"

type S3 struct {
	minio.Save
	minio.Get
	minio.Delete
}
