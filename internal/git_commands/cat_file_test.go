package git_commands

import (
	"bytes"
	"compress/zlib"
	"github.com/codecrafters-io/git-starter-go/internal/utils"
	"io"
	"os"
	"strings"
	"testing"
)

type mockUtils struct {
	tmpFilePath string
}

const dummyTestFileName = "test_object"

func (m mockUtils) GetObjectPathFromFileSha(_ string) (string, error) {
	return m.tmpFilePath, nil
}

func TestCatFile_HappyPath_ContentOutput(t *testing.T) {

	tmpFile, err := os.CreateTemp("", dummyTestFileName)

	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpFile.Name())

	content := []byte("blob 11\x00Hello, world!")

	var b bytes.Buffer
	zlibW := zlib.NewWriter(&b)

	if _, err := zlibW.Write(content); err != nil {
		t.Fatal(err)
	}

	if err := zlibW.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpFile.Write(b.Bytes()); err != nil {
		t.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	// Redirect standard output to a buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err = CatFile(tmpFile.Name(), "-p", mockUtils{tmpFilePath: tmpFile.Name()})

	// Restore standard output
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}

	os.Stdout = oldStdout

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Read the content from the buffer
	out, _ := io.ReadAll(r)
	output := strings.TrimSuffix(string(out), "\n")

	// Check if the output is as expected
	expectedOutput := "Hello, world!"
	if output != expectedOutput {
		t.Errorf("Expected output to be %s, got %s", expectedOutput, output)
	}
}

func TestCatFile_FileOpenError(t *testing.T) {
	err := CatFile(dummyTestFileName, "p", utils.AppUtils{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCatFile_InvalidGitObject(t *testing.T) {
	tmpFile, err := os.CreateTemp("", dummyTestFileName)

	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("invalid")

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	err = CatFile(tmpFile.Name(), "p", mockUtils{tmpFilePath: tmpFile.Name()})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestCatFile_InvalidFlag(t *testing.T) {
	tmpfile, err := os.CreateTemp("", dummyTestFileName)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := []byte("blob 11\x00Hello, world!")

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	err = CatFile(tmpfile.Name(), "invalid", mockUtils{})
	if err == nil {
		t.Error("Expected error, got nil")
	}
}
