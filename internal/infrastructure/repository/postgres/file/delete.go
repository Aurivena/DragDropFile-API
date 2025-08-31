package file

func (r *File) FilesByFileID(id string) error {
	_, err := r.db.Exec(`DELETE FROM "Get" WHERE id = $1`, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *File) FilesBySessionID(sessionID string) error {
	_, err := r.db.Exec(`DELETE FROM "Get"
		USING "SessionID"
		WHERE "SessionID".session = $1
		AND "Get".id = "SessionID".file_id`, sessionID)
	if err != nil {
		return err
	}
	return nil
}
func (r *File) File(id int) error {
	_, err := r.db.Exec(`DELETE FROM "Get" WHERE id = $1`, id)
	if err != nil {
		return err
	}

	return nil
}
