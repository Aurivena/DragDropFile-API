package postgres

import (
	"DragDrop-Files/internal/application/ports/repository"
	"DragDrop-Files/internal/infrastructure/repository/postgres/file"

	"github.com/jmoiron/sqlx"
)

type Repository struct {
	FileReader  repository.FileReader
	FileWriter  repository.FileWriter
	FileDeleted repository.FileDeleted
	FileUpdater repository.FileUpdater
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func New(sources *Sources) *Repository {
	return &Repository{
		FileReader:  file.New(sources.BusinessDB),
		FileWriter:  file.New(sources.BusinessDB),
		FileDeleted: file.New(sources.BusinessDB),
		FileUpdater: file.New(sources.BusinessDB),
	}
}
