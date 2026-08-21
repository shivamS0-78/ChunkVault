package local

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"restic-clone/internal/backend"
)

var _ backend.Backend = (*Local)(nil)

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

func (l *Local) path(id string) string {
	if len(id) < 2 {
		return filepath.Join(l.root, id)
	}
	return filepath.Join(l.root, id[:2], id)
}

func (l *Local) Save(id string, r io.Reader) error {
	p := l.path(id)

	if _, err := os.Stat(p); err == nil {
		return fmt.Errorf("local: save %s: %w", id, errAlready)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("local: save %s: %w", id, err)
	}

	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return fmt.Errorf("local: save %s: create dir: %w", id, err)
	}

	tmp := p + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return fmt.Errorf("local: save %s: create temp: %w", id, err)
	}

	if _, err := io.Copy(f, r); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("local: save %s: write: %w", id, err)
	}

	if err := f.Sync(); err != nil {
		f.Close()
		os.Remove(tmp)
		return fmt.Errorf("local: save %s: sync: %w", id, err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("local: save %s: close: %w", id, err)
	}

	if err := os.Rename(tmp, p); err != nil {
		os.Remove(tmp)
		return fmt.Errorf("local: save %s: rename: %w", id, err)
	}

	return nil
}

func (l *Local) Load(id string) (io.ReadCloser, error) {
	f, err := os.Open(l.path(id))
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("local: load %s: %w", id, backend.ErrNotFound)
		}
		return nil, fmt.Errorf("local: load %s: %w", id, err)
	}
	return f, nil
}

func (l *Local) Stat(id string) (int64, error) {
	fi, err := os.Stat(l.path(id))
	if err != nil {
		if os.IsNotExist(err) {
			return 0, fmt.Errorf("local: stat %s: %w", id, backend.ErrNotFound)
		}
		return 0, fmt.Errorf("local: stat %s: %w", id, err)
	}
	return fi.Size(), nil
}

func (l *Local) List() ([]string, error) {
	var ids []string
	err := filepath.WalkDir(l.root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".tmp") {
			return nil
		}

		ids = append(ids, d.Name())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("local: list: %w", err)
	}
	return ids, nil
}

func (l *Local) Remove(id string) error {
	if err := os.Remove(l.path(id)); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("local: remove %s: %w", id, backend.ErrNotFound)
		}
		return fmt.Errorf("local: remove %s: %w", id, err)
	}
	return nil
}
