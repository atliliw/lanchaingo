package memory

import (
	"fmt"
	"strings"

	"github.com/atliliw/lanchaingo/schema/messages"
)

// ConversationBufferWindowMemory keeps only the last k turns of conversation.
type ConversationBufferWindowMemory struct {
	ChatMemory     *ChatMessageHistory
	K              int
	InputKey       string
	OutputKey      string
	MemoryKey      string
	ReturnMessages bool
}

func NewConversationBufferWindowMemory(k int) *ConversationBufferWindowMemory {
	if k <= 0 {
		k = 5
	}
	return &ConversationBufferWindowMemory{
		ChatMemory: NewChatMessageHistory(),
		K:          k,
		InputKey:   "input",
		OutputKey:  "output",
		MemoryKey:  "history",
	}
}

func (m *ConversationBufferWindowMemory) WithInputKey(key string) *ConversationBufferWindowMemory {
	m.InputKey = key
	return m
}

func (m *ConversationBufferWindowMemory) WithOutputKey(key string) *ConversationBufferWindowMemory {
	m.OutputKey = key
	return m
}

func (m *ConversationBufferWindowMemory) WithMemoryKey(key string) *ConversationBufferWindowMemory {
	m.MemoryKey = key
	return m
}

func (m *ConversationBufferWindowMemory) WithReturnMessages(v bool) *ConversationBufferWindowMemory {
	m.ReturnMessages = v
	return m
}

func (m *ConversationBufferWindowMemory) getWindowMessages() []messages.Message {
	all := m.ChatMemory.msgs
	total := len(all)
	maxMsgs := m.K * 2
	if total <= maxMsgs {
		cp := make([]messages.Message, total)
		copy(cp, all)
		return cp
	}
	cp := make([]messages.Message, maxMsgs)
	copy(cp, all[total-maxMsgs:])
	return cp
}

func (m *ConversationBufferWindowMemory) bufferAsString() string {
	var parts []string
	for _, msg := range m.getWindowMessages() {
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

func (m *ConversationBufferWindowMemory) MemoryVariables() []string {
	return []string{m.MemoryKey}
}

func (m *ConversationBufferWindowMemory) LoadMemoryVariables(_ map[string]string) (map[string]any, error) {
	result := make(map[string]any)
	if m.ReturnMessages {
		result[m.MemoryKey] = m.getWindowMessages()
	} else {
		result[m.MemoryKey] = m.bufferAsString()
	}
	return result, nil
}

func (m *ConversationBufferWindowMemory) SaveContext(inputs, outputs map[string]string) error {
	if input, ok := inputs[m.InputKey]; ok {
		m.ChatMemory.AddUserMessage(input)
	}
	if output, ok := outputs[m.OutputKey]; ok {
		m.ChatMemory.AddAIMessage(output)
	}
	return nil
}

func (m *ConversationBufferWindowMemory) Clear() error {
	m.ChatMemory.Clear()
	return nil
}
