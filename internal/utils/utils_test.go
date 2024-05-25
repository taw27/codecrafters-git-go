package utils

import (
	"os"
	"testing"
)

func TestGetHashFromFile_HappyPath(t *testing.T) {
	file, err := os.CreateTemp("", "test")

	if err != nil {
		t.Fatalf("%v", err)
	}

	defer os.Remove(file.Name())

	hash, err := AppUtils{}.GetHashFromFile(file)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(hash) != 40 {
		t.Errorf("Expected hash with length 40, got  %d", len(hash))
	}
}
func TestGetHashFromFile_FileReadError(t *testing.T) {
	file, _ := os.Open("doesnotexist")

	_, err := AppUtils{}.GetHashFromFile(file)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetObjectPathFromFileSha_HappyPath(t *testing.T) {
	fileSha := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	path, _ := AppUtils{}.GetObjectPathFromFileSha(fileSha)

	expectedPath := ".git/objects/a9/4a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	if path != expectedPath {
		t.Errorf("Expected path to be %s, got %s", expectedPath, path)
	}
}

func TestGetObjectPathFromFileSha_InvalidSha(t *testing.T) {
	fileSha := "invalid"

	_, err := AppUtils{}.GetObjectPathFromFileSha(fileSha)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}