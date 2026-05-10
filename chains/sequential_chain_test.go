package chains

import (
	"testing"
)

// mockChain 鏄竴涓畝鍗曠殑 BaseChain 瀹炵幇锛岀敤浜庢祴璇?SequentialChain
type mockChain struct {
	inputKey  string
	outputKey string
	transform func(string) string
}

func (c *mockChain) InputKeys() []string           { return []string{c.inputKey} }
func (c *mockChain) OutputKeys() []string          { return []string{c.outputKey} }
func (c *mockChain) Name() string                  { return "mock" }
func (c *mockChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.inputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing: "+c.inputKey, nil)
	}
	return nil
}
func (c *mockChain) Invoke(inputs map[string]any) (ChainResult, error) {
	in, _ := inputs[c.inputKey].(string)
	return ChainResult{c.outputKey: c.transform(in)}, nil
}

// 娴嬭瘯 SequentialChain 椤哄簭鎵ц
func TestSequentialChain(t *testing.T) {
	upper := &mockChain{inputKey: "text", outputKey: "upper", transform: func(s string) string {
		return s + "_upper"
	}}
	reverse := &mockChain{inputKey: "upper", outputKey: "result", transform: func(s string) string {
		return s + "_reversed"
	}}

	chain := NewSequentialChain().
		AddChain(upper, []string{"text"}, []string{"upper"}).
		AddChain(reverse, []string{"upper"}, []string{"result"})

	result, err := chain.Invoke(map[string]any{"text": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result["result"] != "hello_upper_reversed" {
		t.Errorf("expected hello_upper_reversed, got %v", result["result"])
	}
	if result["upper"] != "hello_upper" {
		t.Errorf("expected hello_upper, got %v", result["upper"])
	}
}

// 娴嬭瘯缂哄皯杈撳叆閿欒
func TestSequentialChainMissingInput(t *testing.T) {
	chain := NewSequentialChain().
		AddChain(&mockChain{inputKey: "x", outputKey: "y", transform: func(s string) string { return s }}, []string{"x"}, []string{"y"})

	_, err := chain.Invoke(map[string]any{})
	if err == nil {
		t.Error("expected error for missing input")
	}
}

// 娴嬭瘯 InputKeys/OutputKeys
func TestSequentialChainKeys(t *testing.T) {
	chain := NewSequentialChain().
		AddChain(&mockChain{inputKey: "a", outputKey: "b", transform: func(s string) string { return s }}, []string{"a"}, []string{"b"})

	inputs := chain.InputKeys()
	if len(inputs) != 1 || inputs[0] != "a" {
		t.Errorf("unexpected input keys: %v", inputs)
	}
}
