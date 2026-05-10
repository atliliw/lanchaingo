package messages

import (
	"testing"

	"github.com/atliliw/lanchaingo/core"
)

// 娴嬭瘯 NewSystemMessage锛欳ontent 鍜?MessageType 姝ｇ‘
func TestNewSystemMessage(t *testing.T) {
	msg := NewSystemMessage("You are a helpful assistant")
	if msg.Content != "You are a helpful assistant" {
		t.Errorf("unexpected content: %s", msg.Content)
	}
	if msg.MessageType != MessageTypeSystem {
		t.Errorf("expected MessageTypeSystem, got %v", msg.MessageType)
	}
}

// 娴嬭瘯 NewHumanMessage锛歁essageType 涓?Human
func TestNewHumanMessage(t *testing.T) {
	msg := NewHumanMessage("Hello")
	if msg.MessageType != MessageTypeHuman {
		t.Errorf("expected MessageTypeHuman, got %v", msg.MessageType)
	}
}

// 娴嬭瘯 NewAIMessage锛歁essageType 涓?AI
func TestNewAIMessage(t *testing.T) {
	msg := NewAIMessage("I'm AI")
	if msg.MessageType != MessageTypeAI {
		t.Errorf("expected MessageTypeAI, got %v", msg.MessageType)
	}
}

// 娴嬭瘯 NewToolMessage锛歁essageType 涓?Tool锛孴oolCallID 姝ｇ‘璧嬪€?func TestNewToolMessage(t *testing.T) {
	msg := NewToolMessage("result", "call_123")
	if msg.MessageType != MessageTypeTool {
		t.Errorf("expected MessageTypeTool, got %v", msg.MessageType)
	}
	if msg.ToolCallID != "call_123" {
		t.Errorf("expected call_123, got %s", msg.ToolCallID)
	}
}

// 娴嬭瘯鍥涚 Message 鐨?IsSystem/IsHuman/IsAI/IsTool 绫诲瀷妫€鏌ユ柟娉?// 浣跨敤琛ㄩ┍鍔ㄦ祴璇曡鐩栨墍鏈夎鑹?func TestMessageTypeCheckers(t *testing.T) {
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

// 娴嬭瘯 Message 鎼哄甫 ToolCalls锛氶獙璇?AI 娑堟伅鍙互闄勫甫澶氫釜宸ュ叿璋冪敤
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

// 娴嬭瘯 MessageType.String()锛氭墍鏈夌被鍨嬪拰鏈煡绫诲瀷閮借兘姝ｇ‘杞崲涓哄瓧绗︿覆
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
