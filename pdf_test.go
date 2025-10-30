package main

import (
	"os"
	"testing"
)

func TestExtractText(t *testing.T) {
	t.Run("TestExtractText non-existent file", func(t *testing.T) {
		// Attempt to extract text from a file that does not exist.
		_, err := ExtractText("non-existent-file.pdf")
		if err == nil {
			t.Fatal("Expected an error when opening a non-existent file, but got none")
		}
		// Check if the error is a path error, which is expected for a missing file.
		if !os.IsNotExist(err) {
			t.Fatalf("Expected a file not exist error, but got a different error: %v", err)
		}
	})
}
