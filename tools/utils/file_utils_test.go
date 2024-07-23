package utils

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestCopyFile(t *testing.T) {
	// Setup - create a temporary source file
	srcContent := []byte("Hello, world!")
	srcFile, err := os.CreateTemp("", "src")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(srcFile.Name())
	if _, err := srcFile.Write(srcContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	srcFile.Close()

	t.Run("Successful copy with forceOverwrite true", func(t *testing.T) {
		destFile := filepath.Join(os.TempDir(), "destFile")
		defer os.Remove(destFile)

		fu := FileUtils{}
		if err := fu.CopyFile(srcFile.Name(), destFile, true); err != nil {
			t.Errorf("CopyFile returned an error: %v", err)
		}

		verifyFileContent(t, destFile, srcContent)
	})

	t.Run("No copy when forceOverwrite is false and destination exists", func(t *testing.T) {
		destFile := filepath.Join(os.TempDir(), "destFileExists")
		err := os.WriteFile(destFile, []byte("Existing content"), 0644)
		if err != nil {
			return
		}
		defer os.Remove(destFile)

		fu := FileUtils{}
		if err := fu.CopyFile(srcFile.Name(), destFile, false); err != nil {
			t.Errorf("CopyFile returned an error: %v", err)
		}

		content, _ := os.ReadFile(destFile)
		if string(content) == string(srcContent) {
			t.Errorf("File should not have been overwritten")
		}
	})

	t.Run("Error when source file does not exist", func(t *testing.T) {
		destFile := filepath.Join(os.TempDir(), "dest")
		defer os.Remove(destFile)

		fu := FileUtils{}
		if err := fu.CopyFile("nonexistent", destFile, true); err == nil {
			t.Errorf("Expected an error when source file does not exist")
		}
	})

	t.Run("Error when unable to create destination directory", func(t *testing.T) {
		destFile := "/root/destFile" // Assuming /root is not writable for the current user

		fu := FileUtils{}
		if err := fu.CopyFile(srcFile.Name(), destFile, true); err == nil {
			t.Errorf("Expected an error when unable to create destination directory")
		}
	})
}

func verifyFileContent(t *testing.T, filePath string, expectedContent []byte) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != string(expectedContent) {
		t.Errorf("Content mismatch: got %v, want %v", string(content), string(expectedContent))
	}
}

func TestWriteToFile(t *testing.T) {
	// Setup test environment
	tempDir := os.TempDir()

	mockData := "test data"
	mockReader := strings.NewReader(mockData)
	filePath := filepath.Join(tempDir, "testfile")

	// Capture log output
	var logOutput bytes.Buffer
	log.SetOutput(&logOutput)

	fu := FileUtils{}

	// Test WriteToFile success
	err := fu.WriteToFile(filePath, mockReader, 200)
	if err != nil {
		t.Errorf("WriteToFile failed: %v", err)
	}

	// Verify file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	if string(content) != mockData {
		t.Errorf("File content mismatch: got %v, want %v", string(content), mockData)
	}

	// Verify log output
	expectedLogPattern := `Successfully written 9 bytes to ` + regexp.QuoteMeta(filePath) + `\. HTTP status code: 200\n`
	matched, err := regexp.MatchString(expectedLogPattern, logOutput.String())
	if err != nil {
		t.Fatalf("Regex match error: %v", err)
	}
	if !matched {
		t.Errorf("Log output mismatch: expected to match %v", expectedLogPattern)
	}

	// Test WriteToFile failure with invalid path
	invalidPath := "/invalid/path"
	err = fu.WriteToFile(invalidPath, mockReader, 400)
	if err == nil {
		t.Errorf("Expected WriteToFile to fail with invalid path")
	}
}
