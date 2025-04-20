package persistence

import (
	"DragDrop-Files/model"

	"github.com/jmoiron/sqlx"
)

type File interface {
	Save(id string, input *model.FileSave) (bool, error)
	Delete(id string) error
	Get(id string) (*model.File, error)
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
