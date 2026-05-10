package language_models

import "context"

// BaseLanguageModel is the base interface for all language models.
//
// It extends the Runnable interface for unified invocation.
// All LLM wrappers (OpenAI, Ollama, Anthropic, etc.) implement this.
//
// Maps to Rust langchainrust::core::language_models::base::BaseLanguageModel.
type BaseLanguageModel interface {
	// ModelName returns the model identifier (e.g., "gpt-4", "llama3").
	ModelName() string

	// GetNumTokens estimates the number of tokens in the given text.
	GetNumTokens(text string) int

	// Temperature returns the temperature setting, if configured.
	Temperature() *float32

	// MaxTokens returns the max tokens limit, if configured.
	MaxTokens() *int

	// WithTemperature sets the temperature and returns the model for chaining.
	WithTemperature(temp float32) BaseLanguageModel

	// WithMaxTokens sets the max tokens limit and returns the model for chaining.
	WithMaxTokens(max int) BaseLanguageModel
}

// BaseChatModel is the interface for chat-oriented language models.
//
// Accepts a list of messages and returns an LLMResult.
// Supports both synchronous chat and streaming responses.
//
// Maps to Rust langchainrust::core::language_models::chat::BaseChatModel.
type BaseChatModel interface {
	BaseLanguageModel

	// Chat sends a list of messages and returns the model response.
	Chat(ctx context.Context, messages []Message) (*LLMResult, error)

	// StreamChat sends messages and returns a channel yielding response chunks.
	StreamChat(ctx context.Context, messages []Message) (<-chan string, error)

	// BindTools registers tool definitions for function calling.
	BindTools(tools []ToolDefinition)
}
