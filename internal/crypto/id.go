package crypto

import (
	"crypto/sha256"
	"encoding/hex"
)

type ID [32]byte //32-byte SHA256 hash

func (id ID) String() string {
	return hex.EncodeToString(id[:])
}

func ContentID(data []byte) ID {
	return sha256.Sum256(data)
}
