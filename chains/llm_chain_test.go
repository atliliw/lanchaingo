package chains

import (
	"context"
	"testing"

	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// mockLLM 瀹炵幇 BaseChatModel 鎺ュ彛鐢ㄤ簬娴嬭瘯
type mockLLM struct {
	respond func(msgs []messages.Message) (*lm.LLMResult, error)
}

func (m *mockLLM) ModelName() string                        { return "mock" }
func (m *mockLLM) GetNumTokens(text string) int             { return len(text) / 4 }
func (m *mockLLM) Temperature() *float32                    { return nil }
func (m *mockLLM) MaxTokens() *int                          { return nil }
func (m *mockLLM) WithTemperature(f float32) lm.BaseLanguageModel { return m }
func (m *mockLLM) WithMaxTokens(i int) lm.BaseLanguageModel       { return m }
func (m *mockLLM) Chat(_ context.Context, msgs []messages.Message) (*lm.LLMResult, error) {
	return m.respond(msgs)
}
func (m *mockLLM) StreamChat(_ context.Context, _ []messages.Message) (<-chan string, error) {
	ch := make(chan string)
	close(ch)
	return ch, nil
}
func (m *mockLLM) BindTools(_ []lm.ToolDefinition) {}

// 娴嬭瘯 LLMChain.Invoke锛氳緭鍏ユ浛鎹㈠悗璋冪敤 LLM 杩斿洖缁撴灉
func TestLLMChainInvoke(t *testing.T) {
	llm := &mockLLM{
		respond: func(msgs []messages.Message) (*lm.LLMResult, error) {
			return &lm.LLMResult{Content: "Hello, World!"}, nil
		},
	}
	chain := NewLLMChain(llm, "Say {input}")

	result, err := chain.Invoke(map[string]any{"input": "hi"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["text"] != "Hello, World!" {
		t.Errorf("expected Hello World!, got %v", result["text"])
	}
}

// 娴嬭瘯缂哄皯杈撳叆杩斿洖閿欒
func TestLLMChainMissingInput(t *testing.T) {
	llm := &mockLLM{respond: func(_ []messages.Message) (*lm.LLMResult, error) {
		return &lm.LLMResult{Content: "ok"}, nil
	}}
	chain := NewLLMChain(llm, "{input}")

	_, err := chain.Invoke(map[string]any{})
	if err == nil {
		t.Error("expected error for missing input")
	}
}

// 娴嬭瘯鑷畾涔夎緭鍏?杈撳嚭閿悕
func TestLLMChainCustomKeys(t *testing.T) {
	llm := &mockLLM{respond: func(_ []messages.Message) (*lm.LLMResult, error) {
		return &lm.LLMResult{Content: "42"}, nil
	}}
	chain := NewLLMChain(llm, "answer").
		WithInputKey("question").
		WithOutputKey("answer")

	result, err := chain.Invoke(map[string]any{"question": "life"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["answer"] != "42" {
		t.Errorf("expected 42, got %v", result["answer"])
	}
}

// 娴嬭瘯 LLMChainBuilder
func TestLLMChainBuilder(t *testing.T) {
	llm := &mockLLM{respond: func(_ []messages.Message) (*lm.LLMResult, error) {
		return &lm.LLMResult{Content: "result"}, nil
	}}
	chain := NewLLMChainBuilder(llm, "template").
		InputKey("q").
		OutputKey("a").
		Name("test").
		Build()

	if chain.InputKey != "q" || chain.OutputKey != "a" || chain.Name() != "test" {
		t.Errorf("builder fields mismatch: %+v", chain)
	}
}

// 娴嬭瘯 InputKeys/OutputKeys
func TestLLMChainKeys(t *testing.T) {
	llm := &mockLLM{respond: func(_ []messages.Message) (*lm.LLMResult, error) {
		return &lm.LLMResult{Content: "x"}, nil
	}}
	chain := NewLLMChain(llm, "template").WithInputKey("in").WithOutputKey("out")

	inputs := chain.InputKeys()
	if len(inputs) != 1 || inputs[0] != "in" {
		t.Errorf("unexpected input keys: %v", inputs)
	}
	outputs := chain.OutputKeys()
	if len(outputs) != 1 || outputs[0] != "out" {
		t.Errorf("unexpected output keys: %v", outputs)
	}
}
