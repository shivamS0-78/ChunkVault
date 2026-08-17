package backend

import (
	"errors"
	"io"
)

var ErrNotFound = errors.New("blob not found")

type Backend interface {
	Save(id string, r io.Reader) error

	Load(id string) (io.ReadCloser, error)

	Stat(id string) (int64, error)

	List() ([]string, error)

	Delete(id string) error
}
