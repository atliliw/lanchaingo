package chains

import (
	"fmt"
	"strings"

	"github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// LLMChain is the most basic chain combining a prompt template with an LLM.
type LLMChain struct {
	llm             language_models.BaseChatModel
	PromptTemplate  string
	InputKey        string
	OutputKey       string
	name            string
}

func NewLLMChain(llm language_models.BaseChatModel, promptTemplate string) *LLMChain {
	return &LLMChain{
		llm:            llm,
		PromptTemplate: promptTemplate,
		InputKey:       "input",
		OutputKey:      "text",
		name:           "llm_chain",
	}
}

func (c *LLMChain) WithInputKey(key string) *LLMChain {
	c.InputKey = key
	return c
}

func (c *LLMChain) WithOutputKey(key string) *LLMChain {
	c.OutputKey = key
	return c
}

func (c *LLMChain) WithName(name string) *LLMChain {
	c.name = name
	return c
}

func (c *LLMChain) renderPrompt(inputs map[string]any) (string, error) {
	prompt := c.PromptTemplate
	for key, val := range inputs {
		placeholder := "{" + key + "}"
		prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("%v", val))
	}
	return prompt, nil
}

func (c *LLMChain) InputKeys() []string {
	return []string{c.InputKey}
}

func (c *LLMChain) OutputKeys() []string {
	return []string{c.OutputKey}
}

func (c *LLMChain) ValidateInputs(inputs map[string]any) error {
	for _, key := range c.InputKeys() {
		if _, ok := inputs[key]; !ok {
			return NewChainError(ErrMissingInput, "missing input: "+key, nil)
		}
	}
	return nil
}

func (c *LLMChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	prompt, err := c.renderPrompt(inputs)
	if err != nil {
		return nil, NewChainError(ErrExecution, "failed to render prompt", err)
	}

	msg := messages.NewHumanMessage(prompt)
	result, err := c.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return nil, NewChainError(ErrExecution, "LLM call failed", err)
	}

	return ChainResult{c.OutputKey: result.Content}, nil
}

func (c *LLMChain) Name() string {
	return c.name
}

// LLMChainBuilder provides a builder pattern for LLMChain.
type LLMChainBuilder struct {
	llm            language_models.BaseChatModel
	promptTemplate string
	inputKey       string
	outputKey      string
	name           string
}

func NewLLMChainBuilder(llm language_models.BaseChatModel, promptTemplate string) *LLMChainBuilder {
	return &LLMChainBuilder{
		llm:            llm,
		promptTemplate: promptTemplate,
		inputKey:       "input",
		outputKey:      "text",
		name:           "llm_chain",
	}
}

func (b *LLMChainBuilder) InputKey(key string) *LLMChainBuilder {
	b.inputKey = key
	return b
}

func (b *LLMChainBuilder) OutputKey(key string) *LLMChainBuilder {
	b.outputKey = key
	return b
}

func (b *LLMChainBuilder) Name(name string) *LLMChainBuilder {
	b.name = name
	return b
}

func (b *LLMChainBuilder) Build() *LLMChain {
	return &LLMChain{
		llm:            b.llm,
		PromptTemplate: b.promptTemplate,
		InputKey:       b.inputKey,
		OutputKey:      b.outputKey,
		name:           b.name,
	}
}
