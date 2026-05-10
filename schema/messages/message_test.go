package messages

import (
	"testing"

	"github.com/atliliw/lanchaingo/core"
)

func TestNewSystemMessage(t *testing.T) {
	msg := NewSystemMessage("You are a helpful assistant")
	if msg.Content != "You are a helpful assistant" {
		t.Errorf("unexpected content: %s", msg.Content)
	}
	if msg.MessageType != MessageTypeSystem {
		t.Errorf("expected MessageTypeSystem, got %v", msg.MessageType)
	}
}

func TestNewHumanMessage(t *testing.T) {
	msg := NewHumanMessage("Hello")
	if msg.MessageType != MessageTypeHuman {
		t.Errorf("expected MessageTypeHuman, got %v", msg.MessageType)
	}
}

func TestNewAIMessage(t *testing.T) {
	msg := NewAIMessage("I'm AI")
	if msg.MessageType != MessageTypeAI {
		t.Errorf("expected MessageTypeAI, got %v", msg.MessageType)
	}
}

func TestNewToolMessage(t *testing.T) {
	msg := NewToolMessage("result", "call_123")
	if msg.MessageType != MessageTypeTool {
		t.Errorf("expected MessageTypeTool, got %v", msg.MessageType)
	}
	if msg.ToolCallID != "call_123" {
		t.Errorf("expected call_123, got %s", msg.ToolCallID)
	}
}

func TestMessageTypeCheckers(t *testing.T) {
	tests := []struct {
		name     string
		msg      Message
		isSystem bool
		isHuman  bool
		isAI     bool
		isTool   bool
	}{
		{"system", NewSystemMessage("s"), true, false, false, false},
		{"human", NewHumanMessage("h"), false, true, false, false},
		{"ai", NewAIMessage("a"), false, false, true, false},
		{"tool", NewToolMessage("t", "c"), false, false, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.msg.IsSystem() != tt.isSystem {
				t.Errorf("IsSystem = %v, want %v", tt.msg.IsSystem(), tt.isSystem)
			}
			if tt.msg.IsHuman() != tt.isHuman {
				t.Errorf("IsHuman = %v, want %v", tt.msg.IsHuman(), tt.isHuman)
			}
			if tt.msg.IsAI() != tt.isAI {
				t.Errorf("IsAI = %v, want %v", tt.msg.IsAI(), tt.isAI)
			}
			if tt.msg.IsTool() != tt.isTool {
				t.Errorf("IsTool = %v, want %v", tt.msg.IsTool(), tt.isTool)
			}
		})
	}
}

func TestMessageWithToolCalls(t *testing.T) {
	msg := Message{
		Content:     "I'll use calculator",
		MessageType: MessageTypeAI,
		ToolCalls: []core.ToolCall{
			{ID: "call_1", Name: "calculator", Arguments: `{"x":1}`},
		},
	}
	if len(msg.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(msg.ToolCalls))
	}
	if msg.ToolCalls[0].Name != "calculator" {
		t.Errorf("expected calculator, got %s", msg.ToolCalls[0].Name)
	}
}

func TestMessageStringer(t *testing.T) {
	tests := []struct {
		mt   MessageType
		want string
	}{
		{MessageTypeSystem, "system"},
		{MessageTypeHuman, "human"},
		{MessageTypeAI, "ai"},
		{MessageTypeTool, "tool"},
		{MessageType(99), "unknown"},
	}
	for _, tt := range tests {
		if got := tt.mt.String(); got != tt.want {
			t.Errorf("MessageType(%d).String() = %s, want %s", tt.mt, got, tt.want)
		}
	}
}
