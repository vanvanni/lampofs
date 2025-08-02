package errors

import "errors"

var (
	ErrFileNotFound     = errors.New("file not found")
	ErrFileExists       = errors.New("file already exists")
	ErrPermissionDenied = errors.New("permission denied")
)
