// Package main implements a document retrieval and question-answering server.
package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	if err := run(); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	indexFile := os.Getenv("INDEX_FILE")
	if indexFile == "" {
		indexFile = "index.json"
	}

	store, err := LoadStore(indexFile)
	if err != nil {
		return fmt.Errorf("failed to load store: %w", err)
	}

	if err := buildIndex(store); err != nil {
		return fmt.Errorf("failed to build index: %w", err)
	}

	if err := store.Save(indexFile); err != nil {
		return fmt.Errorf("failed to save store: %w", err)
	}

	slog.Info("Index created", "documents", len(store.Docs))

	return startServer(store)
}

func buildIndex(store *VectorStore) error {
	pdfDir := os.Getenv("PDF_DIR")
	if pdfDir == "" {
		pdfDir = "docs/"
	}

	files, err := filepath.Glob(filepath.Join(pdfDir, "*.pdf"))
	if err != nil {
		return fmt.Errorf("failed to find PDF files: %w", err)
	}

	for _, f := range files {
		slog.Info("Processing file", "path", f)
		text, err := ExtractText(f)
		if err != nil {
			slog.Error("Failed to extract text", "file", f, "error", err)
			continue
		}
		chunks := splitText(text, 800)

		for _, ch := range chunks {
			emb, err := GetEmbedding(ch)
			if err != nil {
				slog.Error("Failed to get embedding", "chunk", ch, "error", err)
				continue
			}
			store.Add(Document{Text: ch, Embedding: emb})
		}
	}
	return nil
}

func startServer(store *VectorStore) error {
	app := fiber.New()
	app.Get("/ask", func(c *fiber.Ctx) error {
		start := time.Now()

		q := c.Query("q")
		if q == "" {
			return c.Status(400).SendString("Missing query parameter 'q'")
		}

		qEmb, err := GetEmbedding(q)
		if err != nil {
			return c.Status(500).SendString("Failed to get query embedding")
		}
		results := store.Search(qEmb, 3)

		var context strings.Builder
		for _, r := range results {
			context.WriteString(r.Text)
			context.WriteString("\n")
		}

		model := os.Getenv("MODEL_NAME")
		if model == "" {
			model = "llama3"
		}

		prompt := fmt.Sprintf(
			"Použij následující kontext k odpovědi na dotaz:\n\n%s\n\nOtázka: %s",
			context.String(), q,
		)
		inputTokens := estimateTokens(prompt)

		c.Set("Content-Type", "text/plain; charset=utf-8")
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			ch := make(chan string)
			go StreamOllama(model, prompt, ch)

			var outBuilder strings.Builder

			for chunk := range ch {
				outBuilder.WriteString(chunk)
				_, _ = w.WriteString(chunk)
				_ = w.Flush()
			}
			_, _ = w.WriteString("\n")
			_ = w.Flush()

			outputText := outBuilder.String()
			outputTokens := estimateTokens(outputText)
			elapsed := time.Since(start)

			slog.Info("Query processed",
				"prompt", q,
				"duration", elapsed,
				"input_tokens", inputTokens,
				"output_tokens", outputTokens,
			)
		})

		return nil
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "9090"
	}

	slog.Info("Server is running", "port", port)
	return app.Listen(":" + port)
}
