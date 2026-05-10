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

	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// OpenAICompatChat is an OpenAI-compatible chat provider.
type OpenAICompatChat struct {
	APIKey   string
	BaseURL  string
	Model    string
	client   *http.Client
}

func NewDeepSeekChat(apiKey string) *OpenAICompatChat {
	return &OpenAICompatChat{
		APIKey:  apiKey,
		BaseURL: "https://api.deepseek.com/v1",
		Model:   "deepseek-chat",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func NewMoonshotChat(apiKey string) *OpenAICompatChat {
	return &OpenAICompatChat{
		APIKey:  apiKey,
		BaseURL: "https://api.moonshot.cn/v1",
		Model:   "moonshot-v1-8k",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func NewZhipuChat(apiKey string) *OpenAICompatChat {
	return &OpenAICompatChat{
		APIKey:  apiKey,
		BaseURL: "https://open.bigmodel.cn/api/paas/v4",
		Model:   "glm-4",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

func NewQwenChat(apiKey string) *OpenAICompatChat {
	return &OpenAICompatChat{
		APIKey:  apiKey,
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
		Model:   "qwen-plus",
		client:  &http.Client{Timeout: 60 * time.Second},
	}
}

type oaiReqMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type oaiReq struct {
	Model    string      `json:"model"`
	Messages []oaiReqMsg `json:"messages"`
	Stream   bool        `json:"stream,omitempty"`
}

type oaiChoice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

type oaiResp struct {
	Model   string     `json:"model"`
	Choices []oaiChoice `json:"choices"`
}

func (c *OpenAICompatChat) ModelName() string { return c.Model }
func (c *OpenAICompatChat) GetNumTokens(text string) int { return len(text) / 4 }
func (c *OpenAICompatChat) Temperature() *float32 { return nil }
func (c *OpenAICompatChat) MaxTokens() *int { return nil }
func (c *OpenAICompatChat) WithTemperature(f float32) lm.BaseLanguageModel { return c }
func (c *OpenAICompatChat) WithMaxTokens(i int) lm.BaseLanguageModel { return c }

func (c *OpenAICompatChat) Chat(ctx context.Context, msgs []messages.Message) (*lm.LLMResult, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var oaiMsgs []oaiReqMsg
	for _, m := range msgs {
		role := "user"
		if m.MessageType == messages.MessageTypeAI { role = "assistant" }
		if m.MessageType == messages.MessageTypeSystem { role = "system" }
		oaiMsgs = append(oaiMsgs, oaiReqMsg{Role: role, Content: m.Content})
	}

	body, _ := json.Marshal(oaiReq{Model: c.Model, Messages: oaiMsgs})
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: request failed: %w", c.Model, err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var oai oaiResp
	if err := json.Unmarshal(respBody, &oai); err != nil {
		return nil, fmt.Errorf("%s: parse failed: %w", c.Model, err)
	}

	content := ""
	if len(oai.Choices) > 0 {
		content = oai.Choices[0].Message.Content
	}
	return &lm.LLMResult{Content: content, Model: oai.Model}, nil
}

func (c *OpenAICompatChat) StreamChat(ctx context.Context, msgs []messages.Message) (<-chan string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	var oaiMsgs []oaiReqMsg
	for _, m := range msgs {
		role := "user"
		if m.MessageType == messages.MessageTypeAI { role = "assistant" }
		if m.MessageType == messages.MessageTypeSystem { role = "system" }
		oaiMsgs = append(oaiMsgs, oaiReqMsg{Role: role, Content: m.Content})
	}

	body, _ := json.Marshal(oaiReq{Model: c.Model, Messages: oaiMsgs, Stream: true})
	httpReq, _ := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewReader(body))
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("%s: stream failed: %w", c.Model, err)
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}
			var chunk struct {
				Choices []struct {
					Delta struct {
						Content string `json:"content"`
					} `json:"delta"`
				} `json:"choices"`
			}
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				continue
			}
			for _, c := range chunk.Choices {
				if c.Delta.Content != "" {
					ch <- c.Delta.Content
				}
			}
		}
	}()
	return ch, nil
}

func (c *OpenAICompatChat) BindTools(_ []lm.ToolDefinition) {}
