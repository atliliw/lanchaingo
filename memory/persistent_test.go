package memory

import (
	"testing"
)

// 娴嬭瘯 DefaultPersistenceConfig 榛樿鍊?func TestPersistenceConfigDefault(t *testing.T) {
	cfg := DefaultPersistenceConfig()
	if !cfg.AutoSave {
		t.Error("expected auto_save true")
	}
	if !cfg.AutoLoad {
		t.Error("expected auto_load true")
	}
	if cfg.MaxMessages != 100 {
		t.Errorf("expected max_messages 100, got %d", cfg.MaxMessages)
	}
	if cfg.TokenLimit != 4000 {
		t.Errorf("expected token_limit 4000, got %d", cfg.TokenLimit)
	}
}

// 娴嬭瘯 PersistenceConfig 閾惧紡璁剧疆
func TestPersistenceConfigCustom(t *testing.T) {
	cfg := DefaultPersistenceConfig().
		WithAutoSave(false).
		WithMaxMessages(50).
		WithAutoLoad(false)

	if cfg.AutoSave {
		t.Error("expected auto_save false")
	}
	if cfg.AutoLoad {
		t.Error("expected auto_load false")
	}
	if cfg.MaxMessages != 50 {
		t.Errorf("expected 50, got %d", cfg.MaxMessages)
	}
}

// 娴嬭瘯 NewMemoryData
func TestMemoryDataNew(t *testing.T) {
	data := NewMemoryData("session_123")
	if data.SessionID != "session_123" {
		t.Errorf("expected session_123, got %s", data.SessionID)
	}
	if len(data.Messages) != 0 {
		t.Errorf("expected empty messages, got %d", len(data.Messages))
	}
}

// 娴嬭瘯 MemoryData.AddMessage
func TestMemoryDataAddMessage(t *testing.T) {
	data := NewMemoryData("s1")
	data.AddMessage(NewHumanMessage("hello"))
	if len(data.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(data.Messages))
	}
}
