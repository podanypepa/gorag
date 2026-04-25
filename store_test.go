package main

import (
	"os"
	"testing"
)

func TestVectorStore(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "badger_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Add and Persistence", func(t *testing.T) {
		store, err := NewVectorStore(tmpDir)
		if err != nil {
			t.Fatalf("failed to create store: %v", err)
		}

		doc1 := Document{Text: "Hello world", Source: "test.md", Embedding: []float64{1, 2, 3}}
		doc2 := Document{Text: "Hello gopher", Source: "test.md", Embedding: []float64{4, 5, 6}}
		
		if err := store.Add(doc1); err != nil {
			t.Fatalf("failed to add doc1: %v", err)
		}
		if err := store.Add(doc2); err != nil {
			t.Fatalf("failed to add doc2: %v", err)
		}

		if len(store.Docs) != 2 {
			t.Fatalf("expected 2 documents, got %d", len(store.Docs))
		}
		store.Close()

		// Re-open and check persistence
		store2, err := NewVectorStore(tmpDir)
		if err != nil {
			t.Fatalf("failed to re-open store: %v", err)
		}
		defer store2.Close()

		if len(store2.Docs) != 2 {
			t.Fatalf("expected 2 documents after reload, got %d", len(store2.Docs))
		}
	})

	t.Run("Search", func(t *testing.T) {
		searchDir, _ := os.MkdirTemp("", "badger_search_*")
		defer os.RemoveAll(searchDir)
		
		store, _ := NewVectorStore(searchDir)
		defer store.Close()

		doc1 := Document{Text: "cat", Embedding: []float64{1, 0, 0}}
		doc2 := Document{Text: "dog", Embedding: []float64{0, 1, 0}}
		doc3 := Document{Text: "fish", Embedding: []float64{0, 0, 1}}
		doc4 := Document{Text: "kitten", Embedding: []float64{0.9, 0.1, 0}}

		store.Add(doc1)
		store.Add(doc2)
		store.Add(doc3)
		store.Add(doc4)

		query := []float64{1, 0, 0}
		results := store.Search(query, 2)

		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		if results[0].Text != "cat" {
			t.Errorf("expected first result to be 'cat', got '%s'", results[0].Text)
		}
		// DeepEqual might fail because ordering of reload depends on Badger iteration
		// But in Search, it should be ordered by score
		if results[1].Text != "kitten" {
			t.Errorf("expected second result to be 'kitten', got '%s'", results[1].Text)
		}
	})
}
