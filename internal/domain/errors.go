package domain

import "errors"

var (
	ErrFileDeleted     = errors.New("file deleted")
	ErrPasswordInvalid = errors.New("password invalid")
	ErrDuplicateFile   = errors.New("file duplicate")
)
