package chains

import (
	"context"
	"testing"

	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

type docChainLLM struct{}

func (d *docChainLLM) ModelName() string { return "mock" }
func (d *docChainLLM) GetNumTokens(string) int { return 0 }
func (d *docChainLLM) Temperature() *float32 { return nil }
func (d *docChainLLM) MaxTokens() *int { return nil }
func (d *docChainLLM) WithTemperature(float32) lm.BaseLanguageModel { return d }
func (d *docChainLLM) WithMaxTokens(int) lm.BaseLanguageModel { return d }
func (d *docChainLLM) Chat(_ context.Context, _ []messages.Message) (*lm.LLMResult, error) {
	return &lm.LLMResult{Content: "summary"}, nil
}
func (d *docChainLLM) StreamChat(_ context.Context, _ []messages.Message) (<-chan string, error) {
	ch := make(chan string); close(ch); return ch, nil
}
func (d *docChainLLM) BindTools(_ []lm.ToolDefinition) {}

// 娴嬭瘯 StuffDocumentsChain
func TestStuffDocumentsChain(t *testing.T) {
	c := NewStuffDocumentsChain(&docChainLLM{}, "summarize:")
	result, err := c.Invoke(map[string]any{"input_documents": []string{"doc1", "doc2"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["output_text"] != "summary" {
		t.Errorf("expected summary, got %v", result["output_text"])
	}
}

// 娴嬭瘯 RefineDocumentsChain
func TestRefineDocumentsChain(t *testing.T) {
	c := NewRefineDocumentsChain(&docChainLLM{}, "initial:", "refine:")
	result, err := c.Invoke(map[string]any{"input_documents": []string{"doc1", "doc2"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["output_text"] != "summary" {
		t.Errorf("expected summary, got %v", result["output_text"])
	}
}

// 娴嬭瘯 MapReduceDocumentsChain
func TestMapReduceDocumentsChain(t *testing.T) {
	c := NewMapReduceDocumentsChain(&docChainLLM{}, "map:", "reduce:")
	result, err := c.Invoke(map[string]any{"input_documents": []string{"a", "b"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["output_text"] != "summary" {
		t.Errorf("expected summary, got %v", result["output_text"])
	}
}
