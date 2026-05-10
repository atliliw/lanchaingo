package prompts

import (
	"testing"

	"github.com/atliliw/lanchaingo/schema/messages"
)

// 娴嬭瘯 NewChatPromptTemplate锛氭柊寤虹殑妯℃澘涓嶄负 nil
func TestNewChatPromptTemplate(t *testing.T) {
	cpt := NewChatPromptTemplate()
	if cpt == nil {
		t.Fatal("expected non-nil ChatPromptTemplate")
	}
}

// 娴嬭瘯 AddSystemMessage + AddHumanMessage + FormatMessages锛?// 娑堟伅绫诲瀷姝ｇ‘锛屾ā鏉夸腑鐨?{variable} 琚浛鎹?func TestChatPromptTemplateAddAndFormat(t *testing.T) {
	cpt := NewChatPromptTemplate()
	cpt.AddSystemMessage("You are a {role}")
	cpt.AddHumanMessage("Hello, my name is {name}")

	msgs, err := cpt.FormatMessages(map[string]string{
		"role": "assistant",
		"name": "Alice",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}

	if msgs[0].MessageType != messages.MessageTypeSystem {
		t.Errorf("expected system message, got %v", msgs[0].MessageType)
	}
	if msgs[0].Content != "You are a assistant" {
		t.Errorf("expected 'You are a assistant', got '%s'", msgs[0].Content)
	}

	if msgs[1].MessageType != messages.MessageTypeHuman {
		t.Errorf("expected human message, got %v", msgs[1].MessageType)
	}
	if msgs[1].Content != "Hello, my name is Alice" {
		t.Errorf("expected 'Hello, my name is Alice', got '%s'", msgs[1].Content)
	}
}

// 娴嬭瘯涓夌渚垮埄鏂规硶锛欰ddSystemMessage / AddHumanMessage / AddAIMessage
// 楠岃瘉绫诲瀷鍜屽唴瀹归兘姝ｇ‘
func TestChatPromptTemplateConvenienceMethods(t *testing.T) {
	cpt := NewChatPromptTemplate()
	cpt.AddSystemMessage("system msg")
	cpt.AddHumanMessage("human msg")
	cpt.AddAIMessage("ai msg")

	msgs, err := cpt.FormatMessages(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(msgs))
	}

	if msgs[0].MessageType != messages.MessageTypeSystem || msgs[0].Content != "system msg" {
		t.Error("first message mismatch")
	}
	if msgs[1].MessageType != messages.MessageTypeHuman || msgs[1].Content != "human msg" {
		t.Error("second message mismatch")
	}
	if msgs[2].MessageType != messages.MessageTypeAI || msgs[2].Content != "ai msg" {
		t.Error("third message mismatch")
	}
}

// 娴嬭瘯妯℃澘鍙橀噺鏇挎崲锛欻umanMessage 妯℃澘涓殑 {name} 鍜?{age} 琚纭浛鎹?func TestChatPromptTemplateWithVariables(t *testing.T) {
	cpt := NewChatPromptTemplate()
	cpt.AddHumanMessage("Hi, I'm {name} and I'm {age} years old")

	msgs, err := cpt.FormatMessages(map[string]string{
		"name": "Bob",
		"age":  "25",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "Hi, I'm Bob and I'm 25 years old"
	if msgs[0].Content != expected {
		t.Errorf("expected '%s', got '%s'", expected, msgs[0].Content)
	}
}
