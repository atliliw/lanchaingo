package tools

import (
	"encoding/json"
	"fmt"
)

// StructuredOutput defines the output schema for structured tool responses.
// Maps to Rust langchainrust::core::tools::structured_output::StructuredOutput.
type StructuredOutput struct {
	Schema   map[string]any `json:"schema"`
	DataType string         `json:"data_type"`
	Example  any            `json:"example,omitempty"`
}

// NewStructuredOutput creates a StructuredOutput from a Go struct.
// The struct should have JSON tags for field mapping.
func NewStructuredOutput(example any) (*StructuredOutput, error) {
	data, err := json.Marshal(example)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal example: %w", err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to generate schema: %w", err)
	}

	return &StructuredOutput{
		Schema:   schema,
		DataType: fmt.Sprintf("%T", example),
		Example:  example,
	}, nil
}
