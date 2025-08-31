package domain

import "errors"

var (
	ErrFileDeleted   = errors.New("file deleted")
	ErrDuplicateFile = errors.New("file duplicate")
)
