package tools

import (
	"context"
	"testing"
)

// mockTool 鏄竴涓畝鍗曠殑 Tool 鎺ュ彛瀹炵幇锛岀敤浜庢祴璇?ToolRegistry
type mockTool struct {
	name        string
	description string
}

func (m *mockTool) Name() string                       { return m.name }
func (m *mockTool) Description() string                { return m.description }
func (m *mockTool) Run(_ context.Context, _ string) (string, error) { return "ok", nil }

// 娴嬭瘯 NewToolRegistry锛氭柊寤虹殑娉ㄥ唽鍣ㄥ簲涓虹┖锛孡en() 杩斿洖 0
func TestNewToolRegistry(t *testing.T) {
	r := NewToolRegistry()
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
	if r.Len() != 0 {
		t.Errorf("expected empty registry, got %d items", r.Len())
	}
}

// 娴嬭瘯 Register + Get锛氭敞鍐屽伐鍏峰悗鑳介€氳繃鍚嶇О姝ｇ‘鏌ユ壘鍒?func TestRegisterAndGet(t *testing.T) {
	r := NewToolRegistry()
	tool := &mockTool{name: "calculator", description: "does math"}

	err := r.Register(tool)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := r.Get("calculator")
	if !ok {
		t.Fatal("expected to find calculator")
	}
	if got.Name() != "calculator" {
		t.Errorf("expected calculator, got %s", got.Name())
	}
}

// 娴嬭瘯 Register 閲嶅娉ㄥ唽锛氬悓鍚嶅伐鍏峰簲杩斿洖閿欒
func TestRegisterDuplicate(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name: "calc", description: "calc"})

	err := r.Register(&mockTool{name: "calc", description: "calc2"})
	if err == nil {
		t.Error("expected error for duplicate registration")
	}
}

// 娴嬭瘯 MustGet锛氫笉瀛樺湪鐨勫伐鍏峰簲 panic
// 浣跨敤 defer recover 鎹曡幏 panic 鏉ラ獙璇?func TestMustGet(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name: "exists", description: "exists"})

	defer func() {
		if err := recover(); err == nil {
			t.Error("expected panic for missing tool")
		}
	}()
	r.MustGet("nonexistent")
}

// 娴嬭瘯 Remove锛氱Щ闄ゅ悗娉ㄥ唽鍣ㄥ彉绌?func TestRemove(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name: "calc", description: "calc"})
	r.Remove("calc")

	if r.Len() != 0 {
		t.Error("expected empty registry after remove")
	}
}

// 娴嬭瘯 List锛氳繑鍥炴墍鏈夊凡娉ㄥ唽鐨勫伐鍏峰悕绉板垪琛?func TestList(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name: "a", description: "a"})
	r.Register(&mockTool{name: "b", description: "b"})

	names := r.List()
	if len(names) != 2 {
		t.Fatalf("expected 2 names, got %d", len(names))
	}
}

// 娴嬭瘯 RegisterAll锛氭壒閲忔敞鍐屽涓伐鍏凤紝楠岃瘉鏁伴噺鍜屽唴瀹?func TestRegisterAll(t *testing.T) {
	r := NewToolRegistry()
	err := r.RegisterAll(
		&mockTool{name: "a", description: "a"},
		&mockTool{name: "b", description: "b"},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 2 {
		t.Errorf("expected 2 tools, got %d", r.Len())
	}
}

// 娴嬭瘯 ToToolDefinitions锛氶獙璇佹敞鍐岀殑宸ュ叿鑳芥纭浆鎹负 LLM 鍙瘑鍒殑 ToolDefinition 鍒楄〃
func TestToToolDefinitions(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name: "calc", description: "calculator"})

	defs := r.ToToolDefinitions()
	if len(defs) != 1 {
		t.Fatalf("expected 1 definition, got %d", len(defs))
	}
	if defs[0].Name != "calc" {
		t.Errorf("expected calc, got %s", defs[0].Name)
	}
	if defs[0].Description != "calculator" {
		t.Errorf("expected calculator, got %s", defs[0].Description)
	}
}
