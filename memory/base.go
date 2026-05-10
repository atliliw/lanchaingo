package memory

import (
	"fmt"
	"strings"

	"github.com/atliliw/lanchaingo/schema/messages"
)

// BaseMemory defines the interface for conversation memory.
type BaseMemory interface {
	MemoryVariables() []string
	LoadMemoryVariables(inputs map[string]string) (map[string]any, error)
	SaveContext(inputs, outputs map[string]string) error
	Clear() error
}

// ChatMessageHistory stores a list of chat messages in memory.
type ChatMessageHistory struct {
	msgs []messages.Message
}

func NewChatMessageHistory() *ChatMessageHistory {
	return &ChatMessageHistory{
		msgs: make([]messages.Message, 0),
	}
}

func NewChatMessageHistoryFromMessages(msgs []messages.Message) *ChatMessageHistory {
	cp := make([]messages.Message, len(msgs))
	copy(cp, msgs)
	return &ChatMessageHistory{msgs: cp}
}

func (h *ChatMessageHistory) AddMessage(msg messages.Message) {
	h.msgs = append(h.msgs, msg)
}

func (h *ChatMessageHistory) AddUserMessage(content string) {
	h.AddMessage(messages.NewHumanMessage(content))
}

func (h *ChatMessageHistory) AddAIMessage(content string) {
	h.AddMessage(messages.NewAIMessage(content))
}

func (h *ChatMessageHistory) AddSystemMessage(content string) {
	h.AddMessage(messages.NewSystemMessage(content))
}

func (h *ChatMessageHistory) Messages() []messages.Message {
	cp := make([]messages.Message, len(h.msgs))
	copy(cp, h.msgs)
	return cp
}

func (h *ChatMessageHistory) Clear() {
	h.msgs = make([]messages.Message, 0)
}

func (h *ChatMessageHistory) Len() int {
	return len(h.msgs)
}

func (h *ChatMessageHistory) IsEmpty() bool {
	return len(h.msgs) == 0
}

func (h *ChatMessageHistory) String() string {
	var parts []string
	for _, msg := range h.msgs {
		role := "Unknown"
		switch msg.MessageType {
		case messages.MessageTypeSystem:
			role = "System"
		case messages.MessageTypeHuman:
			role = "Human"
		case messages.MessageTypeAI:
			role = "AI"
		case messages.MessageTypeTool:
			role = "Tool"
		}
		parts = append(parts, fmt.Sprintf("%s: %s", role, msg.Content))
	}
	return strings.Join(parts, "\n")
}
