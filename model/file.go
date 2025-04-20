package model

import (
	"github.com/minio/minio-go/v7"
	"time"
)

type File struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"date_created" db:"date_created"`
	Data
}

type FileSave struct {
	FileBase64 []string `json:"file"`
	Data
}

type GetFileOutput struct {
	File *minio.Object `json:"file"`
	Name string        `json:"name"`
}

type Data struct {
	Name             string  `json:"name" db:"name"`
	Password         *string `json:"password" db:"password"`
	DateDeleted      *uint8  `json:"date_deleted" db:"date_deleted"`
	CountDownload    *uint8  `json:"count_download" db:"count_download"`
	CountDiscoveries *uint8  `json:"count_discoveries" db:"count_discoveries"`
	CountDay         *uint8  `json:"count_day" db:"count_day"`
}
