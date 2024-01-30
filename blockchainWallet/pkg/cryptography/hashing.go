package cryptography

import (
	"crypto/md5"
	"crypto/sha256"
	"hash"
	"io"
)

func HashMD5(text string) ([]byte, error) {
	return hashWith(md5.New(), text)
}

func HashSHA256(text string) ([]byte, error) {
	return hashWith(sha256.New(), text)
}

func hashWith(provider hash.Hash, text string) ([]byte, error) {
	_, err := io.WriteString(provider, text)
	if err != nil {
		return nil, err
	}

	hash := provider.Sum(nil)
	return hash, nil
}
