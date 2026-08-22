package chunker

import (
	"io"

	rchunker "github.com/restic/chunker"
)

const (
	MinSize = 512 * 1024
	MaxSize = 8 * 1024 * 1024
)

// chunk is one-content defined piece of larger stream
type Chunk struct {
	Data   []byte
	Offset uint
	Length uint
}

// generates a random irreducible polynomial for use as the rolling hash's chunking parameter
func NewPolynomial() (rchunker.Pol, error) {
	return rchunker.RandomPolynomial()
}

// Split streams r through the chunker and invokes cb once per chunk
func Split(r io.Reader, poly rchunker.Pol, cb func(Chunk) error) error {
	c := rchunker.New(r, poly)
	buf := make([]byte, MaxSize)

	for {
		chunk, err := c.Next(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		data := make([]byte, len(chunk.Data))
		copy(data, chunk.Data)

		if err := cb(Chunk{Data: data, Offset: chunk.Start, Length: chunk.Length}); err != nil {
			return err
		}
	}
}
