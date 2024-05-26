package models

import (
	"bytes"
	"compress/zlib"
	"errors"
	"fmt"
	"github.com/codecrafters-io/git-starter-go/internal/utils"
	"os"
	"path/filepath"
)

type GitObject struct {
	Type    string
	Size    int
	Content string
}

func (g *GitObject) PrettyPrintContent() {
	fmt.Printf(g.Content)
}

func (g *GitObject) GenerateSha() string {
	appUtils := utils.AppUtils{}

	return appUtils.GetShaFromString(g.GetObjectFileInput())
}

func (g *GitObject) GetObjectFileInput() string {
	return fmt.Sprintf("%s %d\x00%s", g.Type, g.Size, g.Content)
}

func (g *GitObject) CreateBlob(filePath string) error {

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("Error making directory for object: %v", err))
	}

	file, err := os.Create(filePath)

	if err != nil {
		return errors.New(fmt.Sprintf("Error opening object file: %v", err))
	}

	defer file.Close()

	fileContents := []byte(fmt.Sprintf("%s %d\x00%s", g.Type, g.Size, g.Content))

	var compressed bytes.Buffer
	zlibWriter := zlib.NewWriter(&compressed)

	if _, err := zlibWriter.Write(fileContents); err != nil {
		return errors.New(fmt.Sprintf("Error compressing object: %v", err))
	}
	if err := zlibWriter.Close(); err != nil {
		return errors.New(fmt.Sprintf("Error closing zlib writer: %v", err))
	}

	compressedData := compressed.Bytes()

	if _, err = file.Write(compressedData); err != nil {
		return errors.New(fmt.Sprintf("Error writing to object file: %v", err))
	}

	return nil
}
