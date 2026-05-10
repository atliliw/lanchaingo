package prompts

import (
	"fmt"
	"strings"

	"github.com/atliliw/lanchaingo/schema/messages"
)

// messageTemplate pairs a message role with a prompt template string.
type messageTemplate struct {
	Role     messages.MessageType
	Template string
}

// ChatPromptTemplate manages a list of message templates for chat models.
// Maps to Rust langchainrust::prompts::chat_prompt_template::ChatPromptTemplate.
type ChatPromptTemplate struct {
	messageTemplates []messageTemplate
}

// NewChatPromptTemplate creates an empty ChatPromptTemplate.
func NewChatPromptTemplate() *ChatPromptTemplate {
	return &ChatPromptTemplate{}
}

// AddMessage appends a message template to the list.
func (cpt *ChatPromptTemplate) AddMessage(role messages.MessageType, template string) {
	cpt.messageTemplates = append(cpt.messageTemplates, messageTemplate{
		Role:     role,
		Template: template,
	})
}

// AddSystemMessage adds a system message template.
func (cpt *ChatPromptTemplate) AddSystemMessage(template string) {
	cpt.AddMessage(messages.MessageTypeSystem, template)
}

// AddHumanMessage adds a human/user message template.
func (cpt *ChatPromptTemplate) AddHumanMessage(template string) {
	cpt.AddMessage(messages.MessageTypeHuman, template)
}

// AddAIMessage adds an AI/assistant message template.
func (cpt *ChatPromptTemplate) AddAIMessage(template string) {
	cpt.AddMessage(messages.MessageTypeAI, template)
}

// FormatMessages formats all message templates with the given values
// and returns a slice of Messages ready for LLM invocation.
func (cpt *ChatPromptTemplate) FormatMessages(values map[string]string) ([]messages.Message, error) {
	if values == nil {
		values = make(map[string]string)
	}

	result := make([]messages.Message, 0, len(cpt.messageTemplates))

	for _, mt := range cpt.messageTemplates {
		content := mt.Template
		for key, val := range values {
			placeholder := "{" + key + "}"
			content = strings.ReplaceAll(content, placeholder, val)
		}

		msg, err := newMessageForRole(mt.Role, content)
		if err != nil {
			return nil, fmt.Errorf("chat prompt: %w", err)
		}
		result = append(result, msg)
	}

	return result, nil
}

func newMessageForRole(role messages.MessageType, content string) (messages.Message, error) {
	switch role {
	case messages.MessageTypeSystem:
		return messages.NewSystemMessage(content), nil
	case messages.MessageTypeHuman:
		return messages.NewHumanMessage(content), nil
	case messages.MessageTypeAI:
		return messages.NewAIMessage(content), nil
	default:
		return messages.Message{}, fmt.Errorf("unsupported message type: %v", role)
	}
}
