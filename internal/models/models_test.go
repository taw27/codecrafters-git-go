package models

import (
	"compress/zlib"
	"io"
	"os"
	"testing"
)

func TestGitObject_GenerateSha(t *testing.T) {
	obj := GitObject{
		Type:    "blob",
		Size:    4,
		Content: "test",
	}

	sha := obj.GenerateSha()

	expected := "30d74d258442c7c65512eafab474568dd706c430"
	if sha != expected {
		t.Errorf("Expected SHA to be %s, got %s", expected, sha)
	}
}

func TestGitObject_GetObjectFileInput(t *testing.T) {
	obj := GitObject{
		Type:    "blob",
		Size:    4,
		Content: "test",
	}

	input := obj.GetObjectFileInput()

	expected := "blob 4\x00test"
	if input != expected {
		t.Errorf("Expected object file input to be '%s', got %s", expected, input)
	}
}

func TestGitObject_CreateBlob_ErrorMakingDirectory(t *testing.T) {
	obj := GitObject{
		Type:    "blob",
		Size:    4,
		Content: "test",
	}

	err := obj.CreateBlob("/invalid/path")

	if err == nil {
		t.Error("Expected error for invalid path, got nil")
	}
}

func TestGitObject_CreateBlob_Success(t *testing.T) {
	obj := GitObject{
		Type:    "blob",
		Size:    4,
		Content: "test",
	}

	filePath := "/tmp/test"
	err := obj.CreateBlob(filePath)

	if err != nil {
		t.Errorf("Expected no error for valid path, got %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("Expected file to exist at %s", filePath)
	}

	// Open and read the file
	file, err := os.Open(filePath)

	if err != nil {
		t.Errorf("Expected to open file, got error: %v", err)
	}
	defer file.Close()

	// Create a zlib reader on the file
	zr, err := zlib.NewReader(file)
	if err != nil {
		t.Errorf("Expected to create zlib reader, got error: %v", err)
	}
	defer zr.Close()

	decompressed, err := io.ReadAll(zr)

	if err != nil {
		t.Errorf("Expected to read decompressed contents, got error: %v", err)
	}

	// Check the decompressed contents
	expected := "blob 4\x00test"
	got := string(decompressed)

	if got != expected {
		t.Errorf("Expected decompressed contents to be '%s', got '%s'", expected, got)
	}

	// remove file for cleanup
	err = os.Remove(filePath)

	if err != nil {
		t.Errorf("Expected to cleanup test files, got error %v", err)
	}
}
