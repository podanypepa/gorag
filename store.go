package main

import (
	"encoding/json"
	"math"
	"os"
	"sort"
)

type Document struct {
	Text      string
	Embedding []float64
}

type VectorStore struct {
	Docs []Document
}

func (vs *VectorStore) Add(doc Document) {
	vs.Docs = append(vs.Docs, doc)
}

func (vs *VectorStore) Save(path string) error {
	data, _ := json.MarshalIndent(vs, "", "  ")
	return os.WriteFile(path, data, 0644)
}

func LoadStore(path string) (*VectorStore, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return &VectorStore{}, nil
	}
	var vs VectorStore
	json.Unmarshal(data, &vs)
	return &vs, nil
}

func cosine(a, b []float64) float64 {
	var dot, na, nb float64
	for i := range a {
		dot += a[i] * b[i]
		na += a[i] * a[i]
		nb += b[i] * b[i]
	}
	return dot / (math.Sqrt(na) * math.Sqrt(nb))
}

func (vs *VectorStore) Search(query []float64, k int) []Document {
	type scored struct {
		doc   Document
		score float64
	}
	var scores []scored
	for _, d := range vs.Docs {
		score := cosine(d.Embedding, query)
		scores = append(scores, scored{d, score})
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	var top []Document
	for i := 0; i < k && i < len(scores); i++ {
		top = append(top, scores[i].doc)
	}
	return top
}
