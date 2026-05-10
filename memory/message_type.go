package memory

import "github.com/atliliw/lanchaingo/schema/messages"

// NewHumanMessage creates a human message (re-export for convenience).
func NewHumanMessage(content string) messages.Message {
	return messages.NewHumanMessage(content)
}

// NewAIMessage creates an AI message (re-export for convenience).
func NewAIMessage(content string) messages.Message {
	return messages.NewAIMessage(content)
}

// Message is a type alias for schema/messages.Message.
type Message = messages.Message
