package pack

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

const TargetSize = 24 * 1024 * 1024

type Entry struct {
	ID     string `json:"id"`
	Offset uint64 `json:"offset"`
	Length uint32 `json:"length"`
}

type Writer struct {
	entries []Entry
	buf     bytes.Buffer
}

func NewWriter() *Writer {
	return &Writer{}
}

func (w *Writer) Add(id string, encrypted []byte) {
	w.entries = append(w.entries, Entry{
		ID:     id,
		Offset: uint64(w.buf.Len()),
		Length: uint32(len(encrypted)),
	})
	w.buf.Write(encrypted)
}

func (w *Writer) Size() int {
	return w.buf.Len()
}

func (w *Writer) Len() int {
	return len(w.entries)
}

// Finalize returns the complete pack bytes and freshly generated random pack ID
func (w *Writer) Finalize() (packID string, data []byte, entries []Entry, err error) {
	if len(w.entries) == 0 {
		return "", nil, nil, fmt.Errorf("pack : finalize called with no chunks")
	}
	footer, err := json.Marshal(w.entries)
	if err != nil {
		return "", nil, nil, fmt.Errorf("pack : finalize : marshal footer: %w", err)
	}

	out := make([]byte, 0, w.buf.Len()+len(footer)+4)
	out = append(out, w.buf.Bytes()...)
	out = append(out, footer...)

	var lenBuf [4]byte
	binary.LittleEndian.PutUint32(lenBuf[:], uint32(len(footer)))
	out = append(out, lenBuf[:]...)

	id, err := randomID()
	if err != nil {
		return "", nil, nil, fmt.Errorf("pack : generate pack id: %w", err)
	}

	return id, out, w.entries, nil

}

func randomID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
