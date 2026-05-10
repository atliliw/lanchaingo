package core

import "testing"

func TestNewRunnableConfig(t *testing.T) {
	cfg := NewRunnableConfig()
	if cfg == nil { t.Fatal("expected non-nil config") }
	if cfg.Metadata == nil { t.Error("expected non-nil Metadata map") }
}
func TestRunnableConfigWithTag(t *testing.T) {
	cfg := NewRunnableConfig().WithTag("test-tag")
	if len(cfg.Tags) != 1 || cfg.Tags[0] != "test-tag" { t.Errorf("expected [test-tag], got %v", cfg.Tags) }
}
func TestTokenUsage(t *testing.T) {
	tu := TokenUsage{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30}
	if tu.PromptTokens != 10 { t.Error("TokenUsage fields mismatch") }
}
func TestLLMResult(t *testing.T) {
	r := LLMResult{Content: "hello", Model: "gpt-4"}
	if r.Content != "hello" { t.Error("LLMResult fields mismatch") }
}
func TestNewDocument(t *testing.T) {
	doc := NewDocument("test content")
	if doc.PageContent != "test content" { t.Errorf("expected test content") }
	if doc.Metadata == nil { t.Error("expected non-nil Metadata") }
}
func TestCacheEntryExpired(t *testing.T) {
	e := CacheEntry{Value: "hello"}
	if e.IsExpired() { t.Error("entry without expiry should not be expired") }
}
func TestRunnableConfigWithRunName(t *testing.T) {
	cfg := NewRunnableConfig().WithRunName("my-run")
	if cfg.RunName != "my-run" { t.Errorf("expected my-run") }
}
func TestToolCallRoundTrip(t *testing.T) {
	tc := ToolCall{ID: "call_123", Name: "calculator", Arguments: `{"a":1}`}
	if tc.ID != "call_123" { t.Errorf("ToolCall fields mismatch") }
}
func TestAgentAction(t *testing.T) {
	aa := AgentAction{Tool: "search", ToolInput: "query", Log: "thinking..."}
	if aa.Tool != "search" { t.Errorf("AgentAction fields mismatch") }
}
