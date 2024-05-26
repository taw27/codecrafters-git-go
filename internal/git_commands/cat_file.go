package git_commands

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"github.com/codecrafters-io/git-starter-go/internal/models"
	"io"
	"os"
	"strconv"
)

type catFileUtils interface {
	GetObjectPathFromFileSha(fileSha string) (string, error)
}

func CatFile(fileSha string, flag string, u catFileUtils) error {
	pathToObject, err := u.GetObjectPathFromFileSha(fileSha)

	if err != nil {
		return errors.New(fmt.Sprintf("Error: %s\n", err))
	}

	file, err := os.Open(pathToObject)

	if err != nil {
		return errors.New(fmt.Sprintf("Error opening file: %s\n", err))

	}

	defer file.Close()

	reader, err := zlib.NewReader(file)

	if err != nil {
		return errors.New(fmt.Sprintf("Error creating zlib reader: %s\n", err))
	}

	defer reader.Close()

	s, err := io.ReadAll(reader)

	if err != nil {
		return errors.New(fmt.Sprintf("Error reading reader: %s\n", err))
	}

	parts := bytes.Split(s, []byte("\x00"))

	if len(parts) < 2 {
		return errors.New(fmt.Sprintf("Error: not a git object\n"))
	}

	content := string(parts[1])

	parts = bytes.Split(parts[0], []byte(" "))

	if len(parts) != 2 {
		return errors.New(fmt.Sprintf("Error: not a git object\n"))
	}

	objectType := string(parts[0])
	sizeStr := string(parts[1])

	size, err := strconv.Atoi(string(sizeStr))

	if err != nil {
		return errors.New(fmt.Sprintf("Error: not a git object\n"))
	}

	gitObject := models.GitObject{
		Type:    objectType,
		Size:    size,
		Content: content,
	}

	if flag != "-p" {
		return errors.New(fmt.Sprintf("Error: command '%s' not recognized", flag))
	}

	gitObject.PrettyPrintContent()

	return nil
}
