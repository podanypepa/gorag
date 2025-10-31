// Package main
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// OllamaResponse represents the structure of the response from the Ollama API.
type OllamaResponse struct {
	Model              string  `json:"model"`
	CreatedAt          string  `json:"created_at"`
	Response           string  `json:"response"`
	Done               bool    `json:"done"`
	DoneReason         string  `json:"done_reason"`
	TotalDuration      float64 `json:"total_duration"`
	EvalCount          int     `json:"eval_count"`
	EvalDuration       float64 `json:"eval_duration"`
	LoadDuration       float64 `json:"load_duration"`
	PromptEvalCount    int     `json:"prompt_eval_count"`
	PromptEvalDuration float64 `json:"prompt_eval_duration"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Ussage: go run ollama_stream.go \"Your prompt here\"")
		os.Exit(1)
	}
	prompt := os.Args[1]

	start := time.Now()

	payload := map[string]any{
		"model": "mistral",
		// "model":  "llama3",
		"prompt": prompt,
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "http://localhost:11434/api/generate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ollama could not be reached. Is the Ollama server running?")
		os.Exit(1)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var last OllamaResponse

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			break
		}

		if len(line) == 0 {
			continue
		}

		var msg OllamaResponse
		if err := json.Unmarshal(line, &msg); err != nil {
			continue
		}

		if msg.Response != "" {
			fmt.Print(msg.Response)
		}

		last = msg
		if msg.Done {
			break
		}
	}

	elapsed := time.Since(start)

	fmt.Println("\n\n--- 📊 Statistiky ---")
	fmt.Printf("Model: %s\n", last.Model)
	fmt.Printf("Doba trvání: %.2f s\n", elapsed.Seconds())
	fmt.Printf("Tokeny (vstupní/výstupní): %d / %d\n", last.PromptEvalCount, last.EvalCount)
	if last.EvalDuration > 0 {
		speed := float64(last.EvalCount) / (last.EvalDuration / 1e9)
		fmt.Printf("Rychlost: %.1f tokenů/s\n", speed)
	}
	fmt.Printf("Celková doba generování: %.2f s\n", last.TotalDuration/1e9)
}
