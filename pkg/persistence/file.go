package persistence

import (
	"DragDrop-Files/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
)

const (
	errorNoSqlResult = "sql: no rows in result set"
)

var (
	dateDeleted   = time.Now().AddDate(1, 0, 0).UTC()
	countDownload = 365
)

type FilePersistence struct {
	db *sqlx.DB
}

func (p *FilePersistence) GetByID(id string) (*models.FileOutput, error) {
	var out models.FileOutput
	err := p.db.Get(&out, `SELECT S.file_id,name,mime_type,session,password,date_deleted,count_download FROM "File"
			INNER JOIN public."File_Parameters" FP on "File".id = FP.file_id
			INNER JOIN "Session" S ON S.file_id = "File".id
			WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &out, err
}

func (p *FilePersistence) DeleteFilesByFileID(id string) error {
	_, err := p.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func NewFiePersistence(db *sqlx.DB) *FilePersistence {
	return &FilePersistence{db: db}
}

func (p *FilePersistence) Create(ctx context.Context, input models.FileSave) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `INSERT INTO "File" (id, name, mime_type) VALUES ($1,$2,$3)`, input.Id, input.Name, input.MimeType)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO "Session" (file_id, session) VALUES ($1,$2)`, input.Id, input.SessionID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO "File_Parameters" (file_id, date_deleted,count_download,password) VALUES ($1,$2,$3,$4)`, input.Id, dateDeleted, countDownload, nil)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (p *FilePersistence) Delete(id string) error {
	_, err := p.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (p *FilePersistence) GetZipMetaBySession(sessionID string) (*models.FileOutput, error) {
	var out models.FileOutput
	err := p.db.Get(&out, `SELECT file_id, name FROM "File" 
                INNER JOIN "Session" ON "Session".file_id = "File".id
    			WHERE "Session".session = $1 AND "File".name LIKE '%.zip'`, sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &out, nil
}

func (p *FilePersistence) DeleteFilesBySessionID(sessionID string) error {
	_, err := p.db.Exec(`DELETE FROM "File"
		USING "Session"
		WHERE "Session".session = $1
		AND "File".id = "Session".file_id`, sessionID)
	if err != nil {
		return err
	}
	return nil
}

func (p *FilePersistence) UpdateCountDownload(count int, id string) error {
	_, err := p.db.Exec(`UPDATE "File_Parameters" SET count_download = $1 WHERE file_id = $2`, count, id)
	if err != nil {
		return err
	}
	return nil
}
func (p *FilePersistence) UpdateDateDeleted(dateDeleted time.Time, id string) error {
	_, err := p.db.Exec(`UPDATE "File_Parameters" SET date_deleted = $1 WHERE file_id = $2`, dateDeleted, id)
	if err != nil {
		return err
	}
	return nil
}
func (p *FilePersistence) UpdatePassword(password string, id string) error {
	_, err := p.db.Exec(`UPDATE "File_Parameters" SET password = $1 WHERE file_id = $2`, password, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *FilePersistence) GetIdFilesBySession(sessionID string) ([]string, error) {
	var out []string

	err := p.db.Select(&out, `SELECT file_id FROM "Session"
               INNER JOIN public."File" F on F.id = "Session".file_id
               WHERE session = $1 AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil && err.Error() != errorNoSqlResult {
		return nil, err
	}

	return out, nil
}

func (p *FilePersistence) GetFilesBySessionNotZip(sessionID string) ([]models.FileOutput, error) {
	var out []models.FileOutput

	err := p.db.Select(&out, `SELECT F.file_id,name,mime_type,session,password,date_deleted,count_download FROM "Session"
         INNER JOIN public."File_Parameters" F on F.file_id = "Session".file_id
         INNER JOIN public."File" F2 on F2.id = "Session".file_id
         WHERE session = $1
         AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (p *FilePersistence) GetDataFile(id string) (*models.DataOutput, error) {

	var out models.DataOutput
	err := p.db.Get(&out, `SELECT (password IS NOT NULL AND password != '') AS password,date_deleted,count_download FROM "File_Parameters" WHERE file_id =$1`, id)
	if err != nil {
		return nil, err
	}
	return &out, err
}
