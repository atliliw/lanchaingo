package handlers

import "testing"

// 娴嬭瘯 StdOutHandler 鍒涘缓
func TestNewStdOutHandler(t *testing.T) {
	h := NewStdOutHandler()
	if h == nil {
		t.Fatal("expected non-nil handler")
	}
}
