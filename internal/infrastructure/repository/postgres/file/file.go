package file

import "github.com/jmoiron/sqlx"

type File struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *File {
	return &File{
		db: db,
	}
}
