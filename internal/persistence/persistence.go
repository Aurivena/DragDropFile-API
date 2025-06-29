package persistence

import (
	"DragDrop-Files/models"
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type File interface {
	Create(ctx context.Context, input models.FileSave) error
	Delete(id int) error
	GetZipMetaBySession(sessionID string) (*models.FileOutput, error)
	DeleteFilesBySessionID(sessionID string) error
	DeleteFilesByFileID(id string) error
	UpdateCountDownload(count int, session string) error
	UpdateDateDeleted(dateDeleted time.Time, id string) error
	UpdatePassword(password string, session string) error
	UpdateDescription(description string, session string) error
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
