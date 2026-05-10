# lanchaingo

**Go LLM Application Framework**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)

lanchaingo is a Go framework for building LLM applications. It covers the full stack: LLM calls, Agents, RAG pipelines, LangGraph workflows, BM25 search, and hybrid retrieval.

## Quick Start

```go
package main

import (
    "github.com/atliliw/lanchaingo/llms/providers"
    "github.com/atliliw/lanchaingo/schema/messages"
)

func main() {
    llm := providers.NewQwenChat("sk-xxx")
    llm.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
    llm.Model = "qwen-plus"

    msg := messages.NewHumanMessage("Hello")
    result, _ := llm.Chat(nil, []messages.Message{msg})
    println(result.Content)
}
```

## Features

| Module | Status |
|--------|--------|
| LLM Providers (OpenAI, Ollama, Qwen, DeepSeek, Anthropic, Gemini, etc.) | ✅ |
| Agent (ReAct + FunctionCalling + Executor) | ✅ |
| Memory (Buffer, Window, Summary, Persistent, MongoDB) | ✅ |
| Chains (LLMChain, Sequential, Conversation, Router, RAG, Document) | ✅ |
| LangGraph workflow engine | ✅ |
| RAG pipeline (Embedding -> VectorStore -> Retrieval -> Generation) | ✅ |
| BM25 keyword search + Hybrid search (RRF) | ✅ |
| Tools (Calculator, Search, PythonREPL, etc.) | ✅ |
| Structured output `WithStructuredOutput[T]` | ✅ |
| Human-in-the-Loop | ✅ |
| Exponential backoff retry | ✅ |
| Rate limiting | ✅ |

## Documentation

See [docs/en/README.md](docs/en/README.md) for full documentation.

## Testing

```bash
go test ./...                    # Unit tests
go test -tags=integration ./tests/integration/ -v  # Integration tests
```

## License

Proprietary. All rights reserved.
