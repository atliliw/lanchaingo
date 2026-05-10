package chains

import (
	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/prompts"
	"github.com/atliliw/lanchaingo/retrieval"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// RetrievalQA answers questions using retrieved documents.
type RetrievalQA struct {
	llm       lm.BaseChatModel
	retriever retrieval.RetrieverTrait
	template  *prompts.PromptTemplate
	InputKey  string
	OutputKey string
	name      string
}

func NewRetrievalQA(llm lm.BaseChatModel, retriever retrieval.RetrieverTrait) *RetrievalQA {
	return &RetrievalQA{
		llm:       llm,
		retriever: retriever,
		template:  prompts.NewPromptTemplate("Use the following context to answer the question.\nContext: {context}\nQuestion: {question}", []string{"context", "question"}),
		InputKey:  "query",
		OutputKey: "result",
		name:      "retrieval_qa",
	}
}

func (c *RetrievalQA) InputKeys() []string  { return []string{c.InputKey} }
func (c *RetrievalQA) OutputKeys() []string { return []string{c.OutputKey} }
func (c *RetrievalQA) Name() string         { return c.name }

func (c *RetrievalQA) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing: "+c.InputKey, nil)
	}
	return nil
}

func (c *RetrievalQA) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	query, _ := inputs[c.InputKey].(string)
	docs, err := c.retriever.Retrieve(query, 4)
	if err != nil {
		return nil, NewChainError(ErrExecution, "retrieval failed", err)
	}

	var contextStr string
	for i, d := range docs {
		if i > 0 {
			contextStr += "\n\n"
		}
		contextStr += d.Content
	}

	prompt, err := c.template.Format(map[string]string{"context": contextStr, "question": query})
	if err != nil {
		return nil, NewChainError(ErrExecution, "prompt format failed", err)
	}

	msg := messages.NewHumanMessage(prompt)
	result, err := c.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return nil, NewChainError(ErrExecution, "LLM call failed", err)
	}

	return ChainResult{c.OutputKey: result.Content}, nil
}
