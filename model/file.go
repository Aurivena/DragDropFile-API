package model

import (
	"time"
)

type File struct {
	Id          string    `json:"id"`
	DateCreated time.Time `json:"date_created"`
	Data
}

type FileSave struct {
	FileBase64 string `json:"file"`
	Data
}

type Data struct {
	DateDeleted      *uint8 `json:"date_deleted"`
	CountDownload    *uint8 `json:"count_download"`
	CountDiscoveries *uint8 `json:"count_discoveries"`
	CountDay         *uint8 `json:"count_day"`
}
