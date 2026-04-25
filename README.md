# 🚀 GoRAG - Your Local RAG Assistant in Go

GoRAG is a simple yet powerful RAG (Retrieval-Augmented Generation) application written in Go. It allows you to create a smart "knowledge base" from your documents (PDF and Markdown) and interact with it through a modern web interface or API. Everything runs locally using [Ollama](https://ollama.ai/), ensuring your data privacy.

## ✨ Key Features

- **Multi-Format Extraction:** Automatically processes **PDF** and **Markdown** files from your documents directory.
- **Modern Web UI:** Beautiful, responsive chat interface with real-time response streaming.
- **Persistent Vector Store:** Powered by **BadgerDB**, providing high-performance local storage (no more massive JSON files).
- **Source Citations:** The assistant cites specific sources (e.g., `[Source 1]`) in its answers so you can verify the information.
- **Smart Chunking:** Text is split into chunks with **configurable overlap** to preserve context across boundaries.
- **Docker Ready:** Complete `docker-compose.yml` included for easy deployment with Ollama.
- **Graceful Shutdown:** Cleanly handles termination signals to ensure data integrity.
- **Index Management CLI:** Build your index on demand with a dedicated flag.
- **Structured Logging:** Uses Go's `slog` for clean, professional logging.

## 🏁 Getting Started

### Prerequisites

- **Go:** Version 1.25 or higher.
- **Ollama:** Installed and running with your preferred model (e.g., `llama3`).
  ```bash
  ollama pull llama3
  ```

### Installation and Setup (Local)

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/podanypepa/gorag.git
    cd gorag
    ```

2.  **Add your documents:**
    Place your `.pdf` or `.md` files into the `docs/` directory.

3.  **Build the Index:**
    ```bash
    go run . --index
    ```

4.  **Run the Server:**
    ```bash
    go run .
    ```
    Visit `http://localhost:9090` in your browser.

### Installation and Setup (Docker)

1. **Start the containers:**
   ```bash
   docker-compose up -d
   ```
2. **Download the model (first time):**
   ```bash
   docker exec -it gorag-ollama-1 ollama run llama3
   ```
3. **Index your documents:**
   ```bash
   docker-compose run gorag ./gorag --index
   ```

## ⚙️ Configuration

| Variable      | Description                                  | Default Value                    |
|---------------|----------------------------------------------|----------------------------------|
| `OLLAMA_URL`  | The URL of the Ollama API instance.          | `http://localhost:11434`         |
| `MODEL_NAME`  | The name of the model to use.                | `llama3`                         |
| `INDEX_DIR`   | The directory for the BadgerDB store.        | `index_db`                       |
| `PDF_DIR`     | The directory containing documents.          | `docs/`                          |
| `SERVER_PORT` | The port for the web server.                 | `9090`                           |

## 🔌 API Usage

Ask questions via simple GET requests:
```bash
curl "http://localhost:9090/ask?q=What are the key findings?"
```

## 🖥️ CLI Client

A simple streaming CLI client is available in `bin/client`.

```bash
cd bin/client
go build -o ollama_stream .
./ollama_stream "Explain quantum computing."
```

---
*This project explores the power of local RAG systems using Go. Secure, fast, and fully local.*
