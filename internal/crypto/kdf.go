package crypto

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/argon2"
)

// argon2 parameters for key derivation
type KDFParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	Keylen      uint32
}

// provding standard parameters for KDF
func DefaultKDFParams() KDFParams {
	return KDFParams{
		Memory:      64 * 1024, // 64 MB
		Iterations:  3,
		Parallelism: 2,
		Keylen:      32,
	}
}

// derives a 32-byte master/encryption key from a password and salt using Argon2id
func DeriveKey(password, salt []byte, params KDFParams) []byte {
	return argon2.IDKey(password, salt, params.Iterations, params.Memory, params.Parallelism, params.Keylen)
}

// creates a secure random salt of the requested length
func GenerateSalt(length int) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}
	return salt, nil
}
