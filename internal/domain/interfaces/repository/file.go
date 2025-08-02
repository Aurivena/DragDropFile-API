package repository

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"time"
)

type Delete interface {
	FilesBySessionID(sessionID string) error
	FilesByFileID(id string) error
	File(id int) error
}

type Get interface {
	FilesBySessionNotZip(sessionID string) ([]entity.FileOutput, error)
	IdFilesBySession(sessionID string) ([]string, error)
	ByID(id string) (*entity.FileOutput, error)
	DataFile(id string) (*entity.DataOutput, error)
	ZipMetaBySession(sessionID string) (*entity.FileOutput, error)
}

type Save interface {
	File(ctx context.Context, input entity.FileSave) error
}

type Update interface {
	CountDownload(count int, session string) error
	DateDeleted(dateDeleted time.Time, id string) error
	Password(password string, session string) error
	Description(description string, session string) error
}
