package persistence

import (
	"DragDrop-Files/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type File interface {
	Create(input models.FileSave) error
	Delete(id string) error
	GetNameByID(id string) (string, error)
	GetIdFileBySession(sessionID string) ([]string, error)
	GetZipMetaBySession(sessionID string) (*models.FileOutput, error)
	GetMimeTypeByID(id string) (string, error)
	DeleteFilesBySessionID(sessionID string) error
	Get(sessionID string) (*models.Data, error)
	UpdateCountDownload(count int, sessionID string) error
	UpdateDateDeleted(dateDeleted time.Time, sessionID string) error
	UpdatePassword(password string, sessionID string) error
	GetSessionByID(id string) (string, error)
	GetFileBySession(sessionID string) ([]models.FileOutput, error)
}

type Persistence struct {
	File
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func NewPersistence(sources *Sources) *Persistence {
	return &Persistence{
		File: NewFiePersistence(sources.BusinessDB),
	}
}
