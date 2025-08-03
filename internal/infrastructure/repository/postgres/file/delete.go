package file

func (r *File) FilesByFileID(id string) error {
	_, err := r.db.Exec(`DELETE FROM "File" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *File) FilesBySessionID(sessionID string) error {
	_, err := r.db.Exec(`DELETE FROM "File"
		USING "Session"
		WHERE "Session".session = $1
		AND "File".id = "Session".file_id`, sessionID)
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
