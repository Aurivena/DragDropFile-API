package file

import (
	"github.com/jmoiron/sqlx"
)

type Delete struct {
	db *sqlx.DB
}

func NewDelete(db *sqlx.DB) *Delete {
	return &Delete{
		db: db,
	}
}

func (r *Delete) FilesByFileID(id string) error {
	_, err := r.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Delete) FilesBySessionID(sessionID string) error {
	_, err := r.db.Exec(`DELETE FROM "File"
		USING "Session"
		WHERE "Session".session = $1
		AND "File".id = "Session".file_id`, sessionID)
	if err != nil {
		return err
	}
	return nil
}
func (r *Delete) File(id int) error {
	_, err := r.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
