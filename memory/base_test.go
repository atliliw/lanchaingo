package memory

import (
	"testing"

	"github.com/atliliw/lanchaingo/schema/messages"
)

// 娴嬭瘯 ChatMessageHistory 鍒涘缓鍜屽熀纭€鎿嶄綔
func TestChatMessageHistoryNew(t *testing.T) {
	h := NewChatMessageHistory()
	if !h.IsEmpty() {
		t.Error("expected empty history")
	}
	if h.Len() != 0 {
		t.Errorf("expected len 0, got %d", h.Len())
	}
}

// 娴嬭瘯娣诲姞娑堟伅鍜岃幏鍙?func TestChatMessageHistoryAddAndGet(t *testing.T) {
	h := NewChatMessageHistory()
	h.AddUserMessage("hello")
	h.AddAIMessage("hi")
	h.AddSystemMessage("be nice")

	if h.Len() != 3 {
		t.Fatalf("expected 3 messages, got %d", h.Len())
	}

	msgs := h.Messages()
	if msgs[0].Content != "hello" || msgs[1].Content != "hi" || msgs[2].Content != "be nice" {
		t.Error("message content mismatch")
	}
	if msgs[0].MessageType != messages.MessageTypeHuman {
		t.Error("expected human message")
	}
	if msgs[1].MessageType != messages.MessageTypeAI {
		t.Error("expected ai message")
	}
}

// 娴嬭瘯 FromMessages 宸ュ巶鏂规硶
func TestChatMessageHistoryFromMessages(t *testing.T) {
	msgs := []messages.Message{
		messages.NewHumanMessage("q1"),
		messages.NewAIMessage("a1"),
	}
	h := NewChatMessageHistoryFromMessages(msgs)
	if h.Len() != 2 {
		t.Errorf("expected 2 messages, got %d", h.Len())
	}
	// 淇敼鍘熷鍒囩墖涓嶅簲褰卞搷 history
	msgs[0] = messages.NewHumanMessage("changed")
	if h.Messages()[0].Content == "changed" {
		t.Error("FromMessages should copy the slice")
	}
}

// 娴嬭瘯 Clear
func TestChatMessageHistoryClear(t *testing.T) {
	h := NewChatMessageHistory()
	h.AddUserMessage("test")
	h.Clear()
	if !h.IsEmpty() {
		t.Error("expected empty after clear")
	}
}

// 娴嬭瘯 String 杈撳嚭鏍煎紡
func TestChatMessageHistoryString(t *testing.T) {
	h := NewChatMessageHistory()
	h.AddUserMessage("hello")
	h.AddAIMessage("world")

	str := h.String()
	if str != "Human: hello\nAI: world" {
		t.Errorf("unexpected string: %s", str)
	}
}

// 娴嬭瘯 Messages 杩斿洖鐨勬槸鎷疯礉
func TestChatMessageHistoryMessagesCopy(t *testing.T) {
	h := NewChatMessageHistory()
	h.AddUserMessage("test")
	msgs := h.Messages()
	msgs[0].Content = "modified"
	if h.Messages()[0].Content == "modified" {
		t.Error("Messages() should return a copy")
	}
}
