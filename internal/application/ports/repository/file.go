package repository

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"time"
)

type FileDelete interface {
	FilesBySessionID(sessionID string) error
	FilesByFileID(id string) error
	File(id int) error
}

type FileGet interface {
	FilesBySessionNotZip(sessionID string) ([]entity.File, error)
	IdFilesBySession(sessionID string) ([]string, error)
	ByID(id string) (*entity.File, error)
	DataFile(id string) (*entity.FileData, error)
	ZipMetaBySession(sessionID string) (*entity.File, error)
}

type FileSave interface {
	Execute(ctx context.Context, input entity.File, currentTime string) error
}

type FileUpdate interface {
	CountDownload(count int, session string) error
	DateDeleted(dateDeleted time.Time, id string) error
	Password(password string, session string) error
	Description(description string, session string) error
}
