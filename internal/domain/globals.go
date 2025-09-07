package domain

import (
	"errors"
	"runtime"
)

const PrefixZipFile = "dg-"

var (
	ErrFileDeleted     = errors.New("file deleted")
	ErrPasswordInvalid = errors.New("password invalid")
	ErrFileDuplicate   = errors.New("file duplicate")

	WorkerPool = runtime.GOMAXPROCS(0)
)
