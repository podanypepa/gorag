package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
)

// Document represents a text chunk and its corresponding embedding.
type Document struct {
	Text      string    `json:"text"`
	Embedding []float64 `json:"embedding"`
}

// VectorStore holds a collection of documents in memory.
// For larger datasets, consider using a dedicated vector database
// like ChromaDB, Weaviate, or a library like Faiss.
type VectorStore struct {
	Docs []Document `json:"docs"`
}

// Add appends a new document to the vector store.
func (vs *VectorStore) Add(doc Document) {
	vs.Docs = append(vs.Docs, doc)
}

// Save serializes the vector store to a JSON file.
func (vs *VectorStore) Save(path string) error {
	data, err := json.MarshalIndent(vs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal vector store: %w", err)
	}
	return os.WriteFile(path, data, 0644)
}

// LoadStore deserializes the vector store from a JSON file.
// If the file does not exist, it returns a new, empty VectorStore.
func LoadStore(path string) (*VectorStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &VectorStore{}, nil
		}
		return nil, fmt.Errorf("failed to read vector store file: %w", err)
	}

	var vs VectorStore
	if err := json.Unmarshal(data, &vs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vector store: %w", err)
	}
	return &vs, nil
}

// cosine calculates the cosine similarity between two vectors.
func cosine(a, b []float64) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	if na == 0 || nb == 0 {
		return 0.0
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

// Search finds the top k documents most similar to the query vector.
func (vs *VectorStore) Search(query []float64, k int) []Document {
	type scored struct {
		doc   Document
		score float64
	}

	scores := make([]scored, len(vs.Docs))
	for i, d := range vs.Docs {
		score := cosine(d.Embedding, query)
		scores[i] = scored{d, score}
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	top := make([]Document, 0, k)
	for i := 0; i < k && i < len(scores); i++ {
		top = append(top, scores[i].doc)
	}
	return top
}
