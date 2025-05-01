package persistence

import (
	"DragDrop-Files/model"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	errorNoSqlResult = "sql: no rows in result set"
)

type FilePersistence struct {
	db *sqlx.DB
}

func NewFiePersistence(db *sqlx.DB) *FilePersistence {
	return &FilePersistence{db: db}
}

func (p *FilePersistence) Create(input model.FileSave) error {
	_, err := p.db.Exec(`INSERT INTO "File"  (id, name, session, data_base64, password, date_deleted, count_download) VALUES($1,$2,$3,$4,$5,$6,$7)`,
		input.Id, input.Name, input.SessionID, input.DataBase64, nil, nil, nil)
	if err != nil {
		return err
	}

	return nil
}

func (p *FilePersistence) GetMimeTypeByID(id string) (string, error) {
	var dataBase64 string

	err := p.db.Get(&dataBase64, `SELECT data_base64 FROM "File" WHERE id = $1`, id)
	if err != nil {
		return "", err
	}

	return dataBase64, nil
}

func (p *FilePersistence) Delete(id string) error {
	_, err := p.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *FilePersistence) GetNameByID(id string) (string, error) {
	var out string

	err := p.db.Get(&out, `SELECT name FROM "File" WHERE id = $1`, id)
	if err != nil {
		return "", err
	}

	return out, nil
}

func (p *FilePersistence) GetIdFileBySession(sessionID string) ([]string, error) {
	var out []string

	err := p.db.Select(&out, `SELECT id FROM "File" WHERE session = $1 AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil && err.Error() != errorNoSqlResult {
		return nil, err
	}

	return out, nil
}

func (p *FilePersistence) GetZipMetaBySession(sessionID string) (*model.FileOutput, error) {
	var out model.FileOutput
	err := p.db.Get(&out, `SELECT id, name FROM "File" WHERE session = $1 AND name LIKE '%.zip'`, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (p *FilePersistence) DeleteZipMetaBySession(sessionID string) error {
	_, err := p.db.Exec(`DELETE FROM "File" WHERE session = $1 AND name LIKE '%.zip'`, sessionID)
	return err
}

func (p *FilePersistence) DeleteFilesBySessionID(sessionID string) error {
	_, err := p.db.Exec(`DELETE FROM "File" WHERE session = $1`, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (p *FilePersistence) Get(sessionID string) (*model.Data, error) {
	var out model.Data

	err := p.db.Get(&out, `SELECT password,date_deleted,count_download FROM "File" WHERE session = $1`, sessionID)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (p *FilePersistence) UpdateCountDownload(count int, sessionID string) error {
	_, err := p.db.Exec(`UPDATE "File" SET count_download = $1 WHERE session = $2`, count, sessionID)
	if err != nil {
		return err
	}
	return nil
}
func (p *FilePersistence) UpdateDateDeleted(dateDeleted time.Time, sessionID string) error {
	_, err := p.db.Exec(`UPDATE "File" SET date_deleted = $1 WHERE session = $2`, dateDeleted, sessionID)
	if err != nil {
		return err
	}
	return nil
}
func (p *FilePersistence) UpdatePassword(password string, sessionID string) error {
	_, err := p.db.Exec(`UPDATE "File" SET password = $1 WHERE session = $2`, password, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (p *FilePersistence) GetSessionByID(id string) (string, error) {
	var session string

	err := p.db.Get(&session, `SELECT session FROM "File" WHERE id = $1`, id)
	if err != nil {
		return "", err
	}

	return session, nil
}

func (p *FilePersistence) GetFileBySession(sessionID string) ([]model.FileOutput, error) {
	var out []model.FileOutput

	err := p.db.Select(&out, `SELECT id,name FROM "File" WHERE session = $1 AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil {
		return nil, err
	}

	return out, nil
}
