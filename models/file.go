package models

import (
	"github.com/minio/minio-go/v7"
	"time"
)

type File struct {
	FileBase64 string `form:"file"`
	Filename   string `form:"filename"`
}

type FileOutput struct {
	Id   string `json:"id"`
	Name string `json:"name" db:"name"`
}

type FileSave struct {
	Id        string
	Name      string
	SessionID string
	MimeType  string
}

type FilSaveOutput struct {
	Size  int64  `json:"size"`
	Count int    `json:"count"`
	ID    string `json:"id"`
}

type GetFileOutput struct {
	File *minio.Object `json:"file"`
	Name string        `json:"name"`
}

type Data struct {
	Password      *string    `json:"password" db:"password"`
	DateDeleted   *time.Time `json:"date_deleted" db:"date_deleted"`
	CountDownload *int       `json:"count_download" db:"count_download"`
}

type DataOutput struct {
	Password      bool       `json:"password" db:"password"`
	DateDeleted   *time.Time `json:"date_deleted" db:"date_deleted"`
	CountDownload *int       `json:"count_download" db:"count_download"`
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
