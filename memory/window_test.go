package memory

import (
	"testing"
)

// 娴嬭瘯 WindowMemory锛氬彧淇濈暀鏈€杩?k 杞璇?func TestWindowMemoryOnlyKeepsLastK(t *testing.T) {
	m := NewConversationBufferWindowMemory(2)

	// 娣诲姞 3 杞璇濓紙6 鏉℃秷鎭級
	for i := 1; i <= 3; i++ {
		m.SaveContext(
			map[string]string{"input": "q" + string(rune('0'+i))},
			map[string]string{"output": "a" + string(rune('0'+i))},
		)
	}

	// 瀹屾暣鍘嗗彶搴旀湁 6 鏉℃秷鎭?	if m.ChatMemory.Len() != 6 {
		t.Fatalf("expected 6 messages total, got %d", m.ChatMemory.Len())
	}

	// 浣嗗姞杞芥椂搴斿彧杩斿洖鏈€杩?2 杞紙4 鏉★級
	vars, _ := m.LoadMemoryVariables(nil)
	history := vars["history"].(string)

	if m.ChatMemory.Len() != 6 {
		t.Errorf("full history should still be 6, got %d", m.ChatMemory.Len())
	}

	expected := "Human: q2\nAI: a2\nHuman: q3\nAI: a3"
	if history != expected {
		t.Errorf("expected %q, got %q", expected, history)
	}
}

// 娴嬭瘯 WindowMemory锛氭秷鎭暟灏忎簬 k 鏃惰繑鍥炲叏閮?func TestWindowMemorySmallerThanK(t *testing.T) {
	m := NewConversationBufferWindowMemory(5)

	m.SaveContext(map[string]string{"input": "q1"}, map[string]string{"output": "a1"})
	m.SaveContext(map[string]string{"input": "q2"}, map[string]string{"output": "a2"})

	vars, _ := m.LoadMemoryVariables(nil)
	history := vars["history"].(string)

	if history != "Human: q1\nAI: a1\nHuman: q2\nAI: a2" {
		t.Errorf("unexpected: %s", history)
	}
}

// 娴嬭瘯 Clear
func TestWindowMemoryClear(t *testing.T) {
	m := NewConversationBufferWindowMemory(2)
	m.SaveContext(map[string]string{"input": "t"}, map[string]string{"output": "r"})
	m.Clear()

	vars, _ := m.LoadMemoryVariables(nil)
	history := vars["history"].(string)
	if history != "" {
		t.Errorf("expected empty, got %s", history)
	}
}

// 娴嬭瘯鑷畾涔夐敭鍚嶅拰杩斿洖娑堟伅瀵硅薄
func TestWindowMemoryReturnMessages(t *testing.T) {
	m := NewConversationBufferWindowMemory(3).WithReturnMessages(true)
	m.SaveContext(map[string]string{"input": "hi"}, map[string]string{"output": "hello"})

	vars, _ := m.LoadMemoryVariables(nil)
	if _, ok := vars["history"].([]Message); !ok {
		t.Errorf("expected []Message")
	}
}
