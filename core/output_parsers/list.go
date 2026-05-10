package output_parsers

import "strings"

// CommaSeparatedListOutputParser parses comma-separated values from LLM output.
// Maps to Rust langchainrust::core::output_parsers::list_parser::CommaSeparatedListOutputParser.
type CommaSeparatedListOutputParser struct {
	Separator string
}

func NewCommaSeparatedListOutputParser() *CommaSeparatedListOutputParser {
	return &CommaSeparatedListOutputParser{Separator: ","}
}

func (p *CommaSeparatedListOutputParser) Parse(text string) (any, error) {
	parts := strings.Split(text, p.Separator)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result, nil
}

func (p *CommaSeparatedListOutputParser) GetFormatInstructions() string {
	return "Return a comma-separated list of values."
}

func (p *CommaSeparatedListOutputParser) Type() string {
	return "comma_separated_list"
}
