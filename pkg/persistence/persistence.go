package persistence

import "github.com/jmoiron/sqlx"

type Persistence struct {
}

type Sources struct {
	BusinessDB *sqlx.DB
}

func NewPersistence(sources *Sources) *Persistence {
	return &Persistence{}
}
