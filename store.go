package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"sort"

	badger "github.com/dgraph-io/badger/v4"
)

// Document represents a text chunk and its corresponding embedding.
type Document struct {
	Text      string    `json:"text"`
	Source    string    `json:"source"`
	Embedding []float64 `json:"embedding"`
}

// VectorStore holds a collection of documents with BadgerDB persistence.
type VectorStore struct {
	db   *badger.DB
	Docs []Document // In-memory cache for fast search
}

// NewVectorStore initializes a new VectorStore with BadgerDB.
func NewVectorStore(path string) (*VectorStore, error) {
	opts := badger.DefaultOptions(path).WithLoggingLevel(badger.ERROR)
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger db: %w", err)
	}

	vs := &VectorStore{
		db:   db,
		Docs: []Document{},
	}

	// Load existing documents into memory
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var doc Document
				if err := json.Unmarshal(v, &doc); err != nil {
					return err
				}
				vs.Docs = append(vs.Docs, doc)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load documents from db: %w", err)
	}

	return vs, nil
}

// Add appends a new document to the vector store and persists it.
func (vs *VectorStore) Add(doc Document) error {
	vs.Docs = append(vs.Docs, doc)

	return vs.db.Update(func(txn *badger.Txn) error {
		data, err := json.Marshal(doc)
		if err != nil {
			return err
		}
		// Use SHA256 hash of text and source as key
		h := sha256.New()
		h.Write([]byte(doc.Source + doc.Text))
		key := fmt.Sprintf("%x", h.Sum(nil))
		return txn.Set([]byte(key), data)
	})
}

// Close closes the underlying database.
func (vs *VectorStore) Close() error {
	return vs.db.Close()
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
