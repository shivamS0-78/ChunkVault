package crypto_test

import (
	"bytes"
	"testing"

	"restic-clone/internal/crypto"
)

func TestContentID_Determinism(t *testing.T) {
	data := []byte("identical chunk payload")
	id1 := crypto.ContentID(data)
	id2 := crypto.ContentID(data)

	if id1 != id2 {
		t.Fatalf("expected identical ContentIDs, got %s and %s", id1, id2)
	}
}

func TestDeriveKey_Determinism(t *testing.T) {
	password := []byte("master-passphrase")
	salt := []byte("1234567890123456")
	params := crypto.DefaultKDFParams()

	key1 := crypto.DeriveKey(password, salt, params)
	key2 := crypto.DeriveKey(password, salt, params)

	if !bytes.Equal(key1, key2) {
		t.Fatal("DeriveKey must yield identical keys for identical inputs")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := make([]byte, 32) // 256-bit key
	plaintext := []byte("payload data to encrypt and decrypt")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecrypt_DetectsTampering(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("sensitive chunk content")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Bit-flip the final byte in the AEAD auth tag
	ciphertext[len(ciphertext)-1] ^= 0x01

	_, err = crypto.Decrypt(key, ciphertext)
	if err == nil {
		t.Fatal("expected decryption error on corrupted ciphertext, got nil")
	}
}
