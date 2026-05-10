package language_models

import "errors"

// Sentinel errors for language model operations.
var (
	ErrMissingAPIKey = errors.New("API key is required")
	ErrMissingModel  = errors.New("model name is required")
	ErrChatFailed    = errors.New("chat invocation failed")
	ErrStreamFailed  = errors.New("stream invocation failed")
	ErrNoTools       = errors.New("tool calling requested but no tools bound")
)
