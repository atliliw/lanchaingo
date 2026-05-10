# lanchaingo

**Go 语言 LLM 应用框架 — 1:1 复刻 langchainrust，对标 Python LangChain**

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-Commercial-red)](#license)
[![Go Reference](https://pkg.go.dev/badge/github.com/atliliw/lanchaingo)](https://pkg.go.dev/github.com/atliliw/lanchaingo)

---

lanchaingo 是一个用 Go 语言编写的 LLM 应用开发框架，完整覆盖 **LLM 调用 → Agent → RAG → 工作流编排** 全链路。核心代码 1:1 映射 [langchainrust](https://github.com/atliliw/langchainrust) v0.2.20 的 92 个 Feature 模块，并额外实现了对标 Python LangChain 的生产级商业功能。

## 快速开始

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

    msg := messages.NewHumanMessage("你好")
    result, _ := llm.Chat(nil, []messages.Message{msg})
    println(result.Content)
}
```

## Feature 总览

| 模块 | 覆盖情况 |
|------|---------|
| LLM Provider（OpenAI, Ollama, Qwen, DeepSeek, Anthropic, Gemini 等 8 个） | ✅ |
| Agent（ReAct + Function Calling + Executor） | ✅ |
| Memory（Buffer, Window, Summary, Persistent, MongoDB） | ✅ |
| Chains（LLMChain, Sequential, Conversation, Router, RAG, Document 等 9 种） | ✅ |
| LangGraph 工作流引擎 | ✅ |
| RAG 管线（嵌入 → 向量存储 → 检索 → 生成） | ✅ |
| BM25 关键词搜索 + 混合检索（RRF） | ✅ |
| Tools（Calculator, Search, PythonREPL 等 7 个） | ✅ |
| 结构化输出 `WithStructuredOutput[T]` | ✅ |
| Human-in-the-Loop | ✅ |
| 指数退避重试 | ✅ |
| Rate Limiting | ✅ |
| 集成测试（真实 API） | ✅ |

## 文档

详见 [docs/](docs/) 目录：
- [中文文档](docs/zh/README.md)
- [English Docs](docs/en/README.md)

## 许可

商业私有。未经授权不得用于商业用途。
