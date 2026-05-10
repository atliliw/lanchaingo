package openai

import "time"

const (
	defaultBaseURL = "https://api.openai.com/v1"
	defaultModel   = "gpt-3.5-turbo"
)

// OpenAIConfig holds configuration for the OpenAI chat client.
type OpenAIConfig struct {
	APIKey      string
	BaseURL     string
	Model       string
	Temperature float32
	MaxTokens   int
	Timeout     time.Duration
}

// DefaultOpenAIConfig returns an OpenAIConfig with sensible defaults.
func DefaultOpenAIConfig() OpenAIConfig {
	return OpenAIConfig{
		BaseURL: defaultBaseURL,
		Model:   defaultModel,
		Timeout: 60 * time.Second,
	}
}

func (c OpenAIConfig) WithAPIKey(key string) OpenAIConfig {
	c.APIKey = key
	return c
}

func (c OpenAIConfig) WithModel(model string) OpenAIConfig {
	c.Model = model
	return c
}

func (c OpenAIConfig) WithBaseURL(url string) OpenAIConfig {
	c.BaseURL = url
	return c
}

func (c OpenAIConfig) WithTemperature(temp float32) OpenAIConfig {
	c.Temperature = temp
	return c
}

func (c OpenAIConfig) WithMaxTokens(max int) OpenAIConfig {
	c.MaxTokens = max
	return c
}
