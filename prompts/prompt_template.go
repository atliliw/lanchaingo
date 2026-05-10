package prompts

import (
	"fmt"
	"strings"
)

// PromptTemplate manages prompt templates with variable substitution.
// Maps to Rust langchainrust::prompts::prompt_template::PromptTemplate.
type PromptTemplate struct {
	template       string
	inputVariables []string
}

// NewPromptTemplate creates a new PromptTemplate.
// The template uses {variable} syntax for placeholders.
// inputVariables lists the variable names expected by the template.
func NewPromptTemplate(template string, inputVariables []string) *PromptTemplate {
	return &PromptTemplate{
		template:       template,
		inputVariables: inputVariables,
	}
}

// Format replaces {variable} placeholders with values and returns the result.
// Returns an error if any required variable is missing from values.
func (pt *PromptTemplate) Format(values map[string]string) (string, error) {
	if values == nil {
		return "", fmt.Errorf("prompt: values map is nil")
	}

	result := pt.template

	for _, v := range pt.inputVariables {
		if _, ok := values[v]; !ok {
			return "", fmt.Errorf("prompt: missing required variable %q", v)
		}
	}

	for key, val := range values {
		placeholder := "{" + key + "}"
		result = strings.ReplaceAll(result, placeholder, val)
	}

	return result, nil
}

// GetInputVariables returns the list of input variable names.
func (pt *PromptTemplate) GetInputVariables() []string {
	vars := make([]string, len(pt.inputVariables))
	copy(vars, pt.inputVariables)
	return vars
}

// GetTemplate returns the raw template string.
func (pt *PromptTemplate) GetTemplate() string {
	return pt.template
}
