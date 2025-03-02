package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// Encrypt encrypts a plaintext string using AES-GCM.
func Encrypt(plaintext string, keyStr string) (string, string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return "", "", fmt.Errorf("failed to decode key: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", "", fmt.Errorf("failed to create cipher: %w", err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", "", fmt.Errorf("failed to create GCM: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), base64.StdEncoding.EncodeToString(nonce), nil
}

// Decrypt decrypts a ciphertext string using AES-GCM.
func Decrypt(ciphertextB64, nonceB64 string, keyStr string) (string, error) {

	keyBytes, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode key: %w", err)
	}

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		return "", fmt.Errorf("failed to decode nonce: %w", err)
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}
