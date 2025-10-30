package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

	req := GenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: true,
	}

	data, _ := json.Marshal(req)
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(data))
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
