package output_parsers

import (
	"encoding/json"
	"fmt"
)

// StructuredOutputParser parses LLM output into a predefined structure.
// Maps to Rust langchainrust::core::output_parsers::structured_parser::StructuredOutputParser.
type StructuredOutputParser struct {
	Schema map[string]any
}

func NewStructuredOutputParser(schema map[string]any) *StructuredOutputParser {
	return &StructuredOutputParser{Schema: schema}
}

// NewStructuredOutputParserFromFields creates a parser from field name -> description mappings.
func NewStructuredOutputParserFromFields(fields map[string]string) *StructuredOutputParser {
	schema := make(map[string]any)
	for name, desc := range fields {
		schema[name] = map[string]any{
			"description": desc,
			"type":        "string",
		}
	}
	return &StructuredOutputParser{Schema: schema}
}

func (p *StructuredOutputParser) Parse(text string) (any, error) {
	var result map[string]any
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, NewOutputParserError("failed to parse structured output as JSON", err)
	}

	for key := range p.Schema {
		if _, exists := result[key]; !exists {
			return nil, NewOutputParserError(
				fmt.Sprintf("missing required field %q in structured output", key), nil)
		}
	}

	return result, nil
}

func (p *StructuredOutputParser) GetFormatInstructions() string {
	jsonSchema, _ := json.MarshalIndent(p.Schema, "", "  ")
	return fmt.Sprintf("Return a JSON object with the following schema:\n%s", string(jsonSchema))
}

func (p *StructuredOutputParser) Type() string {
	return "structured"
}
