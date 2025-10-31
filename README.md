# 🚀 GoRAG - Your Local RAG Assistant in Go

GoRAG is a simple yet powerful RAG (Retrieval-Augmented Generation) application written in Go. It allows you to create a smart "knowledge base" from your PDF documents and ask it questions in natural language. Everything runs locally using [Ollama](https://ollama.ai/), so your data remains secure.

## ✨ Key Features

- **PDF Text Extraction:** Automatically reads and processes all PDF files from a specified directory.
- **Vector Embedding Generation:** Creates semantic representations of text chunks using a locally running language model via Ollama.
- **In-Memory Vector Search:** Fast and efficient retrieval of relevant information thanks to a custom in-memory vector store implementation.
- **Answer Generation:** Leverages the power of Large Language Models (LLMs) to synthesize answers based on the retrieved context.
- **Response Streaming:** Answers are streamed in real-time, ensuring a smooth user experience.
- **Simple REST API:** An easy-to-integrate interface for asking questions.
- **Flexible Configuration:** Key parameters can be configured via environment variables.

## 🏁 Getting Started

### Prerequisites

- **Go:** Version 1.22 or higher.
- **Ollama:** An installed and running instance of Ollama with at least one model downloaded (e.g., `llama3`).
  ```bash
  ollama pull llama3
  ```

### Installation and Setup

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/podanypepa/gorag.git
    cd gorag
    ```

2.  **Add your documents:**
    Place your PDF files into the `docs/` directory.

3.  **Run the application:**
    ```bash
    go run .
    ```

On the first run, the application will create an `index.json` file from all found PDF documents. It will then start the web server.

## ⚙️ Configuration

The application can be configured using the following environment variables:

| Variable      | Description                                  | Default Value                    |
|---------------|----------------------------------------------|----------------------------------|
| `OLLAMA_URL`  | The URL of the running Ollama API instance.  | `http://localhost:11434`         |
| `MODEL_NAME`  | The name of the model to use.                | `llama3`                         |
| `INDEX_FILE`  | The path to the vector index file.           | `index.json`                     |
| `PDF_DIR`     | The directory containing PDF documents to index. | `docs/`                          |
| `SERVER_PORT` | The port on which the web server will run.   | `9090`                           |

**Example of running with custom configuration:**
```bash
MODEL_NAME=mistral SERVER_PORT=8888 go run .
```

## 🔌 API Usage

After starting the server, you can ask questions by making a simple GET request to the `/ask` endpoint.

**Example using `curl`:**
```bash
curl "http://localhost:9090/ask?q=What is the main topic of the document?"
```

The response will be streamed as plain text.

## 🖥️ CLI Client

In the `bin/client` directory, you'll find a simple command-line client that can stream responses from the Ollama API directly to your terminal.

### Build

```bash
cd bin/client
go build -o ollama_stream .
```

### Usage

```bash
./ollama_stream "Why is the sky blue?"
```

You can also specify the model:
```bash
./ollama_stream -model llama3 "Tell me a joke about a programmer."
```

## 🔮 Future Improvements

-   [ ] Replace the in-memory store with a robust vector database (e.g., ChromaDB, Weaviate).
-   [ ] Implement more advanced text chunking strategies.
-   [ ] Use a dedicated tokenizer for more accurate token counting.
-   [ ] Extend the API with more endpoints (document management, index status).
-   [ ] Containerize the application with Docker for easier deployment.

---
*This project was created to explore the capabilities of RAG applications in Go. It serves as a great starting point for your own experiments!*