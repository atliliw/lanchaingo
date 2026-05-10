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

// GeminiChat implements Google Gemini API chat.
type GeminiChat struct {
	APIKey  string
	BaseURL string
	Model   string
	client  *http.Client
}

func NewGeminiChat(apiKey string) *GeminiChat {
	return &GeminiChat{
		APIKey:  apiKey,
		BaseURL: "https://generativelanguage.googleapis.com/v1beta",
		Model:   "gemini-2.0-flash",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

type geminiReq struct {
	Contents []geminiContent `json:"contents"`
}

type geminiContent struct {
	Role  string      `json:"role"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiResp struct {
	Candidates []struct {
		Content geminiContent `json:"content"`
	} `json:"candidates"`
}

func (c *GeminiChat) ModelName() string { return c.Model }
func (c *GeminiChat) GetNumTokens(text string) int { return len(text) / 4 }
func (c *GeminiChat) Temperature() *float32 { return nil }
func (c *GeminiChat) MaxTokens() *int { return nil }
func (c *GeminiChat) WithTemperature(f float32) lm.BaseLanguageModel { return c }
func (c *GeminiChat) WithMaxTokens(i int) lm.BaseLanguageModel { return c }

func (c *GeminiChat) Chat(ctx context.Context, msgs []messages.Message) (*lm.LLMResult, error) {
	req := geminiReq{}
	for _, m := range msgs {
		role := "user"
		if m.MessageType == messages.MessageTypeAI { role = "model" }
		req.Contents = append(req.Contents, geminiContent{
			Role:  role,
			Parts: []geminiPart{{Text: m.Content}},
		})
	}

	body, _ := json.Marshal(req)
	retryCfg := core.DefaultRetryConfig()
	result, err := core.DoWithRetry(ctx, retryCfg, func(ctx context.Context) (*lm.LLMResult, error) {
		url := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.BaseURL, c.Model, c.APIKey)
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
		if err != nil { return nil, err }
		httpReq.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(httpReq)
		if err != nil { return nil, fmt.Errorf("gemini: request failed: %w", err) }
		defer resp.Body.Close()

		respBody, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			return nil, fmt.Errorf("gemini: status %d: %s", resp.StatusCode, string(respBody))
		}

		var gr geminiResp
		if err := json.Unmarshal(respBody, &gr); err != nil {
			return nil, fmt.Errorf("gemini: parse failed: %w", err)
		}
		content := ""
		if len(gr.Candidates) > 0 && len(gr.Candidates[0].Content.Parts) > 0 {
			content = gr.Candidates[0].Content.Parts[0].Text
		}
		return &lm.LLMResult{Content: content, Model: c.Model}, nil
	})
	if err != nil { return nil, err }
	return result, nil
}

func (c *GeminiChat) StreamChat(ctx context.Context, msgs []messages.Message) (<-chan string, error) {
	var contents []geminiContent
	for _, m := range msgs {
		role := "user"
		if m.MessageType == messages.MessageTypeAI { role = "model" }
		contents = append(contents, geminiContent{Role: role, Parts: []geminiPart{{Text: m.Content}}})
	}

	body, _ := json.Marshal(geminiReq{Contents: contents})
	url := fmt.Sprintf("%s/models/%s:streamGenerateContent?key=%s&alt=sse", c.BaseURL, c.Model, c.APIKey)
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("gemini: stream failed: %w", err)
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
			var gr geminiResp
			if err := json.Unmarshal([]byte(data), &gr); err != nil { continue }
			if len(gr.Candidates) > 0 && len(gr.Candidates[0].Content.Parts) > 0 {
				ch <- gr.Candidates[0].Content.Parts[0].Text
			}
		}
	}()
	return ch, nil
}

func (c *GeminiChat) BindTools(_ []lm.ToolDefinition) {}
