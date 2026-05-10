package memory

// ConversationBufferMemory stores all conversation history in a buffer.
type ConversationBufferMemory struct {
	ChatMemory     *ChatMessageHistory
	InputKey       string
	OutputKey      string
	MemoryKey      string
	ReturnMessages bool
}

func NewConversationBufferMemory() *ConversationBufferMemory {
	return &ConversationBufferMemory{
		ChatMemory: NewChatMessageHistory(),
		InputKey:   "input",
		OutputKey:  "output",
		MemoryKey:  "history",
	}
}

func (m *ConversationBufferMemory) WithInputKey(key string) *ConversationBufferMemory {
	m.InputKey = key
	return m
}

func (m *ConversationBufferMemory) WithOutputKey(key string) *ConversationBufferMemory {
	m.OutputKey = key
	return m
}

func (m *ConversationBufferMemory) WithMemoryKey(key string) *ConversationBufferMemory {
	m.MemoryKey = key
	return m
}

func (m *ConversationBufferMemory) WithReturnMessages(v bool) *ConversationBufferMemory {
	m.ReturnMessages = v
	return m
}

func (m *ConversationBufferMemory) MemoryVariables() []string {
	return []string{m.MemoryKey}
}

func (m *ConversationBufferMemory) LoadMemoryVariables(_ map[string]string) (map[string]any, error) {
	result := make(map[string]any)
	if m.ReturnMessages {
		result[m.MemoryKey] = m.ChatMemory.Messages()
	} else {
		result[m.MemoryKey] = m.ChatMemory.String()
	}
	return result, nil
}

func (m *ConversationBufferMemory) SaveContext(inputs, outputs map[string]string) error {
	if input, ok := inputs[m.InputKey]; ok {
		m.ChatMemory.AddUserMessage(input)
	}
	if output, ok := outputs[m.OutputKey]; ok {
		m.ChatMemory.AddAIMessage(output)
	}
	return nil
}

func (m *ConversationBufferMemory) Clear() error {
	m.ChatMemory.Clear()
	return nil
}
