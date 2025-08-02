package file

import (
	"github.com/jmoiron/sqlx"
	"time"
)

type Update struct {
	db *sqlx.DB
}

func NewUpdate(db *sqlx.DB) *Update {
	return &Update{db: db}
}

func (r *Update) CountDownload(count int, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET count_download = $1 WHERE session = $2`, count, session)
	if err != nil {
		return err
	}
	return nil
}

func (r *Update) Description(description string, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET description = $1 WHERE session = $2`, description, session)
	if err != nil {
		return err
	}
	return nil
}

func (r *Update) DateDeleted(dateDeleted time.Time, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET date_deleted = $1 WHERE session = $2`, dateDeleted, session)
	if err != nil {
		return err
	}
	return nil
}
func (r *Update) Password(password string, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET password = $1 WHERE session = $2`, password, session)
	if err != nil {
		return err
	}
	return nil
}
