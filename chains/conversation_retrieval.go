package chains

import (
	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/memory"
	"github.com/atliliw/lanchaingo/prompts"
	"github.com/atliliw/lanchaingo/retrieval"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// ConversationRetrievalChain combines retrieval with conversation memory.
type ConversationRetrievalChain struct {
	llm       lm.BaseChatModel
	retriever retrieval.RetrieverTrait
	memory    *memory.ConversationBufferMemory
	template  *prompts.PromptTemplate
	InputKey  string
	OutputKey string
	name      string
}

func NewConversationRetrievalChain(llm lm.BaseChatModel, retriever retrieval.RetrieverTrait, mem *memory.ConversationBufferMemory) *ConversationRetrievalChain {
	return &ConversationRetrievalChain{
		llm:       llm,
		retriever: retriever,
		memory:    mem,
		template:  prompts.NewPromptTemplate("Use context to answer.\nContext: {context}\nChat history:\n{history}\nQuestion: {question}", []string{"context", "history", "question"}),
		InputKey:  "query",
		OutputKey: "result",
		name:      "conversation_retrieval_chain",
	}
}

func (c *ConversationRetrievalChain) InputKeys() []string  { return []string{c.InputKey} }
func (c *ConversationRetrievalChain) OutputKeys() []string { return []string{c.OutputKey} }
func (c *ConversationRetrievalChain) Name() string         { return c.name }
func (c *ConversationRetrievalChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing: "+c.InputKey, nil)
	}
	return nil
}

func (c *ConversationRetrievalChain) Invoke(inputs map[string]any) (ChainResult, error) {
	query, _ := inputs[c.InputKey].(string)

	docs, err := c.retriever.Retrieve(query, 4)
	if err != nil {
		return nil, NewChainError(ErrExecution, "retrieval failed", err)
	}

	var contextStr string
	for i, d := range docs {
		if i > 0 { contextStr += "\n\n" }
		contextStr += d.Content
	}

	memVars, _ := c.memory.LoadMemoryVariables(map[string]string{"input": query})
	history, _ := memVars["history"].(string)

	msg := messages.NewHumanMessage(contextStr + "\n\nHistory: " + history + "\n\nQuestion: " + query)
	result, err := c.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return nil, NewChainError(ErrExecution, "LLM call failed", err)
	}

	c.memory.SaveContext(map[string]string{"input": query}, map[string]string{"output": result.Content})

	return ChainResult{c.OutputKey: result.Content}, nil
}
