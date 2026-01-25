// Package main
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

// OllamaResponse represents the structure of the response from the Ollama API.
type OllamaResponse struct {
	Model              string    `json:"model"`
	CreatedAt          time.Time `json:"created_at"`
	Response           string    `json:"response"`
	Done               bool      `json:"done"`
	DoneReason         string    `json:"done_reason,omitempty"`
	TotalDuration      float64   `json:"total_duration,omitempty"`
	EvalCount          int       `json:"eval_count,omitempty"`
	EvalDuration       float64   `json:"eval_duration,omitempty"`
	LoadDuration       float64   `json:"load_duration,omitempty"`
	PromptEvalCount    int       `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration float64   `json:"prompt_eval_duration,omitempty"`
}

func main() {
	model := flag.String("model", "gpt-oss:20b", "Model to use for the generation")
	url := flag.String("url", "http://localhost:11434/api/generate", "Ollama API URL")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <prompt>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Error: Prompt is required.")
		flag.Usage()
		os.Exit(1)
	}
	prompt := strings.Join(flag.Args(), " ")

	start := time.Now()

	last, err := streamOllama(*url, *model, prompt)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	printStats(last, time.Since(start))
}

func streamOllama(url, model, prompt string) (OllamaResponse, error) {
	var last OllamaResponse

	payload := map[string]any{
		"model":  model,
		"prompt": prompt,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return last, fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return last, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return last, fmt.Errorf("ollama could not be reached. Is the Ollama server running? %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return last, fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	reader := bufio.NewReader(resp.Body)
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
			fmt.Fprintf(os.Stderr, "Warning: failed to unmarshal line: %v\n", err)
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
	return last, nil
}

func printStats(last OllamaResponse, elapsed time.Duration) {
	fmt.Println("\n\n--- 📊 Statistics ---")
	fmt.Printf("Model: %s\n", last.Model)
	fmt.Printf("Duration: %.2f s\n", elapsed.Seconds())
	if last.PromptEvalCount > 0 && last.EvalCount > 0 {
		fmt.Printf("Tokens (prompt/generated): %d / %d\n", last.PromptEvalCount, last.EvalCount)
	}
	if last.EvalDuration > 0 {
		speed := float64(last.EvalCount) / (last.EvalDuration / 1e9)
		fmt.Printf("Speed: %.1f tokens/s\n", speed)
	}
	if last.TotalDuration > 0 {
		fmt.Printf("Total generation time: %.2f s\n", last.TotalDuration/1e9)
	}
}
