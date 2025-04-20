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
	if input.File.DateDeleted != nil {
		t := dateCreated.Add(time.Duration(*input.File.DateDeleted))
		del = &t
	}
	_, err := p.db.Exec(`INSERT INTO "File"  (id, name,password, date_created, date_deleted, count_download, count_discoveries, count_day) VALUES($1,$2,$3,$4,$5,$6,$7,$8)`,
		id, input.Name, input.File.Password, dateCreated, del, input.File.CountDownload, input.File.CountDiscoveries, input.File.CountDay)
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
	var out model.File

	err := p.db.Get(&out, `SELECT * FROM "File" WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
