package postgres

import (
	"DragDrop-Files/internal/domain/interfaces/repository"
	file2 "DragDrop-Files/internal/infrastructure/repository/postgres/file"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	repository.Get
	repository.Save
	repository.Delete
	repository.Update
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func NewRepository(sources *Sources) *Repository {
	return &Repository{
		Get:    file2.NewGet(sources.BusinessDB),
		Save:   file2.NewSave(sources.BusinessDB),
		Delete: file2.NewDelete(sources.BusinessDB),
		Update: file2.NewUpdate(sources.BusinessDB),
	}
}
