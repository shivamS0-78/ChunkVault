package local

import (
	"errors"
	"fmt"
	"os"
)

var errAlready = errors.New("already exists")

type Local struct {
	root string
}

func New(dir string) (*Local, error) {
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("local: create root: %w", err)
	}
	return &Local{root: dir}, nil
}
