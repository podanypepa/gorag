package main

import (
	"fmt"
	"testing"
)

func TestExtractText(t *testing.T) {
	t.Run("TestExtractText", func(t *testing.T) {
		text, err := ExtractText("./invoice-MSTRL-API-664426-002.pdf")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		fmt.Println("Extracted Text:", text)
	})
}
