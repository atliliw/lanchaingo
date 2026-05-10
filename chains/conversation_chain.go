package chains

import (
	"github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/memory"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// ConversationChain is a chain with memory for multi-turn conversations.
type ConversationChain struct {
	llm          language_models.BaseChatModel
	memory       *memory.ConversationBufferMemory
	SystemPrompt string
	InputKey     string
	OutputKey    string
	MemoryKey    string
	name         string
	Verbose      bool
}

func NewConversationChain(llm language_models.BaseChatModel, mem *memory.ConversationBufferMemory) *ConversationChain {
	mem.WithReturnMessages(true)
	return &ConversationChain{
		llm:       llm,
		memory:    mem,
		InputKey:  "input",
		OutputKey: "output",
		MemoryKey: "history",
		name:      "conversation_chain",
	}
}

func (c *ConversationChain) WithSystemPrompt(prompt string) *ConversationChain {
	c.SystemPrompt = prompt
	return c
}

func (c *ConversationChain) WithInputKey(key string) *ConversationChain {
	c.InputKey = key
	return c
}

func (c *ConversationChain) WithOutputKey(key string) *ConversationChain {
	c.OutputKey = key
	return c
}

func (c *ConversationChain) WithMemoryKey(key string) *ConversationChain {
	c.MemoryKey = key
	return c
}

func (c *ConversationChain) WithName(name string) *ConversationChain {
	c.name = name
	return c
}

func (c *ConversationChain) WithVerbose(v bool) *ConversationChain {
	c.Verbose = v
	return c
}

func (c *ConversationChain) ClearMemory() error {
	return c.memory.Clear()
}

func (c *ConversationChain) Predict(input string) (string, error) {
	result, err := c.Invoke(map[string]any{c.InputKey: input})
	if err != nil {
		return "", err
	}
	out, _ := result[c.OutputKey].(string)
	return out, nil
}

func (c *ConversationChain) prepareMessages(input string, history []messages.Message) []messages.Message {
	var msgs []messages.Message
	if c.SystemPrompt != "" {
		msgs = append(msgs, messages.NewSystemMessage(c.SystemPrompt))
	}
	msgs = append(msgs, history...)
	msgs = append(msgs, messages.NewHumanMessage(input))
	return msgs
}

func (c *ConversationChain) InputKeys() []string {
	return []string{c.InputKey}
}

func (c *ConversationChain) OutputKeys() []string {
	return []string{c.OutputKey}
}

func (c *ConversationChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing input: "+c.InputKey, nil)
	}
	return nil
}

func (c *ConversationChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	input, _ := inputs[c.InputKey].(string)

	history := c.memory.ChatMemory.Messages()
	msgs := c.prepareMessages(input, history)

	result, err := c.llm.Chat(nil, msgs)
	if err != nil {
		return nil, NewChainError(ErrExecution, "LLM call failed", err)
	}

	_ = c.memory.SaveContext(
		map[string]string{c.InputKey: input},
		map[string]string{c.OutputKey: result.Content},
	)

	return ChainResult{c.OutputKey: result.Content}, nil
}

func (c *ConversationChain) Name() string {
	return c.name
}

// ConversationChainBuilder provides a builder for ConversationChain.
type ConversationChainBuilder struct {
	llm          language_models.BaseChatModel
	memory       *memory.ConversationBufferMemory
	systemPrompt string
	inputKey     string
	outputKey    string
	memoryKey    string
	name         string
	verbose      bool
}

func NewConversationChainBuilder(llm language_models.BaseChatModel) *ConversationChainBuilder {
	return &ConversationChainBuilder{
		llm:    llm,
		memory: memory.NewConversationBufferMemory(),
	}
}

func (b *ConversationChainBuilder) Memory(mem *memory.ConversationBufferMemory) *ConversationChainBuilder {
	b.memory = mem
	return b
}

func (b *ConversationChainBuilder) SystemPrompt(prompt string) *ConversationChainBuilder {
	b.systemPrompt = prompt
	return b
}

func (b *ConversationChainBuilder) InputKey(key string) *ConversationChainBuilder {
	b.inputKey = key
	return b
}

func (b *ConversationChainBuilder) OutputKey(key string) *ConversationChainBuilder {
	b.outputKey = key
	return b
}

func (b *ConversationChainBuilder) Verbose(v bool) *ConversationChainBuilder {
	b.verbose = v
	return b
}

func (b *ConversationChainBuilder) Build() *ConversationChain {
	chain := NewConversationChain(b.llm, b.memory)
	if b.systemPrompt != "" {
		chain.WithSystemPrompt(b.systemPrompt)
	}
	if b.inputKey != "" {
		chain.WithInputKey(b.inputKey)
	}
	if b.outputKey != "" {
		chain.WithOutputKey(b.outputKey)
	}
	if b.memoryKey != "" {
		chain.WithMemoryKey(b.memoryKey)
	}
	if b.name != "" {
		chain.WithName(b.name)
	}
	chain.Verbose = b.verbose
	return chain
}
