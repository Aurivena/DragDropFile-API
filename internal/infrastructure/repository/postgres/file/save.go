package file

import (
	"DragDrop-Files/internal/domain/entity"
)

const (
	countDownload = 365
)

func (r *File) Execute(file entity.File) (int, error) {
	var id int
	err := r.db.Get(&id, `INSERT INTO "File" (file_id, name, mime_type) VALUES ($1,$2,$3) RETURNING id`, file.FileID, file.Name, file.MimeType)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *File) ExecuteSession(sessionID string, id int) error {

	_, err := r.db.Exec(`INSERT INTO "Session" (file_id, session) VALUES ($1,$2)`, id, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (r *File) ExecuteParameters(file entity.File, currentTime string) error {

	_, err := r.db.Exec(`INSERT INTO "File_Parameters" (file_id,session, date_deleted,count_download,password,description) VALUES ($1,$2,$3,$4,$5,$6)`, file.ID, file.SessionID, currentTime, countDownload, nil, "")
	if err != nil {
		return err
	}

	return nil
}
