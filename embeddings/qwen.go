package embeddings

// QwenEmbeddings uses Qwen's embedding API (OpenAI-compatible).
type QwenEmbeddings struct {
	inner *OpenAIEmbeddings
}

func NewQwenEmbeddings(apiKey string) *QwenEmbeddings {
	cfg := DefaultOpenAIEmbeddingsConfig().
		WithAPIKey(apiKey).
		WithModel("text-embedding-v1")
	cfg.BaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
	return &QwenEmbeddings{inner: NewOpenAIEmbeddings(cfg)}
}

func (q *QwenEmbeddings) EmbedQuery(text string) ([]float32, error)     { return q.inner.EmbedQuery(text) }
func (q *QwenEmbeddings) EmbedDocuments(texts []string) ([][]float32, error) { return q.inner.EmbedDocuments(texts) }
func (q *QwenEmbeddings) Dimension() int    { return q.inner.Dimension() }
func (q *QwenEmbeddings) ModelName() string { return "text-embedding-v1" }
