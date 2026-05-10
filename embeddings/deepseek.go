package embeddings

// DeepSeekEmbeddings uses DeepSeek's embedding API (OpenAI-compatible).
type DeepSeekEmbeddings struct {
	inner *OpenAIEmbeddings
}

func NewDeepSeekEmbeddings(apiKey string) *DeepSeekEmbeddings {
	cfg := DefaultOpenAIEmbeddingsConfig().
		WithAPIKey(apiKey).
		WithModel("deepseek-embed")
	cfg.BaseURL = "https://api.deepseek.com/v1"
	return &DeepSeekEmbeddings{inner: NewOpenAIEmbeddings(cfg)}
}

func (d *DeepSeekEmbeddings) EmbedQuery(text string) ([]float32, error)   { return d.inner.EmbedQuery(text) }
func (d *DeepSeekEmbeddings) EmbedDocuments(texts []string) ([][]float32, error) { return d.inner.EmbedDocuments(texts) }
func (d *DeepSeekEmbeddings) Dimension() int  { return d.inner.Dimension() }
func (d *DeepSeekEmbeddings) ModelName() string { return "deepseek-embed" }
