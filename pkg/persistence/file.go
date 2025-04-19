package persistence

import (
	"DragDrop-Files/model"
	"time"

	"github.com/jmoiron/sqlx"
)

type FilePersistence struct {
	db *sqlx.DB
}

func NewFiePersistence(db *sqlx.DB) *FilePersistence {
	return &FilePersistence{db: db}
}

func (p *FilePersistence) Save(id string, input *model.FileSave) (bool, error) {
	var (
		dateCreated = time.Now().UTC()
		del         *time.Time
	)
	if input.DateDeleted != nil {
		t := dateCreated.Add(time.Duration(*input.DateDeleted))
		del = &t
	}
	_, err := p.db.Exec(`INSERT INTO "File"  (id, date_created, date_deleted, count_download, count_discoveries, count_day) VALUES($1,$2,$3,$4,$5,$6)`,
		id, dateCreated, del, input.CountDownload, input.CountDiscoveries, input.CountDay)
	if err != nil {
		return false, nil
	}

	return true, nil
}

func (p *FilePersistence) Delete(id string) error {
	_, err := p.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *FilePersistence) Get(id string) (*model.File, error) {
	var out *model.File

	err := p.db.Get(&out, `SELECT * FROM "File" WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return out, nil
}
