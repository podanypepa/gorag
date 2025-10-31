# Ollama CLI Streamer

This is a simple command-line application to stream responses from the Ollama API to your terminal.

## Installation

To build the application, run the following command in this directory:

```bash
go build -o ollama_stream .
```

## Usage

To use the application, run the executable with a prompt.

```bash
./ollama_stream "Why is the sky blue?"
```

### Options

You can also specify the model and the Ollama API URL using flags:

- `-model`: The model to use for the generation (default: `mistral`).
- `-url`: The Ollama API URL (default: `http://localhost:11434/api/generate`).

## Example

Here is an example of how to use the application with a different model:

```bash
./ollama_stream -model llama3 "Tell me a joke about a programmer."
```
