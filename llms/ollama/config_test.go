package ollama

import (
	"testing"
	"time"
)

// 娴嬭瘯 DefaultOllamaConfig锛氶粯璁?BaseURL銆丮odel銆乀imeout 姝ｇ‘
func TestDefaultOllamaConfig(t *testing.T) {
	cfg := DefaultOllamaConfig()
	if cfg.BaseURL != defaultBaseURL {
		t.Errorf("expected %s, got %s", defaultBaseURL, cfg.BaseURL)
	}
	if cfg.Model != defaultModel {
		t.Errorf("expected %s, got %s", defaultModel, cfg.Model)
	}
	if cfg.Timeout != 120*time.Second {
		t.Errorf("expected 120s, got %v", cfg.Timeout)
	}
}

// 娴嬭瘯 WithModel 璁剧疆妯″瀷鍚嶇О
func TestOllamaConfigWithModel(t *testing.T) {
	cfg := DefaultOllamaConfig().WithModel("llama3.1")
	if cfg.Model != "llama3.1" {
		t.Errorf("expected llama3.1, got %s", cfg.Model)
	}
}

// 娴嬭瘯 WithBaseURL 鑷畾涔?Ollama 鏈嶅姟鍦板潃
func TestOllamaConfigWithBaseURL(t *testing.T) {
	cfg := DefaultOllamaConfig().WithBaseURL("http://ollama:11434")
	if cfg.BaseURL != "http://ollama:11434" {
		t.Errorf("unexpected BaseURL: %s", cfg.BaseURL)
	}
}

// 娴嬭瘯 WithTemperature 璁剧疆娓╁害鍙傛暟
func TestOllamaConfigWithTemperature(t *testing.T) {
	cfg := DefaultOllamaConfig().WithTemperature(0.8)
	if cfg.Temperature != 0.8 {
		t.Errorf("expected 0.8, got %f", cfg.Temperature)
	}
}

// 娴嬭瘯 WithMaxTokens 璁剧疆鏈€澶?token 鏁?func TestOllamaConfigWithMaxTokens(t *testing.T) {
	cfg := DefaultOllamaConfig().WithMaxTokens(4096)
	if cfg.MaxTokens != 4096 {
		t.Errorf("expected 4096, got %d", cfg.MaxTokens)
	}
}

// 娴嬭瘯閾惧紡璋冪敤锛歁odel銆丅aseURL銆乀emperature銆丮axTokens 鍏ㄩ儴閫氳繃閾惧紡璁剧疆
func TestOllamaConfigChaining(t *testing.T) {
	cfg := DefaultOllamaConfig().
		WithModel("mistral").
		WithBaseURL("http://localhost:11434").
		WithTemperature(0.3).
		WithMaxTokens(500)

	if cfg.Model != "mistral" || cfg.BaseURL != "http://localhost:11434" {
		t.Errorf("chained config fields mismatch: %+v", cfg)
	}
}
