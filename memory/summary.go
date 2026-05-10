package memory

import (
	"fmt"
	"strings"
	"sync"

	"github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

const defaultSummaryPrompt = `Progressively summarize the lines of conversation provided, adding onto the previous summary returning a new summary.

Current summary:
{summary}

New lines of conversation:
{new_lines}

New summary:`

// ConversationSummaryMemory uses LLM to summarize conversation history.
type ConversationSummaryMemory struct {
	llm            language_models.BaseChatModel
	Buffer         string
	ChatMemory     *ChatMessageHistory
	InputKey       string
	OutputKey      string
	MemoryKey      string
	SummaryPrompt  string
	ReturnMessages bool
	mu             sync.RWMutex
}

func NewConversationSummaryMemory(llm language_models.BaseChatModel) *ConversationSummaryMemory {
	return &ConversationSummaryMemory{
		llm:           llm,
		ChatMemory:    NewChatMessageHistory(),
		InputKey:      "input",
		OutputKey:     "output",
		MemoryKey:     "history",
		SummaryPrompt: defaultSummaryPrompt,
	}
}

func (m *ConversationSummaryMemory) WithInputKey(key string) *ConversationSummaryMemory {
	m.InputKey = key
	return m
}

func (m *ConversationSummaryMemory) WithOutputKey(key string) *ConversationSummaryMemory {
	m.OutputKey = key
	return m
}

func (m *ConversationSummaryMemory) WithMemoryKey(key string) *ConversationSummaryMemory {
	m.MemoryKey = key
	return m
}

func (m *ConversationSummaryMemory) WithReturnMessages(v bool) *ConversationSummaryMemory {
	m.ReturnMessages = v
	return m
}

func (m *ConversationSummaryMemory) formatNewLines(input, output string) string {
	return fmt.Sprintf("Human: %s\nAI: %s", input, output)
}

func (m *ConversationSummaryMemory) predictNewSummary(newLines string) (string, error) {
	m.mu.RLock()
	prompt := strings.ReplaceAll(m.SummaryPrompt, "{summary}", m.Buffer)
	m.mu.RUnlock()
	prompt = strings.ReplaceAll(prompt, "{new_lines}", newLines)

	msg := messages.NewHumanMessage(prompt)
	result, err := m.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return "", NewMemoryError(ErrSave, fmt.Sprintf("LLM summary failed: %v", err), err)
	}
	return result.Content, nil
}

func (m *ConversationSummaryMemory) MemoryVariables() []string {
	return []string{m.MemoryKey}
}

func (m *ConversationSummaryMemory) LoadMemoryVariables(_ map[string]string) (map[string]any, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make(map[string]any)
	if m.ReturnMessages {
		summaryMsg := messages.NewSystemMessage(m.Buffer)
		result[m.MemoryKey] = []messages.Message{summaryMsg}
	} else {
		result[m.MemoryKey] = m.Buffer
	}
	return result, nil
}

func (m *ConversationSummaryMemory) SaveContext(inputs, outputs map[string]string) error {
	input := inputs[m.InputKey]
	output := outputs[m.OutputKey]

	m.ChatMemory.AddUserMessage(input)
	m.ChatMemory.AddAIMessage(output)

	newLines := m.formatNewLines(input, output)
	newSummary, err := m.predictNewSummary(newLines)
	if err != nil {
		return err
	}

	m.mu.Lock()
	m.Buffer = newSummary
	m.mu.Unlock()

	return nil
}

func (m *ConversationSummaryMemory) Clear() error {
	m.mu.Lock()
	m.Buffer = ""
	m.mu.Unlock()
	m.ChatMemory.Clear()
	return nil
}
