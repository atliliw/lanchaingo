package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/atliliw/lanchaingo/core"
)

type openAIEmbedReq struct {
	Model string `json:"model"`
	Input any    `json:"input"`
}

type openAIEmbedResp struct {
	Data  []openAIEmbedData `json:"data"`
	Model string            `json:"model"`
}

type openAIEmbedData struct {
	Embedding []float32 `json:"embedding"`
	Index     int       `json:"index"`
	Object    string    `json:"object"`
}

// OpenAIEmbeddingsConfig holds configuration for OpenAI embedding API.
type OpenAIEmbeddingsConfig struct {
	APIKey    string
	BaseURL   string
	Model     string
	BatchSize int
}

func DefaultOpenAIEmbeddingsConfig() OpenAIEmbeddingsConfig {
	return OpenAIEmbeddingsConfig{
		BaseURL:   "https://api.openai.com/v1",
		Model:     "text-embedding-ada-002",
		BatchSize: 2048,
	}
}

func (c OpenAIEmbeddingsConfig) WithAPIKey(key string) OpenAIEmbeddingsConfig {
	c.APIKey = key
	return c
}

func (c OpenAIEmbeddingsConfig) WithModel(model string) OpenAIEmbeddingsConfig {
	c.Model = model
	return c
}

// OpenAIEmbeddings uses the OpenAI embedding API.
type OpenAIEmbeddings struct {
	config    OpenAIEmbeddingsConfig
	client    *http.Client
	dimension int
}

func NewOpenAIEmbeddings(config OpenAIEmbeddingsConfig) *OpenAIEmbeddings {
	dim := 1536
	switch config.Model {
	case "text-embedding-3-large":
		dim = 3072
	case "text-embedding-3-small":
		dim = 1536
	}
	return &OpenAIEmbeddings{
		config:    config,
		client:    &http.Client{},
		dimension: dim,
	}
}

// embed 鍙戦€佸祵鍏ヨ姹傚苟杩斿洖鎵€鏈夊祵鍏ュ悜閲?// input 鍙互鏄?string锛堝崟涓枃鏈級鎴?[]string锛堟壒閲忔枃鏈級
// embed sends embedding request and returns all embedding data with retry
func (e *OpenAIEmbeddings) embed(ctx context.Context, input any) ([]openAIEmbedData, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	url := e.config.BaseURL + "/embeddings"
	reqBody := openAIEmbedReq{Model: e.config.Model, Input: input}
	data, err := json.Marshal(reqBody)
	if err != nil {
		return nil, NewEmbeddingError(ErrParse, "marshal request", err)
	}

	retryCfg := core.DefaultRetryConfig()
	result, err := core.DoWithRetry(ctx, retryCfg, func(ctx context.Context) ([]openAIEmbedData, error) {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(data))
		if err != nil { return nil, NewEmbeddingError(ErrHTTP, "create request", err) }
		httpReq.Header.Set("Authorization", "Bearer "+e.config.APIKey)
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := e.client.Do(httpReq)
		if err != nil { return nil, NewEmbeddingError(ErrHTTP, "request failed", err) }
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			return nil, NewEmbeddingError(ErrAPI, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(body)), nil)
		}

		var embedResp openAIEmbedResp
		if err := json.NewDecoder(resp.Body).Decode(&embedResp); err != nil {
			return nil, NewEmbeddingError(ErrParse, "decode response", err)
		}
		if len(embedResp.Data) == 0 {
			return nil, NewEmbeddingError(ErrAPI, "empty response data", nil)
		}
		return embedResp.Data, nil
	})
	if err != nil { return nil, err }
	return result, nil
}

func (e *OpenAIEmbeddings) EmbedQuery(text string) ([]float32, error) {
	if text == "" {
		return nil, NewEmbeddingError(ErrEmptyInput, "empty input", nil)
	}
	data, err := e.embed(context.Background(), text)
	if err != nil {
		return nil, err
	}
	return data[0].Embedding, nil
}

func (e *OpenAIEmbeddings) EmbedDocuments(texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return nil, nil
	}
	if len(texts) == 1 {
		emb, err := e.EmbedQuery(texts[0])
		if err != nil {
			return nil, err
		}
		return [][]float32{emb}, nil
	}

	data, err := e.embed(context.Background(), texts)
	if err != nil {
		return nil, err
	}

	sort.Slice(data, func(i, j int) bool {
		return data[i].Index < data[j].Index
	})

	results := make([][]float32, len(data))
	for i, d := range data {
		results[i] = d.Embedding
	}
	return results, nil
}

func (e *OpenAIEmbeddings) Dimension() int    { return e.dimension }
func (e *OpenAIEmbeddings) ModelName() string { return e.config.Model }
