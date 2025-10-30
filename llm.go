package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// GenerateRequest represents the request payload for the Ollama API.
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// GenerateChunk represents a chunk of the response from the Ollama API.
type GenerateChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// StreamOllama streams responses from the Ollama API for the given model and prompt.
func StreamOllama(model, prompt string, ch chan string) {
	defer close(ch)

	ollamaURL := os.Getenv("OLLAMA_URL")
	if ollamaURL == "" {
		ollamaURL = "http://localhost:11434"
	}

	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	}

	data, err := json.Marshal(req)
	if err != nil {
		ch <- fmt.Sprintf("[chyba marshalling: %v]", err)
		return
	}

	resp, err := http.Post(ollamaURL+"/api/generate", "application/json", bytes.NewReader(data))
	if err != nil {
		ch <- fmt.Sprintf("[chyba připojení: %v]", err)
		return
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		var chunk GenerateChunk
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		if chunk.Done {
			break
		}
		ch <- chunk.Response
	}
}
