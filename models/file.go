package models

import (
	"github.com/minio/minio-go/v7"
	"time"
)

type SessionOutput struct {
	SessionID string `db:"session"`
}
type File struct {
	FileBase64 string `form:"file"`
	Filename   string `form:"filename"`
}

type FileOutput struct {
	ID            int       `json:"id" db:"id"`
	FileID        string    `json:"file_id" db:"file_id"`
	Name          string    `json:"name" db:"name"`
	MimeType      string    `json:"mimeType" db:"mime_type"`
	Session       string    `json:"session" db:"session"`
	Password      *string   `json:"password" db:"password"`
	DateDeleted   time.Time `json:"dateDeleted" db:"date_deleted"`
	CountDownload int       `json:"countDownload" db:"count_download"`
	Description   string    `json:"description" db:"description"`
}

type FileSave struct {
	Id        string
	Name      string
	SessionID string
	MimeType  string
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

type Data struct {
	Password      *string    `json:"password" db:"password"`
	DateDeleted   *time.Time `json:"date_deleted" db:"date_deleted"`
	CountDownload *int       `json:"count_download" db:"count_download"`
	Description   string     `json:"description" db:"description"`
}

type DataOutput struct {
	Password      bool      `json:"password" db:"password"`
	DateDeleted   time.Time `json:"date_deleted" db:"date_deleted"`
	CountDownload int       `json:"count_download" db:"count_download"`
	Description   string    `json:"description" db:"description"`
}

type FileGetInput struct {
	Password string `json:"password"`
}

type FileGet struct {
	ID       string
	Password string
}

type FileUpdate struct {
	CountDayToDeleted *int    `json:"count_day_to_deleted"`
	Password          *string `json:"password"`
	CountDownload     *int    `json:"count_download"`
	Description       string  `json:"description" db:"description"`
}

type DayDeletedUpdateInput struct {
	CountDayToDeleted int `json:"count_day_to_deleted"`
}

type PasswordUpdateInput struct {
	Password string `json:"password"`
}

type CountDownloadUpdateInput struct {
	CountDownload int `json:"count_download"`
}

type DescriptionUpdateInput struct {
	Description string `json:"description" db:"description"`
}

type FileInfo struct {
	ID     int    `json:"id"`
	FileID string `json:"fileID"`
}
