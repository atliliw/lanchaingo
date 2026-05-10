# lanchaingo 中文文档

> Go 语言 LLM 应用框架 — 完整覆盖 LLM 调用、Agent、RAG、LangGraph 工作流

---

## 目录

1. [快速开始](#快速开始)
2. [LLM 调用](#llm-调用)
3. [结构化输出](#结构化输出)
4. [Prompts](#prompts)
5. [Memory（记忆）](#memory记忆)
6. [Chains（链）](#chains链)
7. [Agent（智能体）](#agent智能体)
8. [Tools（工具）](#tools工具)
9. [RAG 检索增强生成](#rag-检索增强生成)
10. [BM25 + 混合检索](#bm25--混合检索)
11. [LangGraph 工作流](#langgraph-工作流)
12. [Human-in-the-Loop](#human-in-the-loop)
13. [Callbacks（回调）](#callbacks回调)
14. [Embeddings（嵌入）](#embeddings嵌入)
15. [Vector Stores（向量存储）](#vector-stores向量存储)
16. [Rate Limiting（限流）](#rate-limiting限流)
17. [重试机制](#重试机制)
18. [配置说明](#配置说明)
19. [测试](#测试)

---

## 快速开始

```go
package main

import (
    "github.com/atliliw/lanchaingo/llms/providers"
    "github.com/atliliw/lanchaingo/schema/messages"
)

func main() {
    // 使用阿里云通义千问 (DashScope)
    llm := providers.NewQwenChat("sk-xxx")
    llm.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
    llm.Model = "qwen-plus"

    msg := messages.NewHumanMessage("你好，请介绍一下自己")
    result, err := llm.Chat(nil, []messages.Message{msg})
    if err != nil {
        panic(err)
    }
    println(result.Content)
}
```

## LLM 调用

### 支持的 Provider

| Provider | 包路径 | 构造函数 |
|----------|--------|---------|
| OpenAI | `llms/openai` | `openai.NewOpenAIChat(config)` |
| Ollama | `llms/ollama` | `ollama.NewOllamaChat(config)` |
| 通义千问 (DashScope) | `llms/providers` | `providers.NewQwenChat(apiKey)` |
| DeepSeek | `llms/providers` | `providers.NewDeepSeekChat(apiKey)` |
| Moonshot | `llms/providers` | `providers.NewMoonshotChat(apiKey)` |
| 智谱 GLM | `llms/providers` | `providers.NewZhipuChat(apiKey)` |
| Anthropic Claude | `llms/providers` | `providers.NewAnthropicChat(apiKey)` |
| Google Gemini | `llms/providers` | `providers.NewGeminiChat(apiKey)` |

### 流式对话

```go
ch, err := llm.StreamChat(ctx, []messages.Message{msg})
for chunk := range ch {
    fmt.Print(chunk) // 逐块输出
}
```

### 多轮对话

```go
history := []messages.Message{
    messages.NewHumanMessage("我的名字是张三"),
    messages.NewAIMessage("你好张三！"),
    messages.NewHumanMessage("我叫什么？"),
}
result, _ := llm.Chat(ctx, history)
```

## 结构化输出

将 Go 结构体绑定到 Chat Model，LLM 直接返回结构化数据。

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

model := language_models.NewStructuredOutputModel[Person](llm)
person, err := model.Invoke(ctx, "张三今年30岁")
// person.Name = "张三", person.Age = 30
```

## Prompts

```go
// 模板替换
pt := prompts.NewPromptTemplate("Hello {name}", []string{"name"})
result, _ := pt.Format(map[string]string{"name": "Alice"})

// 聊天模板
cpt := prompts.NewChatPromptTemplate()
cpt.AddSystemMessage("You are a {role}")
cpt.AddHumanMessage("Hello {name}")
msgs, _ := cpt.FormatMessages(map[string]string{
    "role": "assistant",
    "name": "Bob",
})

// Few-Shot
ep := prompts.NewPromptTemplate("{input} => {output}", []string{"input", "output"})
fpt := prompts.NewFewShotPromptTemplate(examples, ep, "Now: {q}", nil, "\n", "")
```

## Memory（记忆）

```go
// BufferMemory — 保存全部对话历史
mem := memory.NewConversationBufferMemory()
mem.SaveContext(map[string]string{"input":"你好"}, map[string]string{"output":"你好！"})

// WindowMemory — 只保留最近 K 轮
win := memory.NewConversationBufferWindowMemory(5)

// SummaryMemory — LLM 自动摘要
sum := memory.NewConversationSummaryMemory(llm)

// MongoDB 持久化记忆（需要 mongo driver）
mongoMem, _ := memory.NewMongoPersistentMemory(ctx, memory.MongoMemoryConfig{
    URI: "mongodb://localhost:27017", Database: "mydb", Collection: "memory",
})

// Store（长期记忆，跨会话）
store := memory.NewInMemoryStore()
store.Put("user_123", "preferences", map[string]any{"lang": "zh"})
```

## Chains（链）

```go
// LLMChain
chain := chains.NewLLMChain(llm, "用一句话回答：{question}")
result, _ := chain.Invoke(map[string]any{"question": "什么是 Go 语言？"})

// SequentialChain — 链式调用
seq := chains.NewSequentialChain().
    AddChain(chain1, []string{"input"}, []string{"output1"}).
    AddChain(chain2, []string{"output1"}, []string{"final"})

// ConversationChain — 带记忆的对话
conv := chains.NewConversationChain(llm, memory.NewConversationBufferMemory()).
    WithSystemPrompt("你是一个助手")

// RouterChain — 关键词路由
router := chains.NewRouterChain().
    AddRouteWithKeywords("math", "数学", mathChain, []string{"计算", "数学"}).
    WithDefault(defaultChain)

// RetrievalQA — RAG 问答
qa := chains.NewRetrievalQA(llm, retriever)
result, _ := qa.Invoke(map[string]any{"query": "Rust 语言的特点？"})
```

## Agent（智能体）

```go
// ReAct Agent
agent := react.NewReActAgent(llm, tools, "")
exec := agents.NewAgentExecutor(agent, tools).WithVerbose(true)
result, _ := exec.Invoke("计算 37 + 48 等于多少？")

// FunctionCalling Agent
fcAgent := function_calling.NewFunctionCallingAgent(llm, tools, "")
exec2 := agents.NewAgentExecutor(fcAgent, tools)
result2, _ := exec2.Invoke("查询今天的日期")
```

## Tools（工具）

```go
tools := []tools.Tool{
    tools.NewCalculator(),       // 计算器
    tools.NewSimpleMathTool(),   // 高级数学（sqrt, pow, abs）
    tools.NewDateTimeTool(),     // 日期时间
    tools.NewURLFetchTool(),     // URL 获取
    tools.NewWikipediaTool(),    // 维基百科
    tools.NewDuckDuckGoSearchTool(), // 网页搜索
    tools.NewPythonREPLTool(),   // Python 执行
}
```

## RAG 检索增强生成

```go
// 1. 文档
docs := []vs.Document{
    vs.NewDocument("Rust 是一门系统编程语言"),
    vs.NewDocument("Go 语言擅长并发编程"),
}

// 2. 嵌入
emb := embeddings.NewOpenAIEmbeddings(cfg)
vecs, _ := emb.EmbedDocuments([]string{docs[0].Content, docs[1].Content})

// 3. 存入向量存储
store := vs.NewInMemoryVectorStore()
store.AddDocuments(docs, vecs)

// 4. 检索
retriever := retrieval.NewSimilarityRetriever(store, emb)
results, _ := retriever.Retrieve("并发编程", 2)

// 5. 生成
msg := messages.NewHumanMessage(results[0].Content + " 用一句话总结")
answer, _ := llm.Chat(ctx, []messages.Message{msg})
```

### HyDE（假设文档嵌入）

```go
hyde := retrieval.NewHyDERetriever(llm, emb, store)
results, _ := hyde.RetrieveWithScores("Rust vs Go 对比", 5)
```

### MultiQuery（多查询变体）

```go
mq := retrieval.NewMultiQueryRetriever(llm, retriever)
docs, _ := mq.RetrieveWithMultiQuery("内存安全")
```

### Reranking（重排序）

```go
reranker := retrieval.NewKeywordReranker()
executor := retrieval.NewRerankingExecutor(reranker)
reranked, _ := executor.Rerank("query", results)
```

## BM25 + 混合检索

```go
// BM25 关键词检索
bm25 := bm25.NewBM25Retriever()
bm25.AddDocument(vs.NewDocument("Rust 系统编程语言"))
results := bm25.Search("编程语言", 5)

// 分块 BM25 (Parent-Child)
chunked := bm25.NewChunkedBM25Retriever()
chunked.AddDocument(parentDoc, leafDocs, leafTerms)

// 混合检索 (BM25 + Vector + RRF)
fused := retrieval.ReciprocalRankFusion(bm25Docs, vectorDocs, 60)
```

## LangGraph 工作流

```go
// 定义图
g := langgraph.NewStateGraph()
g.AddNodeFn("step1", func(s StateSchema) (StateUpdate, error) {
    return NewStateUpdate(s.CloneState()), nil
})
g.AddEdge(START, "step1")
g.AddEdge("step1", END)

// 编译和执行
cg, _ := g.Compile()
result, _ := cg.Invoke(AgentState{Input: "test"})

// 条件路由
targets := map[string]string{"a": "process_a", "b": "process_b"}
g.AddConditionalEdges("decide", "router", targets, nil)
g.SetConditionalRouter("router", func(s StateSchema) string {
    if s.(AgentState).Input == "a" { return "a" }
    return "b"
})
```

## Human-in-the-Loop

在工作流中插入人工审批点，执行到指定节点前中断，等待外部输入后恢复。

```go
// 编译时配置中断点
cg, _ := g.Compile()
cg.InterruptBefore("approve") // 执行到 approve 节点前中断

// 第一次执行 — 自动中断
result := cg.Invoke(input)
if result.Interrupted {
    println("中断于:", result.InterruptedAt)
    // 此时可以人工审批，修改 state 等

    // 恢复执行
    result = cg.Resume(updatedInput)
}
```

## Callbacks（回调）

```go
// 标准输出回调（调试用）
cm := callbacks.NewCallbackManager()
cm.AddHandler(&handlers.StdOutHandler{})

// 文件回调
fh, _ := handlers.NewFileCallbackHandler("trace.log")
cm.AddHandler(fh)

// LangSmith 回调
ls := callbacks.NewLangSmithHandler("ls-api-key", "my-project")
cm.AddHandler(ls)
```

## Embeddings（嵌入）

```go
// 文本嵌入
emb := embeddings.NewOpenAIEmbeddings(cfg)
vec, _ := emb.EmbedQuery("今天天气怎么样？")
println(len(vec)) // 1536

// 批量嵌入
vecs, _ := emb.EmbedDocuments([]string{"doc1", "doc2"})

// 余弦相似度
sim := embeddings.CosineSimilarity(vecs[0], vecs[1])

// Mock（测试用）
mock := embeddings.NewMockEmbeddings(128)
```

## Vector Stores（向量存储）

```go
// InMemory（默认）
store := vs.NewInMemoryVectorStore()

// Qdrant（需要 Qdrant 服务）
qdrant, _ := vs.NewQdrantVectorStore(ctx, vs.QdrantConfig{
    URL: "http://192.168.10.100:6334", Collection: "docs", VectorSize: 1536,
})

// MongoDB
mongo, _ := vs.NewMongoVectorStore(ctx, vs.MongoConfig{
    URI: "mongodb://localhost:27017", Database: "mydb", Collection: "vectors",
})
```

## Rate Limiting（限流）

```go
rl := core.NewRateLimiter(10, time.Second) // 每秒 10 次
rl.Wait(ctx) // 阻塞直到拿到令牌
```

## 重试机制

所有 HTTP 调用内置指数退避重试（3 次，1s→2s→4s + 25% 抖动）。

自定义重试：

```go
retryCfg := core.DefaultRetryConfig()
retryCfg.MaxRetries = 5
result, err := core.DoWithRetry(ctx, retryCfg, func(ctx context.Context) (string, error) {
    return httpCall()
})
```

## 配置说明

复制 `config/config.toml.example` 为 `config/config.toml`：

```toml
[openai]
api_key = "sk-xxx"
base_url = "https://dashscope.aliyuncs.com/compatible-mode/v1"
chat_model = "qwen-plus"
embedding_model = "text-embedding-v1"
```

`config.toml` 已加入 `.gitignore`，不会被提交。

## 测试

```bash
# 单元测试（无需 API Key）
go test ./...

# 集成测试（需要 config.toml 中的 API Key）
go test -tags=integration ./tests/integration/ -v

# 端到端测试
go test -tags=e2e ./tests/e2e/ -v
```

---

> **商业使用须知**：本框架为商业私有软件，未经授权不得用于商业用途。完整文档、技术支持和企业授权请联系作者。
