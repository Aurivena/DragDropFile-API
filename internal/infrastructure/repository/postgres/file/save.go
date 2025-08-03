package file

import (
	"DragDrop-Files/internal/domain/entity"
	"context"
	"fmt"
	"time"
)

var (
	dateDeleted   = time.Now().AddDate(1, 0, 0).UTC()
	countDownload = 365
)

func (r *File) Execute(ctx context.Context, input entity.FileSave) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id string
	err = tx.QueryRowContext(ctx, `INSERT INTO "File" (file_id, name, mime_type) VALUES ($1,$2,$3) RETURNING id`, input.Id, input.Name, input.MimeType).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO "Session" (file_id, session) VALUES ($1,$2)`, id, input.SessionID)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, `INSERT INTO "File_Parameters" (file_id,session, date_deleted,count_download,password,description) VALUES ($1,$2,$3,$4,$5,$6)`, id, input.SessionID, dateDeleted, countDownload, nil, "")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
