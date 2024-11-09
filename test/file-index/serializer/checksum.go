package serializer

import (
	"fmt"
	"hash"
	"io"
	"os"
)

func CalculateHash(filePath string, hashFunc func() hash.Hash) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	h := hashFunc()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	hashBytes := h.Sum(nil)
	return fmt.Sprintf("%x", hashBytes), nil
}
