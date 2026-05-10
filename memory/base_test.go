package memory

import "testing"

func TestChatMessageHistoryNew(t *testing.T) {
	h := NewChatMessageHistory()
	if !h.IsEmpty() { t.Error("expected empty") }
}
func TestChatMessageHistoryAdd(t *testing.T) {
	h := NewChatMessageHistory()
	h.AddUserMessage("hello")
	h.AddAIMessage("hi")
	if h.Len() != 2 { t.Errorf("expected 2, got %d", h.Len()) }
}
func TestChatMessageHistoryClear(t *testing.T) {
	h := NewChatMessageHistory()
	h.AddUserMessage("test")
	h.Clear()
	if !h.IsEmpty() { t.Error("expected empty") }
}
func TestChatMessageHistoryString(t *testing.T) {
	h := NewChatMessageHistory()
	h.AddUserMessage("hello")
	h.AddAIMessage("world")
	str := h.String()
	if str != "Human: hello\nAI: world" { t.Errorf("unexpected: %s", str) }
}
func TestBufferMemorySaveAndLoad(t *testing.T) {
	m := NewConversationBufferMemory()
	m.SaveContext(map[string]string{"input":"hello"}, map[string]string{"output":"world"})
	v, _ := m.LoadMemoryVariables(nil)
	if v["history"].(string) != "Human: hello\nAI: world" { t.Error("history mismatch") }
}
func TestBufferMemoryClear(t *testing.T) {
	m := NewConversationBufferMemory()
	m.SaveContext(map[string]string{"input":"t"}, map[string]string{"output":"r"})
	m.Clear()
	v, _ := m.LoadMemoryVariables(nil)
	if v["history"].(string) != "" { t.Error("expected empty") }
}
func TestWindowMemory(t *testing.T) {
	m := NewConversationBufferWindowMemory(2)
	m.SaveContext(map[string]string{"input":"q1"}, map[string]string{"output":"a1"})
	if m.ChatMemory.Len() != 2 { t.Errorf("expected 2 messages") }
}
func TestPersistenceConfig(t *testing.T) {
	cfg := DefaultPersistenceConfig()
	if !cfg.AutoSave { t.Error("expected auto_save true") }
}
