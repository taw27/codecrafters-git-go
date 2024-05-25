package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type Utils interface {
	GetObjectPathFromFileSha(fileSha string) (string, error)
	GetHashFromFile(reader io.Reader) (string, error)
}

type AppUtils struct {
}

func (a AppUtils) GetObjectPathFromFileSha(fileSha string) (string, error) {
	if len(fileSha) != 40 {
		return "", errors.New("file sha is invalid. Needs to be 40 char")
	}

	baseDir := ".git/objects/"
	shaFirstTwoChars := fileSha[0:2]
	fileName := fileSha[2:]

	objectPath := fmt.Sprintf("%s%s/%s", baseDir, shaFirstTwoChars, fileName)

	return objectPath, nil
}

func (a AppUtils) GetHashFromFile(reader io.Reader) (string, error) {
	hash := sha1.New()

	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)[:20]

	sha := hex.EncodeToString(hashBytes)
	return sha, nil
}
