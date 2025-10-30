package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestVectorStore(t *testing.T) {
	t.Run("Add and Save", func(t *testing.T) {
		store := &VectorStore{}
		doc1 := Document{Text: "Hello world", Embedding: []float64{1, 2, 3}}
		doc2 := Document{Text: "Hello gopher", Embedding: []float64{4, 5, 6}}
		store.Add(doc1)
		store.Add(doc2)

		if len(store.Docs) != 2 {
			t.Fatalf("expected 2 documents, got %d", len(store.Docs))
		}

		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test_store.json")

		err := store.Save(testFile)
		if err != nil {
			t.Fatalf("failed to save store: %v", err)
		}

		if _, err := os.Stat(testFile); os.IsNotExist(err) {
			t.Fatalf("saved file does not exist")
		}
	})

	t.Run("LoadStore", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test_store.json")

		// Test loading non-existent file
		store, err := LoadStore(testFile)
		if err != nil {
			t.Fatalf("expected no error for non-existent file, got %v", err)
		}
		if len(store.Docs) != 0 {
			t.Fatalf("expected 0 documents for new store, got %d", len(store.Docs))
		}

		// Test loading existing file
		store = &VectorStore{}
		doc1 := Document{Text: "Hello world", Embedding: []float64{1, 2, 3}}
		store.Add(doc1)
		err = store.Save(testFile)
		if err != nil {
			t.Fatalf("failed to save store for loading test: %v", err)
		}

		loadedStore, err := LoadStore(testFile)
		if err != nil {
			t.Fatalf("failed to load store: %v", err)
		}
		if len(loadedStore.Docs) != 1 {
			t.Fatalf("expected 1 document, got %d", len(loadedStore.Docs))
		}
		if !reflect.DeepEqual(store.Docs[0], loadedStore.Docs[0]) {
			t.Fatalf("loaded document does not match saved document")
		}
	})

	t.Run("Search", func(t *testing.T) {
		store := &VectorStore{}
		doc1 := Document{Text: "cat", Embedding: []float64{1, 0, 0}}
		doc2 := Document{Text: "dog", Embedding: []float64{0, 1, 0}}
		doc3 := Document{Text: "fish", Embedding: []float64{0, 0, 1}}
		doc4 := Document{Text: "kitten", Embedding: []float64{0.9, 0.1, 0}} // very similar to cat

		store.Add(doc1)
		store.Add(doc2)
		store.Add(doc3)
		store.Add(doc4)

		query := []float64{1, 0, 0} // "cat"
		results := store.Search(query, 2)

		if len(results) != 2 {
			t.Fatalf("expected 2 results, got %d", len(results))
		}
		if results[0].Text != "cat" {
			t.Errorf("expected first result to be 'cat', got '%s'", results[0].Text)
		}
		if results[1].Text != "kitten" {
			t.Errorf("expected second result to be 'kitten', got '%s'", results[1].Text)
		}
	})
}
