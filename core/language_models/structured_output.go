package language_models

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/atliliw/lanchaingo/schema/messages"
)

// StructuredOutputModel wraps a BaseChatModel to provide typed structured output.
// Usage:
//
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	model := NewStructuredOutputModel[Person](llm)
//	person, err := model.Invoke(ctx, "张三今年30岁")
type StructuredOutputModel[T any] struct {
	llm    BaseChatModel
	schema string
}

// NewStructuredOutputModel creates a model that returns structured output of type T.
// It injects a JSON Schema constraint into the system prompt to force the LLM
// to return valid JSON matching the type.
func NewStructuredOutputModel[T any](llm BaseChatModel) *StructuredOutputModel[T] {
	var zero T
	schema := generateJSONSchema(zero)
	return &StructuredOutputModel[T]{
		llm:    llm,
		schema: schema,
	}
}

// generateJSONSchema creates a simple JSON schema from a Go type.
// Uses json.Marshal on an empty instance to infer the structure.
func generateJSONSchema(v any) string {
	data, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(data)
}

// Invoke sends a query and returns typed structured output.
func (m *StructuredOutputModel[T]) Invoke(ctx context.Context, query string) (T, error) {
	var zero T

	systemPrompt := fmt.Sprintf(
		"You must respond with valid JSON matching this schema: %s\n"+
			"Return ONLY the JSON object, no other text.", m.schema)

	msgs := []messages.Message{
		messages.NewSystemMessage(systemPrompt),
		messages.NewHumanMessage(query),
	}

	result, err := m.llm.Chat(ctx, msgs)
	if err != nil {
		return zero, fmt.Errorf("structured_output: llm failed: %w", err)
	}

	var output T
	if err := json.Unmarshal([]byte(result.Content), &output); err != nil {
		return zero, fmt.Errorf("structured_output: parse failed: %w", err)
	}
	return output, nil
}
