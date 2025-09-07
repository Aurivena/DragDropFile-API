package repository

import (
	"DragDrop-Files/internal/domain/entity"
	"time"
)

type FileDelete interface {
	FilesBySessionID(sessionID string) error
	ID(id int) error
	FileID(fileID string) error
}

type FileGet interface {
	FilesBySessionNotZip(sessionID string) ([]entity.File, error)
	IdFilesBySession(sessionID string) ([]string, error)
	ByID(id string) (*entity.File, error)
	DataFile(id string) (*entity.FileData, error)
	ZipMetaBySession(sessionID string) (*entity.File, error)
}

type FileSave interface {
	Execute(file entity.File) (int, error)
	ExecuteSession(file entity.File) error
	ExecuteParameters(file entity.File) error
}

type FileUpdate interface {
	CountDownload(count int, session string) error
	DateDeleted(dateDeleted time.Time, session string) error
	Password(password string, session string) error
	Description(description string, session string) error
}
