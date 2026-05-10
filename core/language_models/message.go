package language_models

import "github.com/atliliw/lanchaingo/schema/messages"

// Message is re-exported from schema/messages for convenience.
type Message = messages.Message

// MessageType constants re-exported from schema/messages.
type MessageType = messages.MessageType

const (
	MessageTypeSystem = messages.MessageTypeSystem
	MessageTypeHuman  = messages.MessageTypeHuman
	MessageTypeAI     = messages.MessageTypeAI
	MessageTypeTool   = messages.MessageTypeTool
)

// Convenience message constructors.
var (
	NewSystemMessage = messages.NewSystemMessage
	NewHumanMessage  = messages.NewHumanMessage
	NewAIMessage     = messages.NewAIMessage
	NewToolMessage   = messages.NewToolMessage
)
