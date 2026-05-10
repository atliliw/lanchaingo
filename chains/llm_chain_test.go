package chains

import (
	"context"
	"testing"
	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

type mockLLM struct{}
func (m *mockLLM) ModelName() string { return "mock" }
func (m *mockLLM) GetNumTokens(string) int { return 0 }
func (m *mockLLM) Temperature() *float32 { return nil }
func (m *mockLLM) MaxTokens() *int { return nil }
func (m *mockLLM) WithTemperature(float32) lm.BaseLanguageModel { return m }
func (m *mockLLM) WithMaxTokens(int) lm.BaseLanguageModel { return m }
func (m *mockLLM) Chat(_ context.Context, _ []messages.Message) (*lm.LLMResult, error) {
	return &lm.LLMResult{Content: "response"}, nil
}
func (m *mockLLM) StreamChat(_ context.Context, _ []messages.Message) (<-chan string, error) {
	ch := make(chan string); close(ch); return ch, nil
}
func (m *mockLLM) BindTools(_ []lm.ToolDefinition) {}
type mockChain struct{ in, out string; fn func(string) string }
func (c *mockChain) InputKeys() []string { return []string{c.in} }
func (c *mockChain) OutputKeys() []string { return []string{c.out} }
func (c *mockChain) Name() string { return "mock" }
func (c *mockChain) ValidateInputs(m map[string]any) error { _, ok := m[c.in]; if !ok { return NewChainError(ErrMissingInput, "", nil) }; return nil }
func (c *mockChain) Invoke(m map[string]any) (ChainResult, error) {
	return ChainResult{c.out: c.fn(m[c.in].(string))}, nil
}

func TestLLMChainInvoke(t *testing.T) {
	chain := NewLLMChain(&mockLLM{}, "template")
	r, _ := chain.Invoke(map[string]any{"input": "hi"})
	if r["text"] != "response" { t.Errorf("got %v", r["text"]) }
}
func TestLLMChainMissingInput(t *testing.T) {
	chain := NewLLMChain(&mockLLM{}, "template")
	_, err := chain.Invoke(map[string]any{})
	if err == nil { t.Error("expected error") }
}
func TestLLMChainCustomKeys(t *testing.T) {
	chain := NewLLMChain(&mockLLM{}, "template").WithInputKey("q").WithOutputKey("a")
	r, _ := chain.Invoke(map[string]any{"q": "hi"})
	if r["a"] != "response" { t.Errorf("got %v", r["a"]) }
}
func TestSequentialChain(t *testing.T) {
	u := &mockChain{in:"t", out:"u", fn: func(s string) string { return s+"_u" }}
	d := &mockChain{in:"u", out:"r", fn: func(s string) string { return s+"_d" }}
	c := NewSequentialChain().AddChain(u, []string{"t"}, []string{"u"}).AddChain(d, []string{"u"}, []string{"r"})
	r, _ := c.Invoke(map[string]any{"t": "hello"})
	if r["r"] != "hello_u_d" { t.Errorf("got %v", r["r"]) }
}
func TestStuffDocumentsChain(t *testing.T) {
	c := NewStuffDocumentsChain(&mockLLM{}, "summarize:")
	r, _ := c.Invoke(map[string]any{"input_documents": []string{"doc1","doc2"}})
	if r["output_text"] != "response" { t.Errorf("got %v", r["output_text"]) }
}
func TestMapReduceDocumentsChain(t *testing.T) {
	c := NewMapReduceDocumentsChain(&mockLLM{}, "map:", "reduce:")
	r, _ := c.Invoke(map[string]any{"input_documents": []string{"a","b"}})
	if r["output_text"] != "response" { t.Errorf("got %v", r["output_text"]) }
}
func TestRouterChainKeywordMatch(t *testing.T) {
	mc := &mockChain{in:"input", out:"output", fn: func(s string) string { return "m:"+s }}
	dc := &mockChain{in:"input", out:"output", fn: func(s string) string { return "d:"+s }}
	r := NewRouterChain().AddRouteWithKeywords("m", "m", mc, []string{"math"}).WithDefault(dc)
	result, _ := r.Invoke(map[string]any{"input": "math problem"})
	if result["output"] != "m:math problem" { t.Errorf("got %v", result["output"]) }
}
func TestRouterChainDefault(t *testing.T) {
	mc := &mockChain{in:"input", out:"output", fn: func(s string) string { return "m:"+s }}
	dc := &mockChain{in:"input", out:"output", fn: func(s string) string { return "d:"+s }}
	r := NewRouterChain().AddRouteWithKeywords("m", "m", mc, []string{"math"}).WithDefault(dc)
	result, _ := r.Invoke(map[string]any{"input": "hello"})
	if result["output"] != "d:hello" { t.Errorf("got %v", result["output"]) }
}
