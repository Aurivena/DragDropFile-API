package domain

import (
	"errors"
	"runtime"
)

const PrefixZipFile = "dg-"

var (
	ErrFileDeleted       = errors.New("file deleted")
	PasswordInvalidError = errors.New("password invalid")
	ErrFileDuplicate     = errors.New("file duplicate")
	InternalError        = errors.New("internal error")
	BadRequestError      = errors.New("bad request")
	NotFoundError        = errors.New("not found")
	GoneError            = errors.New("gone error")

	WorkerPool = runtime.GOMAXPROCS(0)
)
