package tools

import (
	"context"
	"testing"
)

func TestNewToolRegistry(t *testing.T) {
	r := NewToolRegistry()
	if r.Len() != 0 { t.Errorf("expected empty, got %d", r.Len()) }
}
func TestRegisterAndGet(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name:"calc",description:"calc"})
	got, ok := r.Get("calc")
	if !ok { t.Fatal("expected to find calc") }
	if got.Name() != "calc" { t.Errorf("expected calc") }
}
func TestRegisterDuplicate(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name:"calc",description:"calc"})
	err := r.Register(&mockTool{name:"calc",description:"calc2"})
	if err == nil { t.Error("expected error") }
}
func TestList(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name:"a",description:"a"})
	r.Register(&mockTool{name:"b",description:"b"})
	names := r.List()
	if len(names) != 2 { t.Fatalf("expected 2, got %d", len(names)) }
}
func TestRegisterAll(t *testing.T) {
	r := NewToolRegistry()
	r.RegisterAll(&mockTool{name:"a",description:"a"}, &mockTool{name:"b",description:"b"})
	if r.Len() != 2 { t.Errorf("expected 2") }
}
func TestToToolDefinitions(t *testing.T) {
	r := NewToolRegistry()
	r.Register(&mockTool{name:"calc",description:"calculator"})
	defs := r.ToToolDefinitions()
	if len(defs) != 1 { t.Fatalf("expected 1") }
	if defs[0].Name != "calc" { t.Errorf("expected calc") }
}
func TestNewStructuredTool(t *testing.T) {
	tool := NewStructuredTool("echo","echoes","string",map[string]any{"type":"string"}, nil)
	if tool.Name() != "echo" { t.Errorf("expected echo") }
}
func TestToolError(t *testing.T) {
	err := NewToolError("calc","division by zero",nil)
	if err.Error() != "tool calc: division by zero" { t.Errorf("unexpected: %s", err.Error()) }
}
func TestNewStructuredOutput(t *testing.T) {
	type TestS struct{ Name string }
	so, err := NewStructuredOutput(TestS{Name:"test"})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if so.Schema == nil { t.Error("expected non-nil schema") }
}

type mockTool struct{ name, description string }
func (m *mockTool) Name() string { return m.name }
func (m *mockTool) Description() string { return m.description }
func (m *mockTool) Run(_ context.Context, _ string) (string, error) { return "ok", nil }
