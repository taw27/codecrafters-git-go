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
	GetShaFromFile(reader io.Reader) (string, error)
	GetShaFromString(reader io.Reader) (string, error)
}

type AppUtils struct {
}

func (a *AppUtils) GetObjectPathFromFileSha(sha string) (string, error) {
	shaRune := []rune(sha)

	if len(shaRune) != 40 {
		return "", errors.New("file sha is invalid. Needs to be 40 char")
	}

	baseDir := ".git/objects/"
	shaFirstTwoChars := shaRune[0:2]
	fileName := shaRune[2:]

	objectPath := fmt.Sprintf("%s%s/%s", baseDir, string(shaFirstTwoChars), string(fileName))

	return objectPath, nil
}

func (a *AppUtils) GetShaFromFile(reader io.Reader) (string, error) {
	hash := sha1.New()

	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}

	hashBytes := hash.Sum(nil)[:20]

	sha := hex.EncodeToString(hashBytes)
	return sha, nil
}

func (a *AppUtils) GetShaFromString(s string) string {
	hash := sha1.New()

	hash.Write([]byte(s))

	hashBytes := hash.Sum(nil)[:20]

	return hex.EncodeToString(hashBytes)
}
