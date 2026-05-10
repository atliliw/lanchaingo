package tools

import (
	"context"
	"encoding/json"
)

// StructuredTool is a tool that accepts structured JSON input.
// Use this for tools that need typed parameters beyond a simple string.
type StructuredTool struct {
	name        string
	description string
	inputSchema map[string]any
	inputType   string
	execute     func(ctx context.Context, params json.RawMessage) (string, error)
}

// NewStructuredTool creates a StructuredTool with the given parameters.
func NewStructuredTool(
	name, description, inputType string,
	inputSchema map[string]any,
	execute func(ctx context.Context, params json.RawMessage) (string, error),
) *StructuredTool {
	if inputSchema == nil {
		inputSchema = make(map[string]any)
	}
	return &StructuredTool{
		name:        name,
		description: description,
		inputSchema: inputSchema,
		inputType:   inputType,
		execute:     execute,
	}
}

func (t *StructuredTool) Name() string {
	return t.name
}

func (t *StructuredTool) Description() string {
	return t.description
}

func (t *StructuredTool) Run(ctx context.Context, input string) (string, error) {
	var params json.RawMessage
	if input != "" {
		params = json.RawMessage(input)
	}
	return t.execute(ctx, params)
}

func (t *StructuredTool) InputSchema() map[string]any {
	return t.inputSchema
}

func (t *StructuredTool) InputType() string {
	return t.inputType
}
