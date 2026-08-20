package local

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"restic-clone/internal/backend"
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

var _ backend.Backend = (*Local)(nil)

func (l *Local) path(id string) string {
	if len(id) < 2 {
		return filepath.Join(l.root, id)
	}
	return filepath.Join(l.root, id[:2], id)
}

// writing packed bundles of encrypted chunks into temp file
func (l *Local) Save(id string, r io.Reader) error {
	p := l.path(id)

	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("local : save %v: %w", id, errAlready)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("local : save %v: %w", id, err)
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return fmt.Errorf("local : save %v: create temp: %w", id, err)
	}

	tmp := p + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("local : save %v: create temp: %w", id, err)
	}

	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("local : save %v: write : %w", id, err)
	}

	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("local : save %v: sync : %w", id, err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("local : save %v: close : %w", id, err)
	}

	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("local : save %v: rename : %w", id, err)
	}

	return nil
}
