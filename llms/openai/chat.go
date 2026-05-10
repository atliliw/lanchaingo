package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	lm "github.com/atliliw/lanchaingo/core/language_models"
)

// OpenAIChat is an OpenAI-compatible chat model implementation.
type OpenAIChat struct {
	config    OpenAIConfig
	client    *http.Client
	bindTools []lm.ToolDefinition
}

// NewOpenAIChat creates a new OpenAIChat with the given config.
func NewOpenAIChat(config OpenAIConfig) *OpenAIChat {
	if config.BaseURL == "" {
		config.BaseURL = defaultBaseURL
	}
	if config.Model == "" {
		config.Model = defaultModel
	}
	if config.Timeout == 0 {
		config.Timeout = 60 * time.Second
	}
	return &OpenAIChat{
		config: config,
		client: &http.Client{Timeout: config.Timeout},
	}
}

func (c *OpenAIChat) ModelName() string {
	return c.config.Model
}

func (c *OpenAIChat) GetNumTokens(text string) int {
	// Simple estimation: ~4 chars per token
	return len(text) / 4
}

func (c *OpenAIChat) Temperature() *float32 {
	if c.config.Temperature == 0 {
		return nil
	}
	return &c.config.Temperature
}

func (c *OpenAIChat) MaxTokens() *int {
	if c.config.MaxTokens == 0 {
		return nil
	}
	return &c.config.MaxTokens
}

func (c *OpenAIChat) WithTemperature(temp float32) lm.BaseLanguageModel {
	c.config.Temperature = temp
	return c
}

func (c *OpenAIChat) WithMaxTokens(max int) lm.BaseLanguageModel {
	c.config.MaxTokens = max
	return c
}

func (c *OpenAIChat) BindTools(tools []lm.ToolDefinition) {
	c.bindTools = tools
}

// requestBody represents the OpenAI chat completion request body.
type requestBody struct {
	Model       string           `json:"model"`
	Messages    []messageBody    `json:"messages"`
	Temperature float32          `json:"temperature,omitempty"`
	MaxTokens   int              `json:"max_tokens,omitempty"`
	Stream      bool             `json:"stream,omitempty"`
	Tools       []toolBody       `json:"tools,omitempty"`
}

type messageBody struct {
	Role       string      `json:"role"`
	Content    string      `json:"content"`
	Name       string      `json:"name,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
	ToolCalls  []toolCallBody `json:"tool_calls,omitempty"`
}

type toolBody struct {
	Type     string        `json:"type"`
	Function functionBody  `json:"function"`
}

type functionBody struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type toolCallBody struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function functionCallBody `json:"function"`
}

type functionCallBody struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// chatResponse represents the OpenAI chat completion response.
type chatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Model   string   `json:"model"`
	Choices []choice `json:"choices"`
	Usage   *usage   `json:"usage,omitempty"`
}

type choice struct {
	Index        int            `json:"index"`
	Message      responseMessage `json:"message"`
	FinishReason string         `json:"finish_reason"`
}

type responseMessage struct {
	Role      string         `json:"role"`
	Content   string         `json:"content"`
	ToolCalls []toolCallBody `json:"tool_calls,omitempty"`
}

type usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// streamChunk represents a single chunk from a streaming response.
type streamChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Model   string        `json:"model"`
	Choices []streamChoice `json:"choices"`
}

type streamChoice struct {
	Index int `json:"index"`
	Delta struct {
		Role      string         `json:"role,omitempty"`
		Content   string         `json:"content,omitempty"`
		ToolCalls []toolCallBody `json:"tool_calls,omitempty"`
	} `json:"delta"`
	FinishReason string `json:"finish_reason"`
}

func (c *OpenAIChat) Chat(ctx context.Context, messages []lm.Message) (*lm.LLMResult, error) {
	req := c.buildRequest(messages, false)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("openai: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.BaseURL+"/chat/completions",
		bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai: failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai: API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var chatResp chatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return nil, fmt.Errorf("openai: failed to decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return nil, fmt.Errorf("openai: empty response choices")
	}

	result := &lm.LLMResult{
		Content: chatResp.Choices[0].Message.Content,
		Model:   chatResp.Model,
	}

	if chatResp.Usage != nil {
		result.TokenUsage = &lm.TokenUsage{
			PromptTokens:     chatResp.Usage.PromptTokens,
			CompletionTokens: chatResp.Usage.CompletionTokens,
			TotalTokens:      chatResp.Usage.TotalTokens,
		}
	}

	if toolCalls := chatResp.Choices[0].Message.ToolCalls; len(toolCalls) > 0 {
		result.ToolCalls = make([]lm.ToolCall, len(toolCalls))
		for i, tc := range toolCalls {
			result.ToolCalls[i] = lm.ToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			}
		}
	}

	return result, nil
}

func (c *OpenAIChat) StreamChat(ctx context.Context, messages []lm.Message) (<-chan string, error) {
	req := c.buildRequest(messages, true)

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("openai: failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.config.BaseURL+"/chat/completions",
		bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("openai: failed to create request: %w", err)
	}

	c.setHeaders(httpReq)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("openai: stream request failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("openai: stream API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	ch := make(chan string)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		for event := range ParseSSE(resp.Body) {
			if isDoneEvent(event.Data) {
				return
			}

			var chunk streamChunk
			if err := json.Unmarshal([]byte(event.Data), &chunk); err != nil {
				continue
			}

			for _, choice := range chunk.Choices {
				if choice.Delta.Content != "" {
					select {
					case ch <- choice.Delta.Content:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch, nil
}

func (c *OpenAIChat) buildRequest(messages []lm.Message, stream bool) requestBody {
	req := requestBody{
		Model:       c.config.Model,
		Stream:      stream,
		Temperature: c.config.Temperature,
	}

	if c.config.MaxTokens > 0 {
		req.MaxTokens = c.config.MaxTokens
	}

	req.Messages = make([]messageBody, len(messages))
	for i, msg := range messages {
		mb := messageBody{
			Content:    msg.Content,
			Role:       messageTypeToRole(msg.MessageType),
			Name:       msg.Name,
			ToolCallID: msg.ToolCallID,
		}
		if len(msg.ToolCalls) > 0 {
			mb.ToolCalls = make([]toolCallBody, len(msg.ToolCalls))
			for j, tc := range msg.ToolCalls {
				mb.ToolCalls[j] = toolCallBody{
					ID:   tc.ID,
					Type: "function",
					Function: functionCallBody{
						Name:      tc.Name,
						Arguments: tc.Arguments,
					},
				}
			}
		}
		req.Messages[i] = mb
	}

	if len(c.bindTools) > 0 {
		req.Tools = make([]toolBody, len(c.bindTools))
		for i, td := range c.bindTools {
			req.Tools[i] = toolBody{
				Type: "function",
				Function: functionBody{
					Name:        td.Name,
					Description: td.Description,
					Parameters:  td.Parameters,
				},
			}
		}
	}

	return req
}

func (c *OpenAIChat) setHeaders(req *http.Request) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)
}

func messageTypeToRole(mt lm.MessageType) string {
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
