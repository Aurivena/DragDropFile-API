package persistence

import (
	"DragDrop-Files/models"
	"time"

	"github.com/jmoiron/sqlx"
)

type File interface {
	Create(input models.FileSave) error
	Delete(id string) error
	GetZipMetaBySession(sessionID string) (*models.FileOutput, error)
	DeleteFilesBySessionID(sessionID string) error
	DeleteFilesByFileID(id string) error
	UpdateCountDownload(count int, id string) error
	UpdateDateDeleted(dateDeleted time.Time, id string) error
	UpdatePassword(password string, id string) error
	GetFilesBySessionNotZip(sessionID string) ([]models.FileOutput, error)
	GetIdFilesBySession(sessionID string) ([]string, error)
	GetByID(id string) (*models.FileOutput, error)
	GetDataFile(id string) (*models.DataOutput, error)
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
