package callbacks

import "testing"

func TestNewCallbackManager(t *testing.T) {
	m := NewCallbackManager()
	if m == nil { t.Fatal("expected non-nil") }
	if !m.IsEmpty() { t.Error("expected empty") }
}
