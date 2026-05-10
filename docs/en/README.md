# lanchaingo Documentation

> Go LLM Application Framework — Covers LLM Calls, Agents, RAG, LangGraph Workflows

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [LLM Calls](#llm-calls)
3. [Structured Output](#structured-output)
4. [Prompts](#prompts)
5. [Memory](#memory)
6. [Chains](#chains)
7. [Agent](#agent)
8. [Tools](#tools)
9. [RAG (Retrieval Augmented Generation)](#rag-retrieval-augmented-generation)
10. [BM25 + Hybrid Search](#bm25--hybrid-search)
11. [LangGraph Workflow](#langgraph-workflow)
12. [Human-in-the-Loop](#human-in-the-loop)
13. [Callbacks](#callbacks)
14. [Embeddings](#embeddings)
15. [Vector Stores](#vector-stores)
16. [Rate Limiting](#rate-limiting)
17. [Retry Mechanism](#retry-mechanism)
18. [Configuration](#configuration)
19. [Testing](#testing)

---

## Quick Start

```go
package main

import (
    "github.com/atliliw/lanchaingo/llms/providers"
    "github.com/atliliw/lanchaingo/schema/messages"
)

func main() {
    // Using Qwen via DashScope (Alibaba Cloud)
    llm := providers.NewQwenChat("sk-xxx")
    llm.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
    llm.Model = "qwen-plus"

    msg := messages.NewHumanMessage("Hello, introduce yourself")
    result, err := llm.Chat(nil, []messages.Message{msg})
    if err != nil {
        panic(err)
    }
    println(result.Content)
}
```

## LLM Calls

### Supported Providers

| Provider | Package | Constructor |
|----------|---------|-------------|
| OpenAI | `llms/openai` | `openai.NewOpenAIChat(config)` |
| Ollama | `llms/ollama` | `ollama.NewOllamaChat(config)` |
| Qwen (DashScope) | `llms/providers` | `providers.NewQwenChat(apiKey)` |
| DeepSeek | `llms/providers` | `providers.NewDeepSeekChat(apiKey)` |
| Moonshot | `llms/providers` | `providers.NewMoonshotChat(apiKey)` |
| Zhipu GLM | `llms/providers` | `providers.NewZhipuChat(apiKey)` |
| Anthropic Claude | `llms/providers` | `providers.NewAnthropicChat(apiKey)` |
| Google Gemini | `llms/providers` | `providers.NewGeminiChat(apiKey)` |

### Streaming

```go
ch, err := llm.StreamChat(ctx, []messages.Message{msg})
for chunk := range ch {
    fmt.Print(chunk)
}
```

### Multi-turn

```go
history := []messages.Message{
    messages.NewHumanMessage("My name is John"),
    messages.NewAIMessage("Hello John!"),
    messages.NewHumanMessage("What's my name?"),
}
result, _ := llm.Chat(ctx, history)
```

## Structured Output

Bind a Go struct to the Chat Model for type-safe structured responses.

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

model := language_models.NewStructuredOutputModel[Person](llm)
person, err := model.Invoke(ctx, "John is 30 years old")
// person.Name = "John", person.Age = 30
```

## Prompts

```go
// Template
pt := prompts.NewPromptTemplate("Hello {name}", []string{"name"})
result, _ := pt.Format(map[string]string{"name": "Alice"})

// Chat template
cpt := prompts.NewChatPromptTemplate()
cpt.AddSystemMessage("You are a {role}")
cpt.AddHumanMessage("Hello {name}")
msgs, _ := cpt.FormatMessages(map[string]string{"role": "assistant", "name": "Bob"})

// Few-Shot
fpt := prompts.NewFewShotPromptTemplate(examples, examplePrompt, "Now: {q}", nil, "\n", "")
```

## Memory

```go
// BufferMemory — full conversation history
mem := memory.NewConversationBufferMemory()
mem.SaveContext(map[string]string{"input":"hello"}, map[string]string{"output":"hi"})

// WindowMemory — last K turns
win := memory.NewConversationBufferWindowMemory(5)

// SummaryMemory — LLM-generated summaries
sum := memory.NewConversationSummaryMemory(llm)

// MongoDB persistent memory
mongoMem, _ := memory.NewMongoPersistentMemory(ctx, memory.MongoMemoryConfig{
    URI: "mongodb://localhost:27017", Database: "mydb", Collection: "memory",
})

// Store (cross-session long-term memory)
store := memory.NewInMemoryStore()
store.Put("user_123", "preferences", map[string]any{"lang": "en"})
```

## Chains

```go
// LLMChain
chain := chains.NewLLMChain(llm, "Answer in one sentence: {question}")
result, _ := chain.Invoke(map[string]any{"question": "What is Go?"})

// SequentialChain
seq := chains.NewSequentialChain().
    AddChain(chain1, []string{"input"}, []string{"output1"}).
    AddChain(chain2, []string{"output1"}, []string{"final"})

// ConversationChain with memory
conv := chains.NewConversationChain(llm, memory.NewConversationBufferMemory()).
    WithSystemPrompt("You are a helpful assistant")

// RouterChain — keyword-based routing
router := chains.NewRouterChain().
    AddRouteWithKeywords("math", "math problems", mathChain, []string{"calculate", "math"}).
    WithDefault(defaultChain)

// RetrievalQA
qa := chains.NewRetrievalQA(llm, retriever)
result, _ := qa.Invoke(map[string]any{"query": "What are Rust's features?"})
```

## Agent

```go
// ReAct Agent
agent := react.NewReActAgent(llm, tools, "")
exec := agents.NewAgentExecutor(agent, tools).WithVerbose(true)
result, _ := exec.Invoke("Calculate 37 + 48")

// FunctionCalling Agent
fcAgent := function_calling.NewFunctionCallingAgent(llm, tools, "")
exec2 := agents.NewAgentExecutor(fcAgent, tools)
result2, _ := exec2.Invoke("What's today's date?")
```

## Tools

```go
tools := []tools.Tool{
    tools.NewCalculator(),       // Math expressions
    tools.NewSimpleMathTool(),   // sqrt, pow, abs
    tools.NewDateTimeTool(),     // Current date/time
    tools.NewURLFetchTool(),     // HTTP GET
    tools.NewWikipediaTool(),    // Wikipedia search
    tools.NewDuckDuckGoSearchTool(), // Web search
    tools.NewPythonREPLTool(),   // Python execution
}
```

## RAG (Retrieval Augmented Generation)

```go
// 1. Documents
docs := []vs.Document{
    vs.NewDocument("Rust is a systems programming language"),
    vs.NewDocument("Go excels at concurrent programming"),
}

// 2. Embed
emb := embeddings.NewOpenAIEmbeddings(cfg)
vecs, _ := emb.EmbedDocuments([]string{docs[0].Content, docs[1].Content})

// 3. Store
store := vs.NewInMemoryVectorStore()
store.AddDocuments(docs, vecs)

// 4. Retrieve
retriever := retrieval.NewSimilarityRetriever(store, emb)
results, _ := retriever.Retrieve("concurrent programming", 2)

// 5. Generate
msg := messages.NewHumanMessage(results[0].Content + " Summarize in one sentence")
answer, _ := llm.Chat(ctx, []messages.Message{msg})
```

### HyDE (Hypothetical Document Embeddings)

```go
hyde := retrieval.NewHyDERetriever(llm, emb, store)
results, _ := hyde.RetrieveWithScores("Rust vs Go comparison", 5)
```

### MultiQuery

```go
mq := retrieval.NewMultiQueryRetriever(llm, retriever)
docs, _ := mq.RetrieveWithMultiQuery("memory safety")
```

### Reranking

```go
reranker := retrieval.NewKeywordReranker()
executor := retrieval.NewRerankingExecutor(reranker)
reranked, _ := executor.Rerank("query", results)
```

## BM25 + Hybrid Search

```go
// BM25 keyword search
bm25 := bm25.NewBM25Retriever()
bm25.AddDocument(vs.NewDocument("Rust systems programming"))
results := bm25.Search("programming language", 5)

// Chunked BM25 (Parent-Child)
chunked := bm25.NewChunkedBM25Retriever()
chunked.AddDocument(parentDoc, leafDocs, leafTerms)

// Hybrid (BM25 + Vector + RRF)
fused := retrieval.ReciprocalRankFusion(bm25Docs, vectorDocs, 60)
```

## LangGraph Workflow

```go
// Define graph
g := langgraph.NewStateGraph()
g.AddNodeFn("step1", func(s StateSchema) (StateUpdate, error) {
    return NewStateUpdate(s.CloneState()), nil
})
g.AddEdge(START, "step1")
g.AddEdge("step1", END)

// Compile & execute
cg, _ := g.Compile()
result, _ := cg.Invoke(AgentState{Input: "test"})

// Conditional routing
targets := map[string]string{"a": "process_a", "b": "process_b"}
g.AddConditionalEdges("decide", "router", targets, nil)
g.SetConditionalRouter("router", func(s StateSchema) string {
    if s.(AgentState).Input == "a" { return "a" }
    return "b"
})
```

## Human-in-the-Loop

Pause workflow execution for human approval, then resume.

```go
cg, _ := g.Compile()
cg.InterruptBefore("approve") // Pause before approve node

result := cg.Invoke(input)
if result.Interrupted {
    println("Interrupted at:", result.InterruptedAt)

    // Resume execution after human intervention
    result = cg.Resume(updatedInput)
}
```

## Callbacks

```go
cm := callbacks.NewCallbackManager()
cm.AddHandler(&handlers.StdOutHandler{})     // Console output
fh, _ := handlers.NewFileCallbackHandler("trace.log")
cm.AddHandler(fh)                              // File logging
ls := callbacks.NewLangSmithHandler("key", "project")
cm.AddHandler(ls)                               // LangSmith tracing
```

## Embeddings

```go
emb := embeddings.NewOpenAIEmbeddings(cfg)
vec, _ := emb.EmbedQuery("How's the weather?")
println(len(vec)) // 1536

vecs, _ := emb.EmbedDocuments([]string{"doc1", "doc2"})
sim := embeddings.CosineSimilarity(vecs[0], vecs[1])

mock := embeddings.NewMockEmbeddings(128)
```

## Vector Stores

```go
// InMemory (default)
store := vs.NewInMemoryVectorStore()

// Qdrant (requires Qdrant service)
qdrant, _ := vs.NewQdrantVectorStore(ctx, vs.QdrantConfig{
    URL: "http://192.168.10.100:6334", Collection: "docs", VectorSize: 1536,
})

// MongoDB
mongo, _ := vs.NewMongoVectorStore(ctx, vs.MongoConfig{
    URI: "mongodb://localhost:27017", Database: "mydb", Collection: "vectors",
})
```

## Rate Limiting

```go
rl := core.NewRateLimiter(10, time.Second) // 10 requests/second
rl.Wait(ctx) // Blocks until token available
```

## Retry Mechanism

All HTTP calls have built-in exponential backoff (3 retries, 1s→2s→4s + 25% jitter).

Custom retry:

```go
retryCfg := core.DefaultRetryConfig()
retryCfg.MaxRetries = 5
result, err := core.DoWithRetry(ctx, retryCfg, func(ctx context.Context) (string, error) {
    return httpCall()
})
```

## Configuration

Copy `config/config.toml.example` to `config/config.toml`:

```toml
[openai]
api_key = "sk-xxx"
base_url = "https://dashscope.aliyuncs.com/compatible-mode/v1"
chat_model = "qwen-plus"
embedding_model = "text-embedding-v1"
```

`config.toml` is in `.gitignore` and will not be committed.

## Testing

```bash
# Unit tests (no API key required)
go test ./...

# Integration tests (requires config.toml with API key)
go test -tags=integration ./tests/integration/ -v

# End-to-end tests
go test -tags=e2e ./tests/e2e/ -v
```

---

> **Commercial Use Notice**: This framework is proprietary commercial software. Unauthorized commercial use is prohibited. For licensing, support, and enterprise plans, contact the author.
