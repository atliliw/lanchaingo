package embeddings

import "testing"

// 娴嬭瘯 OpenAIEmbeddingsConfig 榛樿鍊?func TestOpenAIDefaultConfig(t *testing.T) {
	cfg := DefaultOpenAIEmbeddingsConfig()
	if cfg.Model != "text-embedding-ada-002" {
		t.Errorf("expected ada-002, got %s", cfg.Model)
	}
	if cfg.BaseURL != "https://api.openai.com/v1" {
		t.Errorf("unexpected base URL: %s", cfg.BaseURL)
	}
}

// 娴嬭瘯 OpenAIEmbeddingsConfig 閾惧紡璁剧疆
func TestOpenAIConfigBuilder(t *testing.T) {
	cfg := DefaultOpenAIEmbeddingsConfig().
		WithAPIKey("sk-test").
		WithModel("text-embedding-3-large")

	if cfg.APIKey != "sk-test" {
		t.Errorf("unexpected API key: %s", cfg.APIKey)
	}
	if cfg.Model != "text-embedding-3-large" {
		t.Errorf("unexpected model: %s", cfg.Model)
	}
}

// 娴嬭瘯 OpenAIEmbeddings 缁村害
func TestOpenAIEmbeddingsDimension(t *testing.T) {
	e := NewOpenAIEmbeddings(DefaultOpenAIEmbeddingsConfig())
	if e.Dimension() != 1536 {
		t.Errorf("expected 1536, got %d", e.Dimension())
	}

	e2 := NewOpenAIEmbeddings(DefaultOpenAIEmbeddingsConfig().WithModel("text-embedding-3-large"))
	if e2.Dimension() != 3072 {
		t.Errorf("expected 3072, got %d", e2.Dimension())
	}
}
