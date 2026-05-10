# langchaingo �?1:1 复刻 langchainrust 商业版计划书

> **目标**: 使用 Go 语言 1:1 复刻 [atliliw/langchainrust](https://github.com/atliliw/langchainrust) 的全部功能，构建生产�?LLM 应用框架�?> **许可**: 商业私有（自主版权，不采�?MIT/Apache 等开源许可证�?> **状�?*: 📋 计划阶段

---

## 一、项目概�?
### 1.1 什么是 langchaingo

langchaingo �?langchainrust �?Go 语言 1:1 复刻版。提供一套完整的 LLM 应用开发框架，涵盖�?LLM 调用、Agent、RAG、BM25 混合搜索、LangGraph 工作流到多向量存�?文档存储后端的全链路能力�?
### 1.2 为什么是 1:1 复刻

| 维度 | 说明 |
|------|------|
| **API 语义** | 保持�?langchainrust 一致的组件命名、参数结构、调用模�?|
| **模块划分** | 完全对应 src/ 下的每个模块目录 |
| **功能覆盖�?* | 100% 覆盖所�?Feature（LLM/Agent/Memory/Chain/RAG/BM25/Hybrid/LangGraph/VectorStore/Tools/OutputParser/Callbacks/Cache�?|
| **测试体系** | 单元测试 + 集成测试 + 端到端测试三层架�?|

> **不采�?* LangChain 官方 Python SDK 的设计，而是严格对标 langchainrust 的实现风格（Rust �?trait + async + 类型安全模式映射�?Go �?interface + goroutine + 泛型）�?
---

## 二、架构总览

```
┌──────────────────────────────────────────────────────────────────────�?�?                        langchaingo                                   �?├──────────────────────────────────────────────────────────────────────�?�? LLM Layer (language_models �?llms/)                                �?�? ├── OpenAIChat / OllamaChat                                         �?�? ├── DeepSeek / Moonshot / Zhipu / Qwen (OpenAI 兼容)               �?�? ├── AnthropicChat (Claude API)                                      �?�? ├── GeminiChat                                                      �?�? ├── Function Calling (BindTools)                                    �?�? └── Streaming (StreamChat)                                         �?├──────────────────────────────────────────────────────────────────────�?�? Embeddings Layer                                                    �?�? ├── OpenAIEmbeddings / DeepSeekEmbeddings                           �?�? └── QwenEmbeddings / MockEmbeddings                                 �?├──────────────────────────────────────────────────────────────────────�?�? Agent Layer                                                         �?�? ├── ReActAgent / FunctionCallingAgent                               �?�? ├── AgentExecutor                                                   �?�? └── LangGraph (StateGraph, Subgraph, Parallel)                      �?├──────────────────────────────────────────────────────────────────────�?�? Retrieval Layer                                                     �?�? ├── RAG (TextSplitter, VectorStore)                                 �?�? ├── BM25 (Keyword Search, AutoMerging)                              �?�? ├── Hybrid (BM25 + Vector, RRF Fusion)                              �?�? └── HyDE / MultiQuery / Reranking                                   �?�? └── Storage (InMemory, Qdrant, MongoDB)                             �?├──────────────────────────────────────────────────────────────────────�?�? Utility Layer                                                       �?�? ├── Memory (Buffer, Window, Summary, SummaryBuffer, Persistent)     �?�? ├── Chains (LLMChain, SequentialChain, ConversationChain,           �?�? �?        RouterChain, RetrievalQA, ConversationRetrieval,          �?�? �?        Stuff, Refine, MapReduce)                                 �?�? ├── Prompts (PromptTemplate, ChatPromptTemplate, FewShot)           �?�? ├── Tools (Calculator, DateTime, URLFetch, Wikipedia,               �?�? �?        WebSearch, PythonREPL)                                    �?�? ├── Output Parsers (Str, Json, CommaSeparated, Structured, Typed)   �?�? ├── LLM Cache (TTL 支持)                                            �?�? ├── Callbacks (StdOut, LangSmith, FileHandler)                      �?�? └── Document Loaders (PDF, CSV, Text, JSON, Markdown)               �?└──────────────────────────────────────────────────────────────────────�?```

---

## 三、目录结构（Go 实现�?
> ⚠️ **目录 1:1 对齐说明**：以下结构与 langchainrust v0.2.20 `src/` 完全对应。Go 包名遵循 Go 惯例（短名称），括号内标注了对应�?Rust 模块路径�?
```
langchaingo/
�?├── core/                          # Rust: core/ �?核心抽象层（Go interface + 子包�?�?  ├── types.go                  # 公共类型：LLMResult, TokenUsage, AgentAction �?�?  ├── runnable.go               # Runnable interface (invoke/batch/stream)
�?  ├── runnable_config.go        # Rust: core/runnables/config.rs �?RunnableConfig
�?  ├── language_models/          # Rust: core/language_models/ �?LLM 抽象接口
�?  �?  ├── base.go               # BaseLanguageModel interface
�?  �?  └── chat.go               # BaseChatModel interface
�?  ├── tools/                    # Rust: core/tools/ �?工具抽象
�?  �?  ├── base.go               # BaseTool / Tool interface
�?  �?  ├── registry.go           # Rust: registry.rs �?ToolRegistry
�?  �?  ├── tool_definition.go    # Rust: tool_definition.rs �?ToolDefinition / ToolCall
�?  �?  ├── structured.go         # Rust: structured.rs �?结构化工�?�?  �?  └── structured_output.go  # Rust: structured_output.rs �?StructuredOutput
�?  ├── output_parsers/           # Rust: core/output_parsers/
�?  �?  ├── base.go               # BaseOutputParser
�?  �?  ├── string.go             # StrOutputParser
�?  �?  ├── json.go               # JsonOutputParser
�?  �?  ├── list.go               # CommaSeparatedListOutputParser
�?  �?  ├── structured.go         # StructuredOutputParser
�?  �?  └── typed.go              # TypedOutputParser (泛型)
�?  └── cache/                    # Rust: core/cache/
�?      ├── llm_cache.go          # LLMCache interface
�?      └── ttl_cache.go          # TTL 缓存实现（内存）
�?├── llms/                          # Rust: language_models/ �?LLM 实现
�?  ├── openai/                   # Rust: language_models/openai/
�?  �?  ├── chat.go               # OpenAIChat
�?  �?  ├── config.go             # OpenAIConfig
�?  �?  └── sse.go                # SSE 流式解析
�?  ├── ollama/                   # Rust: language_models/ollama/
�?  �?  ├── chat.go               # OllamaChat
�?  �?  └── config.go             # OllamaConfig
�?  ├── providers/                # Rust: language_models/providers/
�?  �?  ├── deepseek.go           # DeepSeekChat (OpenAI 兼容)
�?  �?  ├── moonshot.go           # MoonshotChat
�?  �?  ├── zhipu.go              # ZhipuChat
�?  �?  ├── qwen.go               # QwenChat
�?  �?  ├── anthropic.go          # AnthropicChat (Claude API)
�?  �?  └── gemini.go             # GeminiChat
�?  └── model_config.go           # ModelConfig 枚举/公共配置
�?├── schema/                        # Rust: schema/ �?消息/数据结构
�?  └── messages/                 # Rust: schema/messages/
�?      ├── message.go            # Message / SystemMessage / HumanMessage / AIMessage / ToolMessage
�?      └── mod.go                # 模块导出
�?├── embeddings/                    # Rust: embeddings/
�?  ├── openai.go                 # OpenAIEmbeddings
�?  ├── deepseek.go               # DeepSeekEmbeddings
�?  ├── qwen.go                   # QwenEmbeddings
�?  └── mock.go                   # MockEmbeddings
�?├── agents/                        # Rust: agents/
�?  ├── base.go                   # BaseAgent interface / AgentAction / AgentError
�?  ├── types.go                  # Agent types / AgentStep / AgentFinish
�?  ├── executor.go               # AgentExecutor
�?  ├── react/                    # Rust: agents/react/
�?  �?  ├── agent.go              # ReActAgent
�?  �?  ├── parser.go             # ReAct 输出解析
�?  �?  └── prompt.go             # ReAct 提示词模�?�?  └── function_calling/         # Rust: agents/function_calling/
�?      └── agent.go              # FunctionCallingAgent
�?├── tools/                         # Rust: tools/ �?内置工具
�?  ├── calculator.go             # Calculator 计算�?�?  ├── math.go                   # SimpleMathTool
�?  ├── datetime.go               # DateTimeTool 日期时间
�?  ├── url_fetch.go              # URLFetchTool
�?  ├── wikipedia.go              # WikipediaTool 搜索
�?  ├── search.go                 # DuckDuckGoSearchTool 搜索引擎
�?  └── python_repl.go            # PythonREPLTool
�?├── memory/                        # Rust: memory/
�?  ├── base.go                   # BaseMemory / ChatMessageHistory
�?  ├── buffer.go                 # ConversationBufferMemory
�?  ├── window.go                 # ConversationBufferWindowMemory
�?  ├── summary.go                # ConversationSummaryMemory
�?  ├── summary_buffer.go         # ConversationSummaryBufferMemory
�?  ├── persistent.go             # PersistentMemory
�?  └── mongo_memory.go           # MongoPersistentMemory (build tag: mongodb-persistence)
�?├── prompts/                       # Rust: prompts/
�?  ├── prompt_template.go        # PromptTemplate
�?  ├── chat_prompt_template.go   # ChatPromptTemplate
�?  └── few_shot.go               # FewShotPromptTemplate
�?├── chains/                        # Rust: chains/
�?  ├── base.go                   # BaseChain interface / ChainError / ChainResult
�?  ├── llm_chain.go              # LLMChain + LLMChainBuilder
�?  ├── sequential_chain.go       # SequentialChain
�?  ├── conversation_chain.go     # ConversationChain
�?  ├── router_chain.go           # RouterChain / LLMRouterChain
�?  ├── retrieval_qa.go           # RetrievalQA
�?  ├── conversation_retrieval.go # ConversationRetrievalChain
�?  └── document_chains.go        # Stuff / Refine / MapReduce / MapRerank 文档�?�?├── retrieval/                     # Rust: retrieval/ �?RAG 检索组�?�?  ├── retriever.go              # Retriever interface / SimilarityRetriever
�?  ├── splitter.go               # TextSplitter / RecursiveCharacterSplitter / FixedSize / Regex
�?  ├── hyde.go                   # HyDE (Hypothetical Document Embeddings)
�?  ├── multi_query.go            # MultiQuery Retriever
�?  ├── reranking.go              # Reranking / KeywordReranker / BM25Reranker
�?  ├── hybrid.go                 # HybridRetriever (BM25 + Vector, RRF Fusion)
�?  ├── chunked_hybrid.go         # ChunkedHybridRetriever
�?  ├── unified_hybrid.go         # UnifiedHybridIndex
�?  ├── bm25/                     # BM25 关键词搜�?�?  �?  ├── bm25.go               # BM25Retriever / BM25Params
�?  �?  ├── algorithm.go          # BM25 核心算法
�?  �?  ├── index.go              # BM25Index
�?  �?  ├── tokenizer.go          # 中英文分词器
�?  �?  └── chunked.go            # ChunkedBM25Retriever
�?  └── loaders/                  # Rust: retrieval/loaders/ �?文档加载�?�?      ├── pdf.go                # PDF Loader
�?      ├── csv.go                # CSV Loader
�?      ├── text.go               # Text Loader
�?      ├── json.go               # JSON Loader
�?      └── markdown.go           # Markdown Loader
�?├── vector_stores/                 # Rust: vector_stores/
�?  ├── provider.go               # VectorStoreProvider / VectorStoreType
�?  ├── in_memory.go              # InMemoryVectorStore
�?  ├── document_store.go         # DocumentStore interface
�?  ├── chunked_vector_store.go   # ChunkedDocumentStore
�?  ├── chromadb.go               # ChromaDBVectorStore
�?  ├── qdrant.go                 # QdrantVectorStore (build tag: qdrant-integration)
�?  ├── mongo_document_store.go   # MongoChunkedDocumentStore (build tag: mongodb-persistence)
�?  ├── redis_store.go            # RedisDocumentStore (build tag: redis-storage)
�?  └── sqlite_store.go           # SQLiteDocumentStore (build tag: sqlite-storage)
�?├── callbacks/                     # Rust: callbacks/
�?  ├── base.go                   # CallbackHandler interface
�?  ├── callback_manager.go       # CallbackManager
�?  ├── run_tree.go               # RunTree / RunType 执行追踪
�?  ├── langsmith_client.go       # LangSmith HTTP 客户�?�?  └── handlers/                 # Rust: callbacks/handlers/
�?      ├── stdout.go             # StdOutHandler
�?      ├── langsmith.go          # LangSmithHandler
�?      └── file.go               # FileCallbackHandler
�?├── langgraph/                     # Rust: langgraph/ �?LangGraph 工作�?�?  ├── graph.go                  # StateGraph / GraphBuilder
�?  ├── compiled.go               # CompiledGraph
�?  ├── node.go                   # GraphNode / AsyncNode / NodeConfig
�?  ├── edge.go                   # GraphEdge / ConditionalEdge
�?  ├── state.go                  # StateSchema / Reducer / AgentState
�?  ├── subgraph.go               # SubgraphNode / SubgraphBuilder
�?  ├── checkpointer.go           # Checkpointer / MemoryCheckpointer / FileCheckpointer
�?  ├── persistence.go            # GraphPersistence / MemoryPersistence / FilePersistence
�?  └── errors.go                 # GraphError
�?├── tests/                         # 测试（与 Rust tests/ 1:1 镜像�?�?  ├── unit/                     # 单元测试（无外部依赖�?�?  �?  ├── schema_test.go
�?  �?  ├── prompts_test.go
�?  �?  ├── tools_test.go
�?  �?  ├── memory_test.go
�?  �?  ├── retrieval_test.go
�?  �?  ├── vectorstores_test.go
�?  �?  ├── callbacks_test.go
�?  �?  ├── output_parsers_test.go
�?  �?  ├── bm25_test.go
�?  �?  ├── multi_query_test.go
�?  �?  ├── hyde_test.go
�?  �?  ├── reranking_test.go
�?  �?  ├── chromadb_test.go
�?  �?  ├── sqlite_store_test.go
�?  �?  ├── gemini_test.go
�?  �?  ├── few_shot_test.go
�?  �?  ├── providers_test.go
�?  �?  ├── conversation_chain_test.go
�?  �?  ├── router_chain_test.go
�?  �?  ├── retrieval_qa_test.go
�?  �?  ├── conversation_retrieval_test.go
�?  �?  ├── summary_memory_test.go
�?  �?  ├── summary_buffer_memory_test.go
�?  �?  ├── parallel_tool_calls_test.go
�?  �?  ├── file_handler_test.go
�?  �?  ├── runnable_stream_test.go
�?  �?  ├── ollama_test.go
�?  �?  └── tool_calling_test.go
�?  ├── integration/              # 集成测试（需�?API Key�?�?  �?  ├── agent_memory_test.go
�?  �?  ├── agent_test.go
�?  �?  ├── agent_react_test.go
�?  �?  ├── chains_test.go
�?  �?  ├── embeddings_test.go
�?  �?  ├── llm_chat_test.go
�?  �?  ├── rag_test.go
�?  �?  ├── rag_full_test.go
�?  �?  ├── callbacks_integration_test.go
�?  �?  ├── callbacks_llm_test.go
�?  �?  ├── tool_callbacks_test.go
�?  �?  ├── ollama_chat_test.go
�?  �?  ├── langsmith_test.go
�?  �?  ├── chain_workflow_test.go
�?  �?  ├── rag_pipeline_test.go
�?  �?  └── bm25_real_docs_test.go
�?  ├── e2e/                      # 端到端测�?�?  �?  └── full_application_test.go
�?  ├── function_calling/         # Function Calling 测试
�?  �?  └── agent_fc_test.go
�?  ├── langgraph/                 # LangGraph 测试
�?  �?  ├── basic_test.go
�?  �?  ├── conditional_test.go
�?  �?  ├── state_test.go
�?  �?  ├── recursion_test.go
�?  �?  ├── edge_test.go
�?  �?  ├── checkpointer_test.go
�?  �?  ├── example_test.go
�?  �?  ├── subgraph_test.go
�?  �?  ├── persistence_test.go
�?  �?  ├── parallel_test.go
�?  �?  ├── openai_demo_test.go
�?  �?  └── subgraph_mechanism_test.go
�?  ├── bm25/                      # BM25 专项测试
�?  �?  ├── chunked_test.go
�?  �?  ├── rag_test.go
�?  �?  ├── llm_integration_test.go
�?  �?  └── hybrid_rag_test.go
�?  ├── loaders/                   # Loader 测试
�?  �?  ├── text_loader_test.go
�?  �?  ├── json_loader_test.go
�?  �?  └── markdown_loader_test.go
�?  └── hyde_reranking/           # HyDE & Reranking 测试
�?      ├── hyde_test.go
�?      └── reranking_test.go
�?├── go.mod                         # Go 模块定义
├── go.sum
├── Makefile                       # 构建/测试/ lint 命令
└── README.md                      # 项目说明
```

---

## 四、Rust �?Go 映射策略

### 4.1 核心概念映射

| Rust 概念 | Go 等价实现 |
|-----------|------------|
| `trait` | `interface{}` |
| `async fn` / `async_trait` | `goroutine` + `channel` + `context.Context` |
| `enum` | `type` + `iota` + 接口分派 |
| `impl Struct` | `(s *Struct) Method()` |
| `derive(Serialize, Deserialize)` | `json:"field_name"` 标签 |
| `Result<T, E>` | `(T, error)` |
| `Option<T>` | `*T` �?`T` + `ok bool` |
| `Box<dyn Trait>` | `interface{}` 变量 |
| `Arc<Mutex<T>>` | `sync.Mutex` / `sync.RWMutex` |
| `#[derive(Clone)]` | 自定�?`Clone()` 方法 |
| feature flags (Cargo) | build tags + 接口分派 |
| `Stream` trait | `chan` 流式输出 |
| `#[serde(skip_serializing_if = "...")]` | `json:",omitempty"` |
| `schemars (JSON Schema 生成)` | `github.com/invopop/jsonschema` |
| `impl<T: Into<...>>` | Go 1.22+ 泛型约束 |

### 4.2 流式处理（Rust Stream �?Go chan�?
```go
// Rust: async fn stream_chat(...) -> impl Stream<Item = ChatResponse>
// Go:
type ChatStream struct {
    Data chan ChatResponse
    Err  chan error
}

func (c *OpenAIChat) StreamChat(ctx context.Context, messages []Message) *ChatStream {
    stream := &ChatStream{
        Data: make(chan ChatResponse),
        Err:  make(chan error, 1),
    }
    go func() {
        defer close(stream.Data)
        // HTTP SSE 请求，�?block 写入 channel
    }()
    return stream
}
```

### 4.3 Runnable 接口统一

```go
// Rust: pub trait Runnable { async fn invoke(...) -> Output; async fn batch(...) -> Vec<Output>; async fn stream(...) -> impl Stream; }
// Go:
type Runnable interface {
    Invoke(ctx context.Context, input map[string]any) (any, error)
    Batch(ctx context.Context, inputs []map[string]any) ([]any, error)
    Stream(ctx context.Context, input map[string]any) (<-chan any, error)
}
```

### 4.4 Tool 接口统一

```go
// Rust: pub trait Tool { fn name(&self) -> String; fn description(&self) -> String; async fn run(&self, input: &str) -> String; }
// Go:
type Tool interface {
    Name() string
    Description() string
    Run(ctx context.Context, input string) (string, error)
}

// ToolRegistry �?对应 Rust core/tools/registry.rs
type ToolRegistry struct {
    tools map[string]Tool
}
func NewToolRegistry() *ToolRegistry
func (r *ToolRegistry) Register(tool Tool)
func (r *ToolRegistry) Get(name string) (Tool, bool)
func (r *ToolRegistry) ToDefinitions() []ToolDefinition
```

### 4.5 错误处理

```go
// Rust: thiserror 派生 + enum AgentError
// Go: 自定�?error 类型 + errors.Is / errors.As
type AgentError struct {
    Kind    AgentErrorKind
    Message string
    Cause   error
}

type AgentErrorKind int
const (
    ErrUnknown AgentErrorKind = iota
    ErrToolExecution
    ErrLLMResponse
    ErrMaxIterations
    ErrInvalidState
)
```

---

## 五、实现阶段（Phase 计划�?
### Phase 1：核心基础设施 + 基础 LLM�?周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| core/ | 所有核�?interface + types + runnable + runnable_config | P0 |
| core/language_models/ | BaseLanguageModel, BaseChatModel 接口 | P0 |
| core/tools/ | BaseTool, Tool, ToolRegistry, ToolDefinition, StructuredOutput | P0 |
| core/cache/ | LLMCache 接口 | P1 |
| core/output_parsers/ | BaseOutputParser 接口 | P1 |
| schema/messages/ | Message, SystemMessage, HumanMessage, AIMessage | P0 |
| llms/openai/ | OpenAIChat（流�?+ 非流�?+ BindTools�?| P0 |
| llms/ollama/ | OllamaChat | P0 |
| prompts/ | PromptTemplate, ChatPromptTemplate | P0 |
| go.mod | 模块初始�?+ 基础依赖 | P0 |

**交付�?*: 可用�?`OpenAIChat` + `OllamaChat`，支�?Chat/StreamChat/BindTools

### Phase 2：Memory + Chains�?周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| memory/ | BufferMemory, WindowMemory, SummaryMemory, SummaryBufferMemory, SimpleMemory | P0 |
| chains/ | LLMChain, SequentialChain, ConversationChain, RouterChain | P0 |

**交付�?*: 完整的记忆系�?+ 基础链式调用（对�?Rust `memory/` + `chains/`�?
### Phase 3：Agent 系统�?周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| agents/ | BaseAgent interface, AgentTypes, AgentExecutor | P0 |
| agents/react/ | ReActAgent, Parser, Prompt | P0 |
| agents/function_calling/ | FunctionCallingAgent | P0 |
| tools/ | Calculator, Math, DateTime, URLFetch | P0 |
| tools/ | Wikipedia, Search(DuckDuckGo), PythonREPL | P1 |

**交付�?*: 完整�?Agent 系统，支�?ReAct + Function Calling + 多工�?
### Phase 4：RAG 检索系统（2周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| embeddings/ | OpenAIEmbeddings, DeepSeekEmbeddings, QwenEmbeddings, MockEmbeddings | P0 |
| vector_stores/ | InMemoryVectorStore, DocumentStore interface | P0 |
| retrieval/ | Retriever interface, SimilarityRetriever | P0 |
| retrieval/ | TextSplitter (RecursiveCharacter, FixedSize, Regex) | P0 |
| retrieval/ | HyDE, MultiQuery, Reranking | P1 |
| retrieval/loaders/ | PDF Loader, CSV Loader, Text Loader, JSON Loader, Markdown Loader | P1 |

**交付�?*: 完整�?RAG 管线，支持文档加载→分割→嵌入→检�?
### Phase 5：BM25 + 混合检索（1周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| retrieval/bm25/ | BM25Retriever, Algorithm, Index, Tokenizer(中英�?, ChunkedBM25 | P0 |
| retrieval/ | HybridRetriever (BM25+Vector RRF Fusion) | P0 |
| retrieval/ | ChunkedHybridRetriever, UnifiedHybridIndex | P1 |

**交付�?*: BM25 关键词搜�?+ BM25/Vector 混合检�?+ 统一混合索引

### Phase 6：LangGraph 工作流（2周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| langgraph/ | StateGraph, GraphBuilder, CompiledGraph | P0 |
| langgraph/ | GraphNode, GraphEdge, ConditionalEdge, State | P0 |
| langgraph/ | Subgraph (SubgraphNode, SubgraphBuilder) | P0 |
| langgraph/ | Checkpointer (MemoryCheckpointer, FileCheckpointer) | P1 |
| langgraph/ | Persistence (MemoryPersistence, FilePersistence) | P1 |
| langgraph/ | GraphError | P0 |

**交付�?*: 完整�?LangGraph 工作流引擎（对应 Rust `langgraph/` �?10 个文件）

### Phase 7：高�?Chain + 向量存储后端�?周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| chains/ | RetrievalQA, ConversationRetrievalChain | P0 |
| chains/document_chains.go | StuffDocumentsChain, RefineDocumentsChain, MapReduceDocumentsChain | P0 |
| vector_stores/ | Qdrant, MongoDB, ChromaDB, Redis, SQLite + ChunkedDocumentStore | P1 |
| core/output_parsers/ | StrOutputParser, JsonOutputParser, ListParser, Structured, Typed | P0 |
| core/cache/ | LLMCache with TTL | P1 |
| callbacks/ | CallbackManager, RunTree, StdOutHandler, LangSmith(handler+client), FileHandler | P1 |

**交付�?*: 高级�?+ 多后端向量存�?+ 输出解析 + 回调系统

### Phase 8：剩�?LLM Provider + 测试全覆盖（1周）

| 模块 | 内容 | 优先�?|
|------|------|--------|
| llms/providers/ | DeepSeek, Moonshot, Zhipu, Qwen, Anthropic, Gemini | P1 |
| tests/ | 对齐 langchainrust 全部测试用例（Rust tests/ 下所有文件） | P0 |

**交付�?*: 全部 LLM Provider + 测试通过

---

## 六、Go 依赖选择

| 功能 | Go �?| 说明 |
|------|-------|------|
| HTTP 客户�?| `net/http` (标准�? | 无需第三�?|
| JSON | `encoding/json` (标准�? | 对标 serde_json |
| 流式 SSE | `bufio` + `net/http` (标准�? | 自行实现 SSE 解析 |
| UUID | `github.com/google/uuid` | 对标 uuid crate |
| 正则 | `regexp` (标准�? | 对标 regex crate |
| 日期时间 | `time` (标准�? | 对标 chrono |
| PDF 解析 | `github.com/ledongthuc/pdf` | 对标 pdf-extract |
| CSV | `encoding/csv` (标准�? | 对标 csv crate |
| MongoDB | `go.mongodb.org/mongo-driver` | feature: mongodb-persistence |
| Qdrant | `github.com/qdrant/go-client` | feature: qdrant-integration |
| Redis | `github.com/redis/go-redis` | feature: redis-storage |
| SQLite | `modernc.org/sqlite` (�?Go) | feature: sqlite-storage |
| 向量运算 | `gonum.org/v1/gonum` | 可选，用于 BM25 数学 |
| JSON Schema 生成 | `github.com/invopop/jsonschema` | 对标 schemars，用�?Tool Calling | |
| 日志 | `log/slog` (标准�? Go 1.21+) | 对标 log crate |

> **原则**: 优先标准库，最小化第三方依赖。所有外部存储后端通过 optional build tags 控制�?
---

## 七、关键设计决�?
### 7.1 go.mod 模块设计

```
module github.com/atliliw/lanchaingo  // TODO: 替换为你的私有仓库路�?
go 1.22

// 核心依赖 - 始终引入
require (
    github.com/google/uuid v1.6.0
    github.com/ledongthuc/pdf v0.0.0-20240201131950-da5b75571b58
    github.com/invopop/jsonschema v0.12.0   // JSON Schema 生成，对�?schemars
)

// 可选依�?- build tags 控制
// go build -tags qdrant,mongodb,redis,sqlite
```

### 7.2 并发模型

```go
// 所�?LLM 调用接受 context.Context 支持超时/取消
type BaseChatModel interface {
    Chat(ctx context.Context, messages []Message) (*ChatResult, error)
    Stream(ctx context.Context, messages []Message) (<-chan ChatResult, error)
}

// Agent 使用 goroutine 执行工具调用
func (e *AgentExecutor) Execute(ctx context.Context, input string) (string, error) {
    // 主循环在单个 goroutine 中运�?    // 工具调用可通过 errgroup 并行执行
}
```

### 7.3 配置模式

```go
// Rust: OpenAIConfig { api_key, base_url, model, ..Default::default() }
// Go: 函数选项模式 + WithXxx

type OpenAIConfig struct {
    APIKey  string
    BaseURL string
    Model   string
    // ... 其他参数
}

func NewOpenAIChat(config OpenAIConfig) *OpenAIChat {
    if config.BaseURL == "" {
        config.BaseURL = "https://api.openai.com/v1"
    }
    if config.Model == "" {
        config.Model = "gpt-3.5-turbo"
    }
    return &OpenAIChat{config: config}
}

// 或使�?Builder 模式
chat := NewOpenAIChat().
    WithAPIKey(os.Getenv("OPENAI_API_KEY")).
    WithModel("gpt-4").
    WithBaseURL("https://api.openai.com/v1").
    Build()
```

### 7.4 构建标签控制可选功�?
```go
// openai.go - 始终编译
// qdrant.go
//go:build qdrant

package vector_stores

import qdrant "github.com/qdrant/go-client/qdrant"

type QdrantVectorStore struct { ... }
```

```makefile
# Makefile
build:
    go build -tags "qdrant mongodb redis sqlite" ./...

test-unit:
    go test -tags "qdrant mongodb redis sqlite" ./tests/unit/...

test-all:
    go test -tags "qdrant mongodb redis sqlite" ./...
```

### 7.5 命名规范对齐

| langchainrust | langchaingo |
|---------------|-------------|
| `OpenAIChat` | `OpenAIChat` |
| `ReActAgent` | `ReActAgent` |
| `PromptTemplate::new("...")` | `NewPromptTemplate("...")` |
| `llm.chat(vec![...])` | `llm.Chat(ctx, messages)` |
| `agent.run(input).await` | `agent.Run(ctx, input)` |
| `DuckDuckGoSearchTool` | `DuckDuckGoSearchTool` | |
| `InMemoryVectorStore` | `InMemoryVectorStore` |

> **命名原则**: Go 导出大写开头，避免 stutter。保持与 source 的语义一致性�?
---

## 八、测试策�?
### 8.1 三层测试体系

| 层级 | 目的 | 执行条件 |
|------|------|----------|
| 单元测试 `tests/unit/` | 验证单个组件逻辑 | 无外部依赖，`go test` 即可运行 |
| 集成测试 `tests/integration/` | 验证组件间协�?| 需�?API Key / 外部服务 |
| 端到端测�?`tests/e2e/` | 验证完整应用流程 | 需要完整环�?|

### 8.2 Mock 策略

```go
// 使用 MockEmbeddings 避免真实 API 调用
embeddings := &MockEmbeddings{Dimension: 384}
docs := []core.Document{...}
vecs, _ := embeddings.EmbedDocuments(ctx, docs)
// 测试向量存储/检索逻辑

// 使用 MockLLM 测试 Agent 逻辑
llm := &MockLLM{
    Responses: map[string]string{
        "think": "I need to use the calculator tool",
        "act":   "Action: calculator\nActionInput: 2+2",
    },
}
```

### 8.3 覆盖率目�?
| 模块 | 目标覆盖�?|
|------|-----------|
| core/ | 90%+ |
| schema/ | 95%+ |
| prompts/ | 90%+ |
| memory/ | 90%+ |
| tools/ | 85%+ |
| output_parsers/ | 90%+ |
| agents/ (核心逻辑) | 85%+ |
| chains/ | 85%+ |
| retrieval/ | 85%+ |
| llms/ | 80%+ (需 Mock) |
| vector_stores/ | 80%+ (需 Mock) |

---

## 九、项目管�?
### 9.1 时间线估�?
| Phase | 内容 | 预估人周 | 依赖 |
|-------|------|---------|------|
| 1 | 核心基础设施 + 基础 LLM | 2 人周 | �?|
| 2 | Memory + Chains | 1 人周 | Phase 1 |
| 3 | Agent 系统 | 2 人周 | Phase 1, 2 |
| 4 | RAG 检索系�?| 2 人周 | Phase 1 |
| 5 | BM25 + 混合检�?| 1 人周 | Phase 4 |
| 6 | LangGraph 工作�?| 2 人周 | Phase 1, 3 |
| 7 | 高级 Chain + 向量存储后端 | 2 人周 | Phase 4 |
| 8 | 剩余 LLM Provider + 测试 | 1 人周 | Phase 1 |
| **合计** | | **13 人周** | |

> �?3 人团队并行开发，�?5-6 周可完成全部功能�?
### 9.2 依赖关系�?
```
Phase 1 (Core + LLM)
    ├── Phase 2 (Memory + Chains)
    �?  └── Phase 3 (Agent)
    �?      └── Phase 6 (LangGraph)
    └── Phase 4 (RAG)
        ├── Phase 5 (BM25 + Hybrid)
        └── Phase 7 (高级 Chain + 向量存储)
Phase 8 (剩余 Provider + 测试) - 贯穿全程
```

### 9.3 质量门禁

- [ ] 所有单元测试通过
- [ ] `go vet ./...` 无错�?- [ ] `golangci-lint run ./...` 无严�?lint 错误
- [ ] 核心模块覆盖�?�?85%
- [ ] API 文档完整（每个导出函数有注释�?- [ ] �?`as any` 等效�?unsafe 类型绕过

---

## 十、商业版注意事项

### 10.1 许可合规

- langchainrust 采用 MIT OR Apache-2.0 许可
- **langchaingo 作为 1:1 复刻，代码完全自主编�?*，不使用 langchainrust 的源码（仅参考接口设计）
- 建议商标和命名避免与 LangChain 官方混淆（`langchaingo` 作为产品名需确认可使用）

### 10.2 私有仓库配置

```bash
# 使用私有 Go module
go env -w GOPRIVATE=github.com/atliliw/lanchaingo

# go.mod
module github.com/atliliw/lanchaingo

# .gitignore
.env
*.key
*.pem
```

### 10.3 生产环境增强�?
| 增强�?| 说明 | 优先�?|
|--------|------|--------|
| 请求重试 + 指数退�?| 所�?LLM 调用自动重试 | P0 |
| 速率限制 | Provider 级别限流 | P1 |
| 请求追踪 | OpenTelemetry 集成 | P1 |
| 熔断�?| 连续失败自动熔断 | P2 |
| 配置热加�?| 动态更�?LLM 配置 | P2 |
| Prometheus 指标 | 请求延迟/成功�?Token 消�?| P1 |
| 结构化日�?| slog 结构化输�?| P0 |

---

## 十一、langchainrust v0.2.20 Feature 清单（检查表�?
### LLM Provider�?�?- [ ] OpenAIChat（流�?+ 非流�?+ SSE�?- [ ] OllamaChat
- [ ] DeepSeekChat
- [ ] MoonshotChat
- [ ] ZhipuChat
- [ ] QwenChat
- [ ] AnthropicChat
- [ ] GeminiChat

### Embeddings�?�?- [ ] OpenAIEmbeddings
- [ ] DeepSeekEmbeddings
- [ ] QwenEmbeddings
- [ ] MockEmbeddings

### Agents�?�?- [ ] ReActAgent（含 Parser + Prompt�?- [ ] FunctionCallingAgent
- [ ] AgentExecutor
- [ ] ToolRegistry（工具注册中心）

### Memory�?�?- [ ] BufferMemory（ConversationBufferMemory�?- [ ] WindowMemory（ConversationBufferWindowMemory�?- [ ] SummaryMemory（ConversationSummaryMemory�?- [ ] SummaryBufferMemory（ConversationSummaryBufferMemory�?- [ ] PersistentMemory
- [ ] MongoPersistentMemory（build tag: mongodb-persistence�?
### Chains�?�?- [ ] LLMChain（含 LLMChainBuilder�?- [ ] SequentialChain
- [ ] ConversationChain
- [ ] RouterChain / LLMRouterChain
- [ ] RetrievalQA
- [ ] ConversationRetrievalChain
- [ ] StuffDocumentsChain
- [ ] RefineDocumentsChain
- [ ] MapReduceDocumentsChain（含 MapRerankDocumentsChain�?
### RAG / Retrieval�?0�?- [ ] RecursiveCharacterSplitter
- [ ] FixedSizeSplitter
- [ ] RegexSplitter
- [ ] SimilarityRetriever
- [ ] RerankerRetriever（KeywordReranker + BM25Reranker�?- [ ] HyDE（Hypothetical Document Embeddings�?- [ ] MultiQuery Retriever
- [ ] Reranking（RerankingExecutor�?- [ ] ChunkedHybridRetriever
- [ ] UnifiedHybridIndex

### BM25�?�?- [ ] BM25Retriever + BM25Params
- [ ] BM25 核心算法（Algorithm�?- [ ] BM25Index
- [ ] 中英�?Tokenizer
- [ ] ChunkedBM25Retriever

### Hybrid�?�?- [ ] BM25 + Vector Hybrid（HybridRetriever�?- [ ] RRF Fusion（reciprocal_rank_fusion�?- [ ] ChunkedHybridRetriever

### LangGraph�?0�?- [ ] StateGraph / GraphBuilder
- [ ] CompiledGraph
- [ ] GraphNode / AsyncNode
- [ ] GraphEdge / ConditionalEdge
- [ ] StateSchema / Reducer（Replace, Append, AppendMessages, AppendSteps�?- [ ] Subgraph（SubgraphNode / SubgraphBuilder�?- [ ] Parallel Execution（ParallelInvocation�?- [ ] Checkpointer（MemoryCheckpointer / FileCheckpointer�?- [ ] Persistence（MemoryPersistence / FilePersistence�?- [ ] Time-travel Debugging

### Vector Stores�?�?- [ ] VectorStoreProvider + VectorStoreType
- [ ] DocumentStore 接口
- [ ] InMemoryVectorStore
- [ ] ChunkedDocumentStore（InMemoryChunkedDocumentStore�?- [ ] QdrantVectorStore（build tag: qdrant-integration�?- [ ] MongoChunkedDocumentStore（build tag: mongodb-persistence�?- [ ] ChromaDBVectorStore
- [ ] RedisDocumentStore（build tag: redis-storage�?- [ ] SQLiteDocumentStore（build tag: sqlite-storage�?
### Tools�?�?- [ ] CalculatorTool
- [ ] SimpleMathTool
- [ ] DateTimeTool
- [ ] URLFetchTool
- [ ] WikipediaTool
- [ ] DuckDuckGoSearchTool
- [ ] PythonREPLTool

### Output Parsers�?�?- [ ] StrOutputParser
- [ ] JsonOutputParser
- [ ] CommaSeparatedListOutputParser
- [ ] StructuredOutputParser
- [ ] TypedOutputParser

### Callbacks�?�?- [ ] CallbackHandler 接口
- [ ] CallbackManager
- [ ] RunTree / RunType（执行追踪）
- [ ] StdOutHandler
- [ ] LangSmithHandler + LangSmithClient
- [ ] FileCallbackHandler

### Cache�?�?- [ ] LLMCache with TTL

### Document Loaders�?�?- [ ] PDF Loader
- [ ] CSV Loader
- [ ] Text Loader
- [ ] JSON Loader
- [ ] Markdown Loader

---

**总计: 92 �?Feature 模块**（含子模块细分），全�?1:1 复刻 langchainrust v0.2.20�?
---

## 十二、开始开�?
### 初始化项�?
```bash
cd D:\BaiduNetdiskDownload\LLM\langchaingo
go mod init github.com/atliliw/lanchaingo   # TODO: 替换为你的私有仓库路�?mkdir -p core/language_models core/tools core/output_parsers core/cache schema/messages llms/openai llms/ollama llms/providers prompts memory agents/react agents/function_calling tools embeddings vector_stores retrieval/bm25 retrieval/loaders callbacks/handlers langgraph
```

### 第一个文件：core/types.go

所有接口定义从这里开始，是所有上层模块的基础�?
---

**本计划书完整定义�?langchaingo 的商业复刻路线，覆盖 langchainrust v0.2.20 全部 92 �?Feature 模块，并提供�?Go 语言特定的实现策略、依赖选择、并发模型和测试方案�?*

---

## 十三、商业竞争力增强（对标 Python LangChain）

> 以下功能是 Python LangChain 具备但尚未实现的，按商业重要度排序。

### 13.1 结构化输出（P0 — 竞争差异点）

对标 LangChain `model.with_structured_output(Schema)`，允许将 Go 结构体绑定到 Chat Model。

```go
type Person struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}
result, err := llm.WithStructuredOutput[Person]().
    Invoke(ctx, "张三今年30岁")
```

### 13.2 Human-in-the-Loop（P0 — 企业客户刚需）

LangGraph 执行中支持中断和恢复，用于人工审批。

```go
cg := graph.Compile()
cg.InterruptBefore("approve")
inv, _ := cg.Invoke(input)
// inv.Interrupted == true
inv, _ = cg.Resume(approvalResult)
```

### 13.3 PostgreSQL Checkpointer（P0 — 生产部署）

```go
cp := langgraph.NewPostgresCheckpointer("postgres://user:pass@host/db")
cg := graph.Compile().WithCheckpointer(cp)
cg.Resume(nil, "thread-123")
```

### 13.4 Rate Limiting（P1 — 生产稳定性）

```go
rl := core.NewRateLimiter(10, time.Second)
llm := openai.NewChat(config).WithRateLimiter(rl)
```

### 13.5 Store 长期记忆（P1）

```go
store := memory.NewInMemoryStore()
store.Put("user_123", "preferences", map[string]any{"lang": "zh"})
```

---

## 十四、版本计划

### v1.0（当前已完成）
| 模块 | 状态 |
|------|------|
| 全部 92 个 langchainrust Feature | ✅ |
| LangGraph StateGraph | ✅ |
| BM25 + Hybrid 检索 | ✅ |
| Agent ReAct + FunctionCalling | ✅ |
| 8 个 LLM Provider | ✅ |
| 重试中间件（指数退避） | ✅ |
| 集成测试（真实 API） | ✅ |
| Qdrant/MongoDB/Redis/SQLite 后端 | ✅ |
| LangSmith + MongoMemory | ✅ |

### v1.1（本次实现）
| 模块 | 状态 |
|------|------|
| 结构化输出 WithStructuredOutput | ⬜ |
| Human-in-the-Loop | ⬜ |
| PostgreSQL Checkpointer | ⬜ |
| Rate Limiting | ⬜ |
| Store 长期记忆 | ⬜ |

