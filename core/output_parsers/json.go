package output_parsers

import (
	"fmt"
)

// JsonOutputParser parses LLM output as JSON.
// Maps to Rust langchainrust::core::output_parsers::json_parser::JsonOutputParser.
type JsonOutputParser struct{}

func NewJsonOutputParser() *JsonOutputParser {
	return &JsonOutputParser{}
}

func (p *JsonOutputParser) Parse(text string) (any, error) {
	result, err := tryParseJSON(text)
	if err != nil {
		return nil, NewOutputParserError("failed to parse JSON output", err)
	}
	return result, nil
}

func (p *JsonOutputParser) GetFormatInstructions() string {
	return "Return your response as valid JSON."
}

func (p *JsonOutputParser) Type() string {
	return "json"
}

// ValidateJSON checks if a string is valid JSON.
func ValidateJSON(text string) error {
	_, err := tryParseJSON(text)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}
