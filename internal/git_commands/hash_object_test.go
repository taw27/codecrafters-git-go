package git_commands

import (
	"bytes"
	"compress/zlib"
	"errors"
	"io"
	"os"
	"testing"
)

type mockHashObjectUtils struct {
	path string
	err  error
}

func (m mockHashObjectUtils) GetObjectPathFromFileSha(_ string) (string, error) {
	return m.path, m.err
}

func TestHashObject_InvalidFlag(t *testing.T) {
	err := HashObject("dummy", "invalid", mockHashObjectUtils{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestHashObject_FileReadError(t *testing.T) {
	err := HashObject("/invalid/path", "", mockHashObjectUtils{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestHashObject_WriteError(t *testing.T) {
	err := HashObject("hash_object_test.go", "-w", mockHashObjectUtils{err: errors.New("dummy error")})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func createTempFile(content string) (string, error) {
	tempFile, err := os.CreateTemp("", "test")

	if err != nil {
		return "", err
	}

	defer tempFile.Close()

	if _, err := tempFile.WriteString(content); err != nil {
		return "", err
	}

	return tempFile.Name(), nil
}

func TestHashObject_SuccessWithoutWrite(t *testing.T) {
	filePath, err := createTempFile("test content")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(filePath)

	expected := "08cf6101416f0ce0dda3c80e627f333854c4085c"

	// Redirect standard output to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = HashObject(filePath, "", mockHashObjectUtils{path: "/tmp/test"})

	// Restore standard output
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	got := buf.String()

	if got != expected {
		t.Errorf("Expected SHA to be %s, got %s", expected, got)
	}
}

func TestHashObject_SuccessWithWrite(t *testing.T) {
	filePath, err := createTempFile("test content")

	if err != nil {
		t.Fatalf("Failed to create temp gotFile: %v", err)
	}
	defer os.Remove(filePath)

	// Redirect standard output to a buffer
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	mockUtils := mockHashObjectUtils{path: filePath}
	err = HashObject(filePath, "-w", mockUtils)

	// Restore standard output
	w.Close()
	os.Stdout = old

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Read the output
	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Open and read the gotFile
	gotFile, err := os.Open(mockUtils.path)

	if err != nil {
		t.Errorf("Expected to open gotFile, got error: %v", err)
	}
	defer gotFile.Close()

	// Create a zlib reader on the gotFile
	zr, err := zlib.NewReader(gotFile)
	if err != nil {
		t.Errorf("Expected to create zlib reader, got error: %v", err)
	}
	defer zr.Close()

	decompressed, err := io.ReadAll(zr)

	if err != nil {
		t.Errorf("Expected to read decompressed contents, got error: %v", err)
	}

	// Check the decompressed contents
	expected := "blob 12\x00test content"
	got := string(decompressed)

	if got != expected {
		t.Errorf("Expected decompressed contents to be '%s', got '%s'", expected, got)
	}

	// Cleanup files
	err = os.Remove(filePath)

	if err != nil {
		t.Errorf("Expected to cleanup test files, got error %v", err)
	}
}
