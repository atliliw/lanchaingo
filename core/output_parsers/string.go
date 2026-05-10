package output_parsers

// StrOutputParser returns the LLM output as-is (identity parser).
// Maps to Rust langchainrust::core::output_parsers::str_parser::StrOutputParser.
type StrOutputParser struct{}

func NewStrOutputParser() *StrOutputParser {
	return &StrOutputParser{}
}

func (p *StrOutputParser) Parse(text string) (any, error) {
	return text, nil
}

func (p *StrOutputParser) GetFormatInstructions() string {
	return ""
}

func (p *StrOutputParser) Type() string {
	return "string"
}
