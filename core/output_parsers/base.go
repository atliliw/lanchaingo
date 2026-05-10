package output_parsers

import (
	"encoding/json"
	"fmt"
)

// BaseOutputParser defines the interface for parsing LLM output.
// Maps to Rust langchainrust::core::output_parsers::base::BaseOutputParser.
type BaseOutputParser interface {
	// Parse parses the LLM output string into a structured result.
	Parse(text string) (any, error)

	// GetFormatInstructions returns instructions for the LLM
	// to generate output in the expected format.
	GetFormatInstructions() string

	// Type returns the parser type identifier.
	Type() string
}

// OutputParserError represents a parsing failure.
type OutputParserError struct {
	Message string
	Cause   error
}

func (e *OutputParserError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("output parser: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("output parser: %s", e.Message)
}

func (e *OutputParserError) Unwrap() error {
	return e.Cause
}

// NewOutputParserError creates a new OutputParserError.
func NewOutputParserError(msg string, cause error) *OutputParserError {
	return &OutputParserError{Message: msg, Cause: cause}
}

func isJSON(s string) bool {
	s = trimSpace(s)
	if len(s) == 0 {
		return false
	}
	return (s[0] == '{' && s[len(s)-1] == '}') ||
		(s[0] == '[' && s[len(s)-1] == ']')
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func tryParseJSON(text string) (any, error) {
	var result any
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return result, nil
}
