package pii

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type SecretProvider[K any] interface {
	SecretForKey(key K) ([]byte, error)
}

type AESAnonymizer[K any] struct {
	secretProvider SecretProvider[K]
}

func NewAESAnonymizer[K any](secretProvider SecretProvider[K]) AESAnonymizer[K] {
	return AESAnonymizer[K]{
		secretProvider: secretProvider,
	}
}

func (a AESAnonymizer[K]) AnonymizeString(key K, value string) (string, error) {
	secret, err := a.secretProvider.SecretForKey(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(value), nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func (a AESAnonymizer[K]) DeanonymizeString(key K, value string) (string, error) {
	secret, err := a.secretProvider.SecretForKey(key)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decoded, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(decoded) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := decoded[:nonceSize], decoded[nonceSize:]
	data, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
