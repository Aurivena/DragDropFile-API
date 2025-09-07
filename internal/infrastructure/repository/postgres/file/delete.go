package file

func (r *File) FilesBySessionID(sessionID string) error {
	_, err := r.db.Exec(`DELETE FROM "Session" s
       USING "File" f
		WHERE s.session = $1
		AND f.id = s.file_id`, sessionID)
	if err != nil {
		return err
	}
	return nil
}
func (r *File) ID(id int) error {
	_, err := r.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *File) FileID(fileID string) error {
	_, err := r.db.Exec(`DELETE FROM "File" WHERE file_id = $1`, fileID)
	if err != nil {
		return err
	}

	return nil
}
