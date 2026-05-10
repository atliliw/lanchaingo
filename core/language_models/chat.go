package language_models

import (
	"time"

	"github.com/atliliw/lanchaingo/core"
)

// LLMResult is the result of a chat model invocation.
// Re-exported from core for convenience.
type LLMResult = core.LLMResult

// TokenUsage tracks token consumption.
// Re-exported from core for convenience.
type TokenUsage = core.TokenUsage

// ToolCall represents a tool invocation.
// Re-exported from core for convenience.
type ToolCall = core.ToolCall

// ToolDefinition defines a tool for function calling.
// Re-exported from core for convenience.
type ToolDefinition = core.ToolDefinition

// ChatStream wraps the streaming response from a chat model.
type ChatStream struct {
	Data chan string
	Err  chan error
}

// NewChatStream creates a ChatStream with buffered channels.
func NewChatStream(bufferSize int) *ChatStream {
	if bufferSize <= 0 {
		bufferSize = 100
	}
	return &ChatStream{
		Data: make(chan string, bufferSize),
		Err:  make(chan error, 1),
	}
}

// Close safely closes the data channel.
func (cs *ChatStream) Close() {
	// Multiple closes of a channel panic, so use recover.
	defer func() { recover() }()
	close(cs.Data)
	close(cs.Err)
}

// ModelConfig provides common configuration for LLM providers.
type ModelConfig struct {
	APIKey      string
	BaseURL     string
	Model       string
	Temperature float32
	MaxTokens   int
	Timeout     time.Duration
}

// Validate checks that required fields are set.
func (c *ModelConfig) Validate() error {
	if c.APIKey == "" {
		return ErrMissingAPIKey
	}
	if c.Model == "" {
		return ErrMissingModel
	}
	return nil
}
