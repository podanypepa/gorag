package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// EmbeddingRequest represents the request payload for the embedding API.
type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// EmbeddingResponse represents the response payload from the embedding API.
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// GetEmbedding fetches the embedding for the given text from the local API.
func GetEmbedding(text string) ([]float64, error) {
	req := EmbeddingRequest{Model: "llama3", Prompt: text}
	data, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewReader(data))
	if err != nil {
		slog.Error("Failed to get embedding", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)

	var result EmbeddingResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Embedding, nil
}
