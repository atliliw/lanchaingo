package callbacks

import "testing"

// 娴嬭瘯 CallbackManager 鍒涘缓
func TestNewCallbackManager(t *testing.T) {
	m := NewCallbackManager()
	if m == nil {
		t.Fatal("expected non-nil manager")
	}
	if !m.IsEmpty() {
		t.Error("expected empty manager")
	}
}
