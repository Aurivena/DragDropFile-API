package entity

import (
	"time"

	"github.com/minio/minio-go/v7"
)

type File struct {
	ID            int `json:"id" db:"id"`
	Prefix        string
	FileBase64    string
	FileID        string     `json:"file_id" db:"file_id"`
	Name          string     `json:"name" db:"name"`
	MimeType      string     `json:"mimeType" db:"mime_type"`
	SessionID     string     `json:"session" db:"session"`
	Password      *string    `json:"password" db:"password"`
	TimeDeleted   *time.Time `json:"timeDeleted" db:"date_deleted"`
	CountDownload int        `json:"countDownload" db:"count_download"`
	Description   string     `json:"description" db:"description"`
}

type FileSaveOutput struct {
	Size  int64  `json:"size"`
	Count int    `json:"count"`
	ID    string `json:"id"`
}

type GetFileOutput struct {
	File        *minio.Object `json:"file"`
	Name        string        `json:"name"`
	Description string        `json:"description" db:"description"`
}

type FileData struct {
	Password      bool      `json:"password" db:"password"`
	DateDeleted   time.Time `json:"dateDeleted" db:"date_deleted"`
	CountDownload int       `json:"countDownload" db:"count_download"`
	Description   string    `json:"description" db:"description"`
}

type FileUpdateInput struct {
	CountDayToDeleted int     `json:"countDayToDeleted,omitempty"`
	Password          *string `json:"password,omitempty"`
	CountDownload     *int    `json:"countDownload,omitempty"`
	Description       *string `json:"description,omitempty" db:"description"`
}
