package file

import (
	"DragDrop-Files/internal/domain/entity"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type Get struct {
	db *sqlx.DB
}

func NewGet(db *sqlx.DB) *Get {
	return &Get{
		db: db,
	}
}

func (r *Get) ByID(id string) (*entity.FileOutput, error) {
	var out entity.FileOutput
	err := r.db.Get(&out, `SELECT S.file_id,name,mime_type,FP.session,password,date_deleted,count_download,description FROM "File"
			INNER JOIN public."File_Parameters" FP on "File".id = FP.file_id
			INNER JOIN "Session" S ON S.file_id = "File".id
			WHERE "File".file_id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &out, err
}

func (r *Get) ZipMetaBySession(sessionID string) (*entity.FileOutput, error) {
	var out entity.FileOutput
	err := r.db.Get(&out, `SELECT "File".id,"File".file_id, name FROM "File" 
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

func (r *Get) IdFilesBySession(sessionID string) ([]string, error) {
	var out []string

	err := r.db.Select(&out, `SELECT F.file_id FROM "Session"
               INNER JOIN public."File" F on F.id = "Session".file_id
               WHERE session = $1 AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return nil, err
	}

	return out, nil
}

func (r *Get) FilesBySessionNotZip(sessionID string) ([]entity.FileOutput, error) {
	var out []entity.FileOutput

	err := r.db.Select(&out, `SELECT F.id, F.file_id,name,mime_type,FP.session,password,date_deleted,count_download,description FROM "Session"
		INNER JOIN public."File_Parameters" FP on FP.file_id = "Session".file_id
		INNER JOIN public."File" F on F.id = "Session".file_id
		WHERE FP.session = $1
		  AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (r *Get) DataFile(id string) (*entity.DataOutput, error) {
	var out entity.DataOutput

	err := r.db.Get(&out, `SELECT (password IS NOT NULL AND password != '') AS password,date_deleted,count_download,description
					FROM "File_Parameters"
					INNER JOIN public."File" F on F.id = "File_Parameters".file_id
					WHERE F.file_id =$1`, id)
	if err != nil {
		return nil, err
	}
	return &out, err
}
