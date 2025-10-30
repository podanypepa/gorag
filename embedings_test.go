package main

import (
	"fmt"
	"testing"
)

func TestGetEmbedding(t *testing.T) {
	t.Run("GetEmbedding", func(t *testing.T) {
		text := "Hello, world!"
		embedding, err := GetEmbedding(text)
		if err != nil {
			t.Fatalf("GetEmbedding failed: %v", err)
		}
		fmt.Printf("Embedding: %v\n", embedding)
		if len(embedding) == 0 {
			t.Fatal("Expected non-empty embedding")
		}
	})
}
