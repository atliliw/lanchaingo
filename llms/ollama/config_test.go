package ollama

import "testing"

func TestDefaultOllamaConfig(t *testing.T) {
	cfg := DefaultOllamaConfig()
	if cfg.Model != defaultModel { t.Errorf("expected %s, got %s", defaultModel, cfg.Model) }
}
func TestOllamaConfigChaining(t *testing.T) {
	cfg := DefaultOllamaConfig().WithModel("llama3").WithTemperature(0.5)
	if cfg.Model != "llama3" || cfg.Temperature != 0.5 { t.Errorf("config mismatch") }
}
