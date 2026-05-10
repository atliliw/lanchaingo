package language_models

import (
	"context"
	"testing"
	"time"
)

func TestNewChatStream(t *testing.T) {
	cs := NewChatStream(0)
	if cs == nil { t.Fatal("expected non-nil ChatStream") }
	if cs.Data == nil { t.Error("expected non-nil Data channel") }
}
func TestChatStreamSend(t *testing.T) {
	cs := NewChatStream(10)
	go func() { cs.Data <- "hello"; cs.Close() }()
	msg := <-cs.Data
	if msg != "hello" { t.Errorf("expected hello, got %s", msg) }
}
func TestChatStreamClose(t *testing.T) {
	cs := NewChatStream(10)
	cs.Close()
	cs.Close()
}
func TestModelConfigValidate(t *testing.T) {
	cfg1 := ModelConfig{APIKey: "sk-xxx", Model: "gpt-4"}
	if err := cfg1.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	cfg2 := ModelConfig{Model: "gpt-4"}
	if err := cfg2.Validate(); err == nil {
		t.Error("expected error")
	}
}
func TestChatStreamTimeout(t *testing.T) {
	cs := NewChatStream(10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	select {
	case <-cs.Data: t.Error("unexpected data")
	case <-ctx.Done():
	}
}
