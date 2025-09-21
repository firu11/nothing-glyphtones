package utils

import (
	"io"
	"os"

	"github.com/zeebo/blake3"
)

func GetCheckSum(file *os.File) ([]byte, error) {
	var hasher *blake3.Hasher = blake3.New()

	_, err := io.Copy(hasher, file)
	if err != nil {
		return nil, err
	}

	return hasher.Sum(nil), nil
}
