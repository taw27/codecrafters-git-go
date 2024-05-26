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

	utils := AppUtils{}

	hash, err := utils.GetShaFromFile(file)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(hash) != 40 {
		t.Errorf("Expected hash with length 40, got  %d", len(hash))
	}
}
func TestGetHashFromFile_FileReadError(t *testing.T) {
	file, _ := os.Open("doesnotexist")

	utils := AppUtils{}
	_, err := utils.GetShaFromFile(file)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetObjectPathFromFileSha_HappyPath(t *testing.T) {
	fileSha := "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	utils := AppUtils{}
	path, _ := utils.GetObjectPathFromFileSha(fileSha)

	expectedPath := ".git/objects/a9/4a8fe5ccb19ba61c4c0873d391e987982fbbd3"

	if path != expectedPath {
		t.Errorf("Expected path to be %s, got %s", expectedPath, path)
	}
}

func TestGetObjectPathFromFileSha_InvalidSha(t *testing.T) {
	fileSha := "invalid"

	utils := AppUtils{}
	_, err := utils.GetObjectPathFromFileSha(fileSha)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestGetShaFromString_EmptyString(t *testing.T) {
	utils := AppUtils{}
	sha := utils.GetShaFromString("")

	if sha != "da39a3ee5e6b4b0d3255bfef95601890afd80709" {
		t.Errorf("Expected SHA for empty string to be da39a3ee5e6b4b0d3255bfef95601890afd80709, got %s", sha)
	}
}

func TestGetShaFromString_NonEmptyString(t *testing.T) {
	utils := AppUtils{}
	sha := utils.GetShaFromString("hello")

	if sha != "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" {
		t.Errorf("Expected SHA for 'hello' to be aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d, got %s", sha)
	}
}
