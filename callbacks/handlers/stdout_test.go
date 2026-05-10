package handlers

import "testing"

func TestNewStdOutHandler(t *testing.T) {
	h := NewStdOutHandler()
	if h == nil { t.Fatal("expected non-nil") }
}
