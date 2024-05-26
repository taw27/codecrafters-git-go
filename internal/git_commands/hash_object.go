package git_commands

import (
	"errors"
	"fmt"
	"github.com/codecrafters-io/git-starter-go/internal/models"
	"os"
)

type hashObjectUtils interface {
	GetObjectPathFromFileSha(fileSha string) (string, error)
}

func HashObject(filePath, flag string, u hashObjectUtils) error {
	if flag != "" && flag != "-w" {
		return errors.New("Error: provided flag invalid.\nAvailable flags:\n- '-w': Write the object into object storage")
	}

	contentBytes, err := os.ReadFile(filePath)

	if err != nil {
		return errors.New(fmt.Sprintf("Error reading content fiile: %v", err))
	}

	content := string(contentBytes)

	gitObject := models.GitObject{
		Type:    "blob",
		Size:    len(content),
		Content: content,
	}

	objectSha := gitObject.GenerateSha()

	fmt.Printf(objectSha)

	if flag == "" {
		return nil
	}

	objectPath, err := u.GetObjectPathFromFileSha(objectSha)

	if err != nil {
		return errors.New(fmt.Sprintf("Error generating objectPath: %v", err))
	}

	err = gitObject.CreateBlob(objectPath)

	if err != nil {
		return err
	}

	return nil
}
