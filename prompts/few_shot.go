package prompts

import (
	"fmt"
	"strings"
)

// ExampleSelector selects which examples to include in a few-shot prompt.
type ExampleSelector interface {
	SelectExamples(input map[string]string) ([]map[string]string, error)
}

// LengthBasedExampleSelector selects examples that fit within a length limit.
type LengthBasedExampleSelector struct {
	examples  []map[string]string
	maxLength int
}

// NewLengthBasedExampleSelector creates a LengthBasedExampleSelector.
func NewLengthBasedExampleSelector(examples []map[string]string, maxLength int) *LengthBasedExampleSelector {
	if maxLength <= 0 {
		maxLength = 2048
	}
	return &LengthBasedExampleSelector{
		examples:  examples,
		maxLength: maxLength,
	}
}

func (s *LengthBasedExampleSelector) SelectExamples(input map[string]string) ([]map[string]string, error) {
	var selected []map[string]string
	totalLength := 0

	for _, example := range s.examples {
		exampleStr := exampleToString(example)
		if totalLength+len(exampleStr) > s.maxLength {
			break
		}
		selected = append(selected, example)
		totalLength += len(exampleStr)
	}

	return selected, nil
}

// AddExample appends a new example to the selector.
func (s *LengthBasedExampleSelector) AddExample(example map[string]string) {
	s.examples = append(s.examples, example)
}

// FewShotPromptTemplate constructs a prompt with few-shot examples.
type FewShotPromptTemplate struct {
	examples         []map[string]string
	examplePrompt    *PromptTemplate
	suffix           string
	inputVariables   []string
	exampleSeparator string
	prefix           string
}

// NewFewShotPromptTemplate creates a FewShotPromptTemplate.
func NewFewShotPromptTemplate(
	examples []map[string]string,
	examplePrompt *PromptTemplate,
	suffix string,
	inputVariables []string,
	exampleSeparator string,
	prefix string,
) *FewShotPromptTemplate {
	if exampleSeparator == "" {
		exampleSeparator = "\n\n"
	}
	return &FewShotPromptTemplate{
		examples:         examples,
		examplePrompt:    examplePrompt,
		suffix:           suffix,
		inputVariables:   inputVariables,
		exampleSeparator: exampleSeparator,
		prefix:           prefix,
	}
}

// Format constructs the full prompt by combining prefix, examples, and suffix.
func (fpt *FewShotPromptTemplate) Format(values map[string]string) (string, error) {
	var parts []string

	if fpt.prefix != "" {
		prefixResult := fpt.prefix
		for key, val := range values {
			prefixResult = strings.ReplaceAll(prefixResult, "{"+key+"}", val)
		}
		parts = append(parts, prefixResult)
	}

	for _, example := range fpt.examples {
		formatted, err := fpt.examplePrompt.Format(example)
		if err != nil {
			return "", fmt.Errorf("few-shot: failed to format example: %w", err)
		}
		parts = append(parts, formatted)
	}

	suffixResult := fpt.suffix
	for key, val := range values {
		suffixResult = strings.ReplaceAll(suffixResult, "{"+key+"}", val)
	}
	parts = append(parts, suffixResult)

	return strings.Join(parts, fpt.exampleSeparator), nil
}

func exampleToString(example map[string]string) string {
	var parts []string
	for k, v := range example {
		parts = append(parts, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(parts, "\n")
}
