package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	store, _ := LoadStore("index.json")

	files, _ := filepath.Glob("docs/*.pdf")
	for _, f := range files {
		fmt.Println("Pocessing:", f)
		text, _ := ExtractText(f)
		chunks := splitText(text, 800)

		for _, ch := range chunks {
			emb, _ := GetEmbedding(ch)
			store.Add(Document{Text: ch, Embedding: emb})
		}
	}

	store.Save("index.json")
	fmt.Println("✅ Index created with", len(store.Docs), "documents.")

	app := fiber.New()
	app.Get("/ask", func(c *fiber.Ctx) error {
		start := time.Now()

		q := c.Query("q")
		if q == "" {
			return c.Status(400).SendString("Missing query parameter 'q'")
		}

		qEmb, _ := GetEmbedding(q)
		results := store.Search(qEmb, 3)

		var context string
		for _, r := range results {
			context += r.Text + "\n"
		}

		prompt := fmt.Sprintf(
			"Použij následující kontext k odpovědi na dotaz:\n\n%s\n\nOtázka: %s",
			context, q,
		)
		inputTokens := estimateTokens(prompt)

		c.Set("Content-Type", "text/plain; charset=utf-8")
		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			ch := make(chan string)
			go StreamOllama("llama3", prompt, ch)

			var outBuilder strings.Builder

			for chunk := range ch {
				outBuilder.WriteString(chunk)
				_, _ = w.WriteString(chunk)
				_ = w.Flush()
			}
			w.WriteString("\n")
			w.Flush()

			outputText := outBuilder.String()
			outputTokens := estimateTokens(outputText)
			elapsed := time.Since(start)

			fmt.Printf(
				"✅ Prompt: %q | duration: %v | input tokens: %d | output tokens: %d\n",
				q, elapsed, inputTokens, outputTokens,
			)
		})

		return nil
	})

	fmt.Println("🌐 Server is running on http://localhost:8080")
	if err := app.Listen(":9090"); err != nil {
		fmt.Println("❌ Server failed to start", err)
		os.Exit(1)
	}
}

func splitText(text string, size int) []string {
	var chunks []string
	words := strings.Fields(text)
	for i := 0; i < len(words); i += size {
		end := i + size
		end = min(end, len(words))
		chunks = append(chunks, strings.Join(words[i:end], " "))
	}
	return chunks
}

func estimateTokens(s string) int {
	words := len(strings.Fields(s))
	runes := len([]rune(s))
	wTok := int(float64(words) / 0.75)
	cTok := runes / 4
	if wTok > cTok {
		return wTok
	}
	return cTok
}
