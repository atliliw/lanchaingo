package ollama

import "time"

const (
	defaultBaseURL = "http://localhost:11434"
	defaultModel   = "llama3"
)

// OllamaConfig holds configuration for the Ollama chat client.
type OllamaConfig struct {
	BaseURL     string
	Model       string
	Temperature float32
	MaxTokens   int
	Timeout     time.Duration
}

// DefaultOllamaConfig returns an OllamaConfig with sensible defaults.
func DefaultOllamaConfig() OllamaConfig {
	return OllamaConfig{
		BaseURL: defaultBaseURL,
		Model:   defaultModel,
		Timeout: 120 * time.Second,
	}
}

func (c OllamaConfig) WithModel(model string) OllamaConfig {
	c.Model = model
	return c
}

func (c OllamaConfig) WithBaseURL(url string) OllamaConfig {
	c.BaseURL = url
	return c
}

func (c OllamaConfig) WithTemperature(temp float32) OllamaConfig {
	c.Temperature = temp
	return c
}

func (c OllamaConfig) WithMaxTokens(max int) OllamaConfig {
	c.MaxTokens = max
	return c
}
