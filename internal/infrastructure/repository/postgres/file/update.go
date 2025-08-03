package file

import (
	"time"
)

func (r *File) CountDownload(count int, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET count_download = $1 WHERE session = $2`, count, session)
	if err != nil {
		return err
	}
	return nil
}

func (r *File) Description(description string, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET description = $1 WHERE session = $2`, description, session)
	if err != nil {
		return err
	}
	return nil
}

func (r *File) DateDeleted(dateDeleted time.Time, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET date_deleted = $1 WHERE session = $2`, dateDeleted, session)
	if err != nil {
		return err
	}
	return nil
}
func (r *File) Password(password string, session string) error {
	_, err := r.db.Exec(`UPDATE "File_Parameters" SET password = $1 WHERE session = $2`, password, session)
	if err != nil {
		return err
	}
	return nil
}
