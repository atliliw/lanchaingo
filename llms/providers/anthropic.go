package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/atliliw/lanchaingo/core"
	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// AnthropicChat implements Claude API chat.
type AnthropicChat struct {
	APIKey       string
	BaseURL      string
	Model        string
	maxToken     int
	client       *http.Client
}

func NewAnthropicChat(apiKey string) *AnthropicChat {
	return &AnthropicChat{
		APIKey:    apiKey,
		BaseURL:   "https://api.anthropic.com/v1",
		Model:     "claude-3-5-sonnet-20241022",
		maxToken: 4096,
		client:    &http.Client{Timeout: 60 * time.Second},
	}
}

type anthropicReq struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []anthropicMsg  `json:"messages"`
}

type anthropicMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResp struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Model string `json:"model"`
}

func (c *AnthropicChat) ModelName() string { return c.Model }
func (c *AnthropicChat) GetNumTokens(text string) int { return len(text) / 4 }
func (c *AnthropicChat) Temperature() *float32 { return nil }
func (c *AnthropicChat) MaxTokens() *int { return &c.maxToken }
func (c *AnthropicChat) WithTemperature(f float32) lm.BaseLanguageModel { return c }
func (c *AnthropicChat) WithMaxTokens(i int) lm.BaseLanguageModel { c.maxToken = i; return c }

func (c *AnthropicChat) Chat(ctx context.Context, msgs []messages.Message) (*lm.LLMResult, error) {
	req := anthropicReq{
		Model:     c.Model,
		MaxTokens: c.maxToken,
	}
	for _, m := range msgs {
		role := "user"
		if m.MessageType == messages.MessageTypeAI { role = "assistant" }
		if m.MessageType == messages.MessageTypeSystem { role = "user" }
		req.Messages = append(req.Messages, anthropicMsg{Role: role, Content: m.Content})
	}

	body, _ := json.Marshal(req)
	retryCfg := core.DefaultRetryConfig()
	result, err := core.DoWithRetry(ctx, retryCfg, func(ctx context.Context) (*lm.LLMResult, error) {
		httpReq, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/messages", bytes.NewReader(body))
		if err != nil { return nil, err }
		httpReq.Header.Set("x-api-key", c.APIKey)
		httpReq.Header.Set("anthropic-version", "2023-06-01")
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(httpReq)
		if err != nil { return nil, fmt.Errorf("anthropic: request failed: %w", err) }
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("anthropic: status %d: %s", resp.StatusCode, string(respBody))
		}

		var ar anthropicResp
		if err := json.Unmarshal(respBody, &ar); err != nil {
			return nil, fmt.Errorf("anthropic: parse failed: %w", err)
		}
		content := ""
		if len(ar.Content) > 0 { content = ar.Content[0].Text }
		return &lm.LLMResult{Content: content, Model: ar.Model}, nil
	})
	if err != nil { return nil, err }
	return result, nil
}

func (c *AnthropicChat) StreamChat(ctx context.Context, msgs []messages.Message) (<-chan string, error) {
	var antMsgs []anthropicMsg
	for _, m := range msgs {
		role := "user"
		if m.MessageType == messages.MessageTypeAI { role = "assistant" }
		antMsgs = append(antMsgs, anthropicMsg{Role: role, Content: m.Content})
	}
	req := struct {
		Model         string        `json:"model"`
		MaxTokens     int           `json:"max_tokens"`
		Messages      []anthropicMsg `json:"messages"`
		Stream        bool          `json:"stream"`
	}{Model: c.Model, MaxTokens: c.maxToken, Messages: antMsgs, Stream: true}

	body, _ := json.Marshal(req)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/messages", bytes.NewReader(body))
	httpReq.Header.Set("x-api-key", c.APIKey)
	httpReq.Header.Set("anthropic-version", "2023-06-01")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("anthropic: stream failed: %w", err)
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") { continue }
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" { return }
			var event struct {
				Type string `json:"type"`
				Delta struct {
					Text string `json:"text"`
				} `json:"delta"`
			}
			if err := json.Unmarshal([]byte(data), &event); err != nil { continue }
			if event.Type == "content_block_delta" && event.Delta.Text != "" {
				ch <- event.Delta.Text
			}
		}
	}()
	return ch, nil
}

func (c *AnthropicChat) BindTools(_ []lm.ToolDefinition) {}
