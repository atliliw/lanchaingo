package memory

import (
	"testing"
)

// 娴嬭瘯 ConversationBufferMemory 鍩虹锛氫繚瀛樹笂涓嬫枃鍚庤兘姝ｇ‘鍔犺浇鍘嗗彶
func TestBufferMemorySaveAndLoad(t *testing.T) {
	m := NewConversationBufferMemory()

	inputs := map[string]string{"input": "浣犲ソ"}
	outputs := map[string]string{"output": "浣犲ソ锛?}
	m.SaveContext(inputs, outputs)

	vars, err := m.LoadMemoryVariables(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	history, ok := vars["history"].(string)
	if !ok {
		t.Fatalf("expected string, got %T", vars["history"])
	}
	if history != "Human: 浣犲ソ\nAI: 浣犲ソ锛? {
		t.Errorf("unexpected history: %s", history)
	}
}

// 娴嬭瘯澶氳疆瀵硅瘽璁板繂绱Н
func TestBufferMemoryMultipleTurns(t *testing.T) {
	m := NewConversationBufferMemory()

	m.SaveContext(map[string]string{"input": "q1"}, map[string]string{"output": "a1"})
	m.SaveContext(map[string]string{"input": "q2"}, map[string]string{"output": "a2"})

	vars, _ := m.LoadMemoryVariables(nil)
	history := vars["history"].(string)

	if history != "Human: q1\nAI: a1\nHuman: q2\nAI: a2" {
		t.Errorf("unexpected history: %s", history)
	}
}

// 娴嬭瘯 Clear 娓呯┖璁板繂
func TestBufferMemoryClear(t *testing.T) {
	m := NewConversationBufferMemory()
	m.SaveContext(map[string]string{"input": "test"}, map[string]string{"output": "ok"})
	m.Clear()

	vars, _ := m.LoadMemoryVariables(nil)
	history := vars["history"].(string)
	if history != "" {
		t.Errorf("expected empty after clear, got: %s", history)
	}
}

// 娴嬭瘯 ReturnMessages=true 杩斿洖娑堟伅瀵硅薄鑰岄潪瀛楃涓?func TestBufferMemoryReturnMessages(t *testing.T) {
	m := NewConversationBufferMemory().WithReturnMessages(true)
	m.SaveContext(map[string]string{"input": "hi"}, map[string]string{"output": "hello"})

	vars, _ := m.LoadMemoryVariables(nil)
	_, ok := vars["history"].([]Message)
	if !ok {
		t.Errorf("expected []Message, got %T", vars["history"])
	}
}

// 娴嬭瘯鑷畾涔夐敭鍚?func TestBufferMemoryCustomKeys(t *testing.T) {
	m := NewConversationBufferMemory().
		WithInputKey("question").
		WithOutputKey("answer").
		WithMemoryKey("context")

	m.SaveContext(map[string]string{"question": "what"}, map[string]string{"answer": "42"})
	vars, _ := m.LoadMemoryVariables(nil)

	history, ok := vars["context"]
	if !ok {
		t.Fatal("expected context key")
	}
	if history.(string) != "Human: what\nAI: 42" {
		t.Errorf("unexpected: %s", history)
	}
}

// 娴嬭瘯 MemoryVariables 杩斿洖鍊?func TestBufferMemoryVariables(t *testing.T) {
	m := NewConversationBufferMemory().WithMemoryKey("my_history")
	vars := m.MemoryVariables()
	if len(vars) != 1 || vars[0] != "my_history" {
		t.Errorf("unexpected vars: %v", vars)
	}
}
