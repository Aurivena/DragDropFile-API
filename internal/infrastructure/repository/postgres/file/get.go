package file

import (
	"DragDrop-Files/internal/domain/entity"
	"database/sql"
	"errors"
)

func (r *File) ByID(id string) (*entity.File, error) {
	var out entity.File
	err := r.db.Get(&out, `SELECT S.file_id,name,mime_type,FP.session,password,date_deleted,count_download,description FROM "FileFFF"
			INNER JOIN public."File_Parameters" FP on "FileFFF".id = FP.file_id
			INNER JOIN "SessionID" S ON S.file_id = "FileFFF".id
			WHERE "FileFFF".file_id = $1`, id)
	if err != nil {
		return nil, err
	}

	return &out, err
}

func (r *File) ZipMetaBySession(sessionID string) (*entity.File, error) {
	var out entity.File
	err := r.db.Get(&out, `SELECT "FileFFF".id,"FileFFF".file_id, name FROM "FileFFF" 
                INNER JOIN "SessionID" ON "SessionID".file_id = "FileFFF".id
    			WHERE "SessionID".session = $1 AND "FileFFF".name LIKE '%.zip'`, sessionID)
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

	err := r.db.Select(&out, `SELECT F.file_id FROM "SessionID"
               INNER JOIN public."FileFFF" F on F.id = "SessionID".file_id
               WHERE session = $1 AND name NOT LIKE '%.zip'`, sessionID)
	if err != nil && err.Error() != sql.ErrNoRows.Error() {
		return nil, err
	}

	return out, nil
}

func (r *File) FilesBySessionNotZip(sessionID string) ([]entity.File, error) {
	var out []entity.File

	err := r.db.Select(&out, `SELECT F.id, F.file_id,name,mime_type,FP.session,password,date_deleted,count_download,description FROM "SessionID"
		INNER JOIN public."File_Parameters" FP on FP.file_id = "SessionID".file_id
		INNER JOIN public."FileFFF" F on F.id = "SessionID".file_id
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
					INNER JOIN public."FileFFF" F on F.id = "File_Parameters".file_id
					WHERE F.file_id =$1`, id)
	if err != nil {
		return nil, err
	}
	return &out, err
}
