package file

func (r *File) FilesByFileID(id string) error {
	_, err := r.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

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
func (r *File) File(id int) error {
	_, err := r.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
