// Package main implements a document retrieval and question-answering server.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	indexOnly := flag.Bool("index", false, "Build the index from documents and exit")
	flag.Parse()

	if err := run(*indexOnly); err != nil {
		slog.Error("Application failed", "error", err)
		os.Exit(1)
	}
}

func run(indexOnly bool) error {
	indexDir := os.Getenv("INDEX_DIR")
	if indexDir == "" {
		indexDir = "index_db"
	}

	slog.Info("Opening vector store", "directory", indexDir)
	store, err := NewVectorStore(indexDir)
	if err != nil {
		return fmt.Errorf("failed to open store: %w", err)
	}
	defer store.Close()

	if indexOnly {
		if err := buildIndex(store); err != nil {
			return fmt.Errorf("failed to build index: %w", err)
		}
		slog.Info("Index built successfully", "documents", len(store.Docs))
		return nil
	}

	if len(store.Docs) == 0 {
		slog.Warn("Vector store is empty. Use --index to build it.")
	}

	return startServer(store)
}

func buildIndex(store *VectorStore) error {
	pdfDir := os.Getenv("PDF_DIR")
	if pdfDir == "" {
		pdfDir = "docs/"
	}
	slog.Info("Building index from files", "directory", pdfDir)

	// Create a map of already indexed files
	indexedFiles := make(map[string]bool)
	for _, doc := range store.Docs {
		indexedFiles[doc.Source] = true
	}

	entries, err := os.ReadDir(pdfDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", pdfDir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		f := filepath.Join(pdfDir, entry.Name())
		if indexedFiles[f] {
			slog.Debug("File already indexed, skipping", "path", f)
			continue
		}

		ext := strings.ToLower(filepath.Ext(f))
		var text string
		var err error

		switch ext {
		case ".md":
			slog.Info("Processing Markdown file", "path", f)
			text, err = ExtractMDText(f)
		case ".pdf":
			slog.Info("Processing PDF file", "path", f)
			text, err = ExtractText(f)
		default:
			continue
		}

		if err != nil {
			slog.Error("Failed to extract text", "file", f, "error", err)
			continue
		}

		chunks := splitText(text, 800, 150)
		for _, ch := range chunks {
			emb, err := GetEmbedding(ch)
			if err != nil {
				slog.Error("Failed to get embedding", "chunk", ch, "error", err)
				continue
			}
			store.Add(Document{
				Text:      ch,
				Source:    f,
				Embedding: emb,
			})
		}
	}
	return nil
}

func startServer(store *VectorStore) error {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("index.html")
	})

	app.Get("/ask", func(c *fiber.Ctx) error {
		start := time.Now()

		q := c.Query("q")
		if q == "" {
			return c.Status(400).SendString("Missing query parameter 'q'")
		}
		slog.Info("Received query", "query", q)

		qEmb, err := GetEmbedding(q)
		if err != nil {
			return c.Status(500).SendString("Failed to get query embedding")
		}
		results := store.Search(qEmb, 3)

		var context strings.Builder
		for i, r := range results {
			context.WriteString(fmt.Sprintf("[Source %d: %s]\n", i+1, filepath.Base(r.Source)))
			context.WriteString(r.Text)
			context.WriteString("\n\n")
		}

		model := os.Getenv("MODEL_NAME")
		if model == "" {
			model = "llama3"
		}

		prompt := fmt.Sprintf(
			"You are a helpful assistant. Use the following context to answer the query. "+
				"If you find the answer in the context, cite the source number (e.g., [Source 1]). "+
				"If the answer is not in the context, state that you do not know.\n\nContext:\n%s\n\nQuestion: %s",
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

	// Graceful shutdown setup
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		slog.Info("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			slog.Error("Server shutdown failed", "error", err)
		}
		close(idleConnsClosed)
	}()

	slog.Info("Server is running", "port", port)
	if err := app.Listen(":" + port); err != nil {
		slog.Error("Server failed to start", "error", err)
	}

	<-idleConnsClosed
	slog.Info("Server stopped")
	return nil
}
