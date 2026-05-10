package output_parsers

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// TypedOutputParser parses LLM output into a specific Go type.
// Uses the type's JSON tags for field mapping.
// Maps to Rust langchainrust::core::output_parsers::typed_parser::TypedOutputParser.
type TypedOutputParser struct {
	targetType reflect.Type
	example    any
}

func NewTypedOutputParser(example any) *TypedOutputParser {
	return &TypedOutputParser{
		targetType: reflect.TypeOf(example),
		example:    example,
	}
}

func (p *TypedOutputParser) Parse(text string) (any, error) {
	if !isJSON(text) {
		return nil, NewOutputParserError("output is not valid JSON", nil)
	}

	target := reflect.New(p.targetType).Interface()
	if err := json.Unmarshal([]byte(text), target); err != nil {
		return nil, NewOutputParserError(
			fmt.Sprintf("failed to parse into %s", p.targetType.Name()), err)
	}

	return reflect.ValueOf(target).Elem().Interface(), nil
}

func (p *TypedOutputParser) GetFormatInstructions() string {
	exampleJSON, err := json.MarshalIndent(p.example, "", "  ")
	if err != nil {
		return fmt.Sprintf("Return a JSON object matching the %s type.", p.targetType.Name())
	}
	return fmt.Sprintf("Return a JSON object with this structure:\n%s", string(exampleJSON))
}

func (p *TypedOutputParser) Type() string {
	return fmt.Sprintf("typed(%s)", p.targetType.Name())
}
