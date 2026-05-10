package openai

import (
	"testing"
	"time"
)

// 娴嬭瘯 DefaultOpenAIConfig锛氶粯璁?BaseURL銆丮odel銆乀imeout 姝ｇ‘
func TestDefaultOpenAIConfig(t *testing.T) {
	cfg := DefaultOpenAIConfig()
	if cfg.BaseURL != defaultBaseURL {
		t.Errorf("expected %s, got %s", defaultBaseURL, cfg.BaseURL)
	}
	if cfg.Model != defaultModel {
		t.Errorf("expected %s, got %s", defaultModel, cfg.Model)
	}
	if cfg.Timeout != 60*time.Second {
		t.Errorf("expected 60s, got %v", cfg.Timeout)
	}
}

// 娴嬭瘯 WithAPIKey 璁剧疆 APIKey
func TestOpenAIConfigWithAPIKey(t *testing.T) {
	cfg := DefaultOpenAIConfig().WithAPIKey("sk-test")
	if cfg.APIKey != "sk-test" {
		t.Errorf("expected sk-test, got %s", cfg.APIKey)
	}
}

// 娴嬭瘯 WithModel 璁剧疆妯″瀷鍚嶇О
func TestOpenAIConfigWithModel(t *testing.T) {
	cfg := DefaultOpenAIConfig().WithModel("gpt-4")
	if cfg.Model != "gpt-4" {
		t.Errorf("expected gpt-4, got %s", cfg.Model)
	}
}

// 娴嬭瘯 WithBaseURL 鑷畾涔?API 鍦板潃
func TestOpenAIConfigWithBaseURL(t *testing.T) {
	cfg := DefaultOpenAIConfig().WithBaseURL("https://custom.api.com")
	if cfg.BaseURL != "https://custom.api.com" {
		t.Errorf("unexpected BaseURL: %s", cfg.BaseURL)
	}
}

// 娴嬭瘯 WithTemperature 璁剧疆娓╁害鍙傛暟
func TestOpenAIConfigWithTemperature(t *testing.T) {
	cfg := DefaultOpenAIConfig().WithTemperature(0.7)
	if cfg.Temperature != 0.7 {
		t.Errorf("expected 0.7, got %f", cfg.Temperature)
	}
}

// 娴嬭瘯 WithMaxTokens 璁剧疆鏈€澶?token 鏁?func TestOpenAIConfigWithMaxTokens(t *testing.T) {
	cfg := DefaultOpenAIConfig().WithMaxTokens(2048)
	if cfg.MaxTokens != 2048 {
		t.Errorf("expected 2048, got %d", cfg.MaxTokens)
	}
}

// 娴嬭瘯閾惧紡璋冪敤锛欰PIKey銆丮odel銆乀emperature銆丮axTokens 鍏ㄩ儴閫氳繃閾惧紡璁剧疆
func TestOpenAIConfigChaining(t *testing.T) {
	cfg := DefaultOpenAIConfig().
		WithAPIKey("sk-key").
		WithModel("gpt-4").
		WithTemperature(0.5).
		WithMaxTokens(1000)

	if cfg.APIKey != "sk-key" || cfg.Model != "gpt-4" || cfg.Temperature != 0.5 || cfg.MaxTokens != 1000 {
		t.Errorf("chained config fields mismatch: %+v", cfg)
	}
}
