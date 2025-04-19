package model

import (
	"github.com/minio/minio-go/v7"
	"time"
)

type File struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"date_created" db:"date_created"`
	Name        string    `json:"name" db:"name"`
	Data
}

type FileSave struct {
	FileBase64 []string `json:"file"`
	Name       string
	Data
}

type GetFileOutput struct {
	File *minio.Object `json:"file"`
	Name string        `json:"name"`
}

type Data struct {
	DateDeleted      *uint8 `json:"date_deleted" db:"date_deleted"`
	CountDownload    *uint8 `json:"count_download" db:"count_download"`
	CountDiscoveries *uint8 `json:"count_discoveries" db:"count_discoveries"`
	CountDay         *uint8 `json:"count_day" db:"count_day"`
}
