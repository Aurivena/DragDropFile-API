package postgres

import (
	"DragDrop-Files/internal/application/ports/repository"
	"DragDrop-Files/internal/infrastructure/repository/postgres/file"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	repository.FileGet
	repository.FileSave
	repository.FileDelete
	repository.FileUpdate
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func New(sources *Sources) *Repository {
	return &Repository{
		FileGet:    file.New(sources.BusinessDB),
		FileSave:   file.New(sources.BusinessDB),
		FileDelete: file.New(sources.BusinessDB),
		FileUpdate: file.New(sources.BusinessDB),
	}
}
