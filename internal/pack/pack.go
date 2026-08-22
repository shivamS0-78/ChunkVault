package pack

import (
	"bytes"
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

func (w *Writer) Add(id string, encrypted []Data) {
	w.entries = append(w.entries, Entry{
		ID:     id,
		Offset: uint64(w.buf.Len()),
		Length: uint32(len(encrypted)),
	})
	w.buf.Write(encrypted)
}
