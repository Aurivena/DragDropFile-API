package service

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
)

type Validate interface {
	File(zipFileID, id, password string) error
	CheckFilesID(sessionID string) (string, []entity.File, error)
}

type Save interface {
	Execute(ctx context.Context, id, sessionID string, newFiles, oldFiles []entity.File) (*entity.FileSaveOutput, error)
}
