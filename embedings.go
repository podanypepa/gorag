package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
	model := os.Getenv("MODEL_NAME")
	if model == "" {
		model = "llama3"
	}

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	req := EmbeddingRequest{Model: model, Prompt: text}
	data, err := json.Marshal(req)
	if err != nil {
		slog.Error("Failed to marshal embedding request", "error", err)
		return nil, err
	}

	resp, err := http.Post(ollamaURL+"/api/embeddings", "application/json", bytes.NewReader(data))
	if err != nil {
		slog.Error("Failed to get embedding", "error", err)
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Printf("Response status: %s\n", resp.Status)

	var result EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		slog.Error("Failed to decode embedding response", "error", err)
		return nil, err
	}
	return result.Embedding, nil
}
