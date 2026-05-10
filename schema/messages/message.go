package messages

import (
	"github.com/atliliw/lanchaingo/core"
)

// MessageType categorizes the role of a message in a conversation.
type MessageType int

const (
	MessageTypeSystem MessageType = iota
	MessageTypeHuman
	MessageTypeAI
	MessageTypeTool
)

func (t MessageType) String() string {
	switch t {
	case MessageTypeSystem:
		return "system"
	case MessageTypeHuman:
		return "human"
	case MessageTypeAI:
		return "ai"
	case MessageTypeTool:
		return "tool"
	default:
		return "unknown"
	}
}

// Message represents a single message in a conversation.
// Maps to Rust langchainrust::schema::messages::message::Message.
type Message struct {
	Content          string              `json:"content"`
	MessageType      MessageType         `json:"message_type"`
	Name             string              `json:"name,omitempty"`
	ToolCallID       string              `json:"tool_call_id,omitempty"`
	ToolCalls        []core.ToolCall     `json:"tool_calls,omitempty"`
	AdditionalKwargs map[string]any      `json:"additional_kwargs,omitempty"`
}

func NewSystemMessage(content string) Message {
	return Message{
		Content:     content,
		MessageType: MessageTypeSystem,
	}
}

func NewHumanMessage(content string) Message {
	return Message{
		Content:     content,
		MessageType: MessageTypeHuman,
	}
}

func NewAIMessage(content string) Message {
	return Message{
		Content:     content,
		MessageType: MessageTypeAI,
	}
}

func NewToolMessage(content string, toolCallID string) Message {
	return Message{
		Content:     content,
		MessageType: MessageTypeTool,
		ToolCallID:  toolCallID,
	}
}

func (m Message) IsSystem() bool {
	return m.MessageType == MessageTypeSystem
}

func (m Message) IsHuman() bool {
	return m.MessageType == MessageTypeHuman
}

func (m Message) IsAI() bool {
	return m.MessageType == MessageTypeAI
}

func (m Message) IsTool() bool {
	return m.MessageType == MessageTypeTool
}
