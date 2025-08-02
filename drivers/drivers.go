package drivers

import (
	"io"
)

type Driver interface {
	Read(path string) (io.ReadCloser, error)
	Write(path string, data []byte) error
	Put(path string, data []byte) error
	Delete(path string) error
	Update(path string, data []byte, prepend bool) error
}
