package ollama

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	lm "github.com/atliliw/lanchaingo/core/language_models"
)

// OllamaChat is a chat model implementation for Ollama.
type OllamaChat struct {
	config OllamaConfig
	client *http.Client
}

// NewOllamaChat creates a new OllamaChat with the given config.
func NewOllamaChat(config OllamaConfig) *OllamaChat {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	if config.Model == "" {
		config.Model = defaultModel
	}
	if config.Timeout == 0 {
		config.Timeout = 120 * time.Second
	}
	return &OllamaChat{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
	}
}

func (c *OllamaChat) ModelName() string {
	return c.config.Model
}

func (c *OllamaChat) GetNumTokens(text string) int {
	return len(text) / 4
}

func (c *OllamaChat) Temperature() *float32 {
	if c.config.Temperature == 0 {
		return nil
	}
	return &c.config.Temperature
}

func (c *OllamaChat) MaxTokens() *int {
	if c.config.MaxTokens == 0 {
		return nil
	}
	return &c.config.MaxTokens
}

func (c *OllamaChat) WithTemperature(temp float32) lm.BaseLanguageModel {
	c.config.Temperature = temp
	return c
}

func (c *OllamaChat) WithMaxTokens(max int) lm.BaseLanguageModel {
	c.config.MaxTokens = max
	return c
}

func (c *OllamaChat) BindTools(_ []lm.ToolDefinition) {
	// Ollama does not natively support OpenAI-style tool calling.
	// Tools can be handled via prompt engineering at the agent level.
}

// ollamaMessage represents a message in Ollama's API format.
type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// chatRequest is the Ollama /api/chat request body.
type chatRequest struct {
	Model    string           `json:"model"`
	Messages []ollamaMessage  `json:"messages"`
	Stream   bool             `json:"stream"`
	Options  *chatOptions     `json:"options,omitempty"`
}

type chatOptions struct {
	Temperature float32 `json:"temperature,omitempty"`
	NumPredict  int     `json:"num_predict,omitempty"`
}

// chatResponse is the Ollama /api/chat non-streaming response.
type chatResponse struct {
	Model     string         `json:"model"`
	CreatedAt string         `json:"created_at"`
	Message   ollamaMessage  `json:"message"`
	Done      bool           `json:"done"`
}

// streamResponse is the Ollama /api/chat streaming response chunk.
type streamResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Message   struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Done bool `json:"done"`
}

func (c *OllamaChat) Chat(ctx context.Context, messages []lm.Message) (*lm.LLMResult, error) {
	req := c.buildRequest(messages, false)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.BaseURL+"/api/chat",
		bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ollama: API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("ollama: failed to decode response: %w", err)
	}

	return &lm.LLMResult{
		Content: chatResp.Message.Content,
		Model:   chatResp.Model,
	}, nil
}

func (c *OllamaChat) StreamChat(ctx context.Context, messages []lm.Message) (<-chan string, error) {
	req := c.buildRequest(messages, true)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.BaseURL+"/api/chat",
		bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("ollama: failed to create stream request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("ollama: stream request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("ollama: stream API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}

			var chunk streamResponse
			if err := json.Unmarshal([]byte(line), &chunk); err != nil {
				continue
			}

			if chunk.Message.Content != "" {
				select {
				case ch <- chunk.Message.Content:
				case <-ctx.Done():
					return
				}
			}

			if chunk.Done {
				return
			}
		}
	}()

	return ch, nil
}

func (c *OllamaChat) buildRequest(messages []lm.Message, stream bool) chatRequest {
	req := chatRequest{
		Model:  c.config.Model,
		Stream: stream,
	}

	req.Messages = make([]ollamaMessage, len(messages))
	for i, msg := range messages {
		req.Messages[i] = ollamaMessage{
			Role:    ollamaRole(msg.MessageType),
			Content: msg.Content,
		}
	}

	if c.config.Temperature > 0 || c.config.MaxTokens > 0 {
		req.Options = &chatOptions{
			Temperature: c.config.Temperature,
			NumPredict:  c.config.MaxTokens,
		}
	}

	return req
}

func ollamaRole(mt lm.MessageType) string {
	switch mt {
	case lm.MessageTypeSystem:
		return "system"
	case lm.MessageTypeHuman:
		return "user"
	case lm.MessageTypeAI:
		return "assistant"
	case lm.MessageTypeTool:
		return "tool"
	default:
		return "user"
	}
}
