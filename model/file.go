package model

import (
	"time"
)

type File struct {
	Id               string    `json:id`
	DateCreated      time.Time `json:date_created`
	DateDeleted      time.Time `json:date_deleted`
	CountDownload    uint8     `json:count_download`
	CountDiscoveries uint8     `json:count_discoveries`
	CountDay         uint8     `json:count_day`
}

type FileSave struct {
	File             string    `json:file`
	DateCreated      time.Time `json:date_created`
	DateDeleted      *uint8    `json:date_deleted`
	CountDownload    *uint8    `json:count_download`
	CountDiscoveries *uint8    `json:count_discoveries`
	CountDay         *uint8    `json:count_day`
}
