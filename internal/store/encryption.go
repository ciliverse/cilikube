package store

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// !!! IMPORTANT SECURITY WARNING !!!
// In production environments, the encryption key (encryptionKey) must never be hardcoded in the code.
// It must be loaded from a secure location, such as environment variables or Secrets Manager.
// This key must be 32 bytes long, corresponding to AES-256.

func Encrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes long for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not create new GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt(data []byte, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes long for AES-256")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("could not create new cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("could not create new GCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt data: %w", err)
	}
	return plaintext, nil
}
