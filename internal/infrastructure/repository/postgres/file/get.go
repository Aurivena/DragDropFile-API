package file

import (
	"DragDrop-Files/internal/domain/entity"
	"database/sql"
	"errors"
)

func (r *File) ByID(id string) (*entity.File, error) {
	var out entity.File
	err := r.db.Get(&out, `SELECT S.file_id,name,mime_type,FP.session,password,date_deleted,count_download,description FROM "File"
			INNER JOIN public."File_Parameters" FP on "File".id = FP.file_id
			INNER JOIN "Session" S ON S.file_id = "File".id
			WHERE "File".file_id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &out, err
}

func (r *File) ZipMetaBySession(sessionID string) (*entity.File, error) {
	var out entity.File
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

func (r *File) IdFilesBySession(sessionID string) ([]string, error) {
	var out []string

	err := r.db.Select(&out, `SELECT F.file_id FROM "Session"
               INNER JOIN public."File" F on F.id = "Session".file_id
               WHERE session = $1 AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return nil, err
	}

	return out, nil
}

func (r *File) FilesBySessionNotZip(sessionID string) ([]entity.File, error) {
	var out []entity.File

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

func (r *File) DataFile(id string) (*entity.FileData, error) {
	var out entity.FileData
	err := r.db.Get(&out, `SELECT (password IS NOT NULL AND password != '') AS password,date_deleted,count_download,description
					FROM "File_Parameters"
					INNER JOIN public."File" F on F.id = "File_Parameters".file_id
					WHERE F.file_id =$1`, id)
	if err != nil {
		return nil, err
	}
	return &out, err
}
