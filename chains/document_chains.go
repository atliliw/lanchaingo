package chains

import (
	"fmt"
	"sort"
	"strings"

	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// StuffDocumentsChain stuffs all documents into a single prompt.
type StuffDocumentsChain struct {
	llm              lm.BaseChatModel
	DocumentPrompt   string
	DocumentSeparator string
	InputKey         string
	OutputKey        string
	name             string
}

func NewStuffDocumentsChain(llm lm.BaseChatModel, documentPrompt string) *StuffDocumentsChain {
	return &StuffDocumentsChain{
		llm:               llm,
		DocumentPrompt:    documentPrompt,
		DocumentSeparator: "\n\n",
		InputKey:          "input_documents",
		OutputKey:         "output_text",
		name:              "stuff_documents_chain",
	}
}

func (c *StuffDocumentsChain) Name() string { return c.name }
func (c *StuffDocumentsChain) InputKeys() []string  { return []string{c.InputKey} }
func (c *StuffDocumentsChain) OutputKeys() []string { return []string{c.OutputKey} }

func (c *StuffDocumentsChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing: "+c.InputKey, nil)
	}
	return nil
}

func (c *StuffDocumentsChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	docs, ok := inputs[c.InputKey].([]string)
	if !ok {
		return nil, NewChainError(ErrExecution, "input_documents must be []string", nil)
	}

	content := strings.Join(docs, c.DocumentSeparator)
	msg := messages.NewHumanMessage(content)

	result, err := c.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return nil, NewChainError(ErrExecution, "LLM call failed", err)
	}

	return ChainResult{c.OutputKey: result.Content}, nil
}

// RefineDocumentsChain iteratively refines answer with each document.
type RefineDocumentsChain struct {
	llm                   lm.BaseChatModel
	InitialPrompt         string
	RefinePrompt          string
	DocumentSeparator     string
	InputKey              string
	OutputKey             string
	name                  string
}

func NewRefineDocumentsChain(llm lm.BaseChatModel, initialPrompt, refinePrompt string) *RefineDocumentsChain {
	return &RefineDocumentsChain{
		llm:               llm,
		InitialPrompt:     initialPrompt,
		RefinePrompt:      refinePrompt,
		DocumentSeparator: "\n\n",
		InputKey:          "input_documents",
		OutputKey:         "output_text",
		name:              "refine_documents_chain",
	}
}

func (c *RefineDocumentsChain) Name() string { return c.name }
func (c *RefineDocumentsChain) InputKeys() []string  { return []string{c.InputKey} }
func (c *RefineDocumentsChain) OutputKeys() []string { return []string{c.OutputKey} }
func (c *RefineDocumentsChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing: "+c.InputKey, nil)
	}
	return nil
}

func (c *RefineDocumentsChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	docs, ok := inputs[c.InputKey].([]string)
	if !ok {
		return nil, NewChainError(ErrExecution, "input_documents must be []string", nil)
	}

	var output string
	for i, doc := range docs {
		if i == 0 {
			msg := messages.NewHumanMessage(c.InitialPrompt + "\n" + doc)
			result, err := c.llm.Chat(nil, []messages.Message{msg})
			if err != nil {
				return nil, NewChainError(ErrExecution, "LLM call failed", err)
			}
			output = result.Content
		} else {
			prompt := fmt.Sprintf("%s\n\nPrevious answer: %s\n\nNew document: %s", c.RefinePrompt, output, doc)
			msg := messages.NewHumanMessage(prompt)
			result, err := c.llm.Chat(nil, []messages.Message{msg})
			if err != nil {
				return nil, NewChainError(ErrExecution, "LLM call failed", err)
			}
			output = result.Content
		}
	}

	return ChainResult{c.OutputKey: output}, nil
}

// MapReduceDocumentsChain maps over documents then reduces results.
type MapReduceDocumentsChain struct {
	llm               lm.BaseChatModel
	MapPrompt         string
	ReducePrompt      string
	InputKey          string
	OutputKey         string
	name              string
}

func NewMapReduceDocumentsChain(llm lm.BaseChatModel, mapPrompt, reducePrompt string) *MapReduceDocumentsChain {
	return &MapReduceDocumentsChain{
		llm:          llm,
		MapPrompt:    mapPrompt,
		ReducePrompt: reducePrompt,
		InputKey:     "input_documents",
		OutputKey:    "output_text",
		name:         "map_reduce_documents_chain",
	}
}

func (c *MapReduceDocumentsChain) InputKeys() []string  { return []string{c.InputKey} }
func (c *MapReduceDocumentsChain) OutputKeys() []string { return []string{c.OutputKey} }
func (c *MapReduceDocumentsChain) Name() string { return c.name }
func (c *MapReduceDocumentsChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing: "+c.InputKey, nil)
	}
	return nil
}

func (c *MapReduceDocumentsChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}

	docs, ok := inputs[c.InputKey].([]string)
	if !ok {
		return nil, NewChainError(ErrExecution, "input_documents must be []string", nil)
	}

	var summaries []string
	for _, doc := range docs {
		msg := messages.NewHumanMessage(c.MapPrompt + "\n" + doc)
		result, err := c.llm.Chat(nil, []messages.Message{msg})
		if err != nil {
			return nil, NewChainError(ErrExecution, "map LLM call failed", err)
		}
		summaries = append(summaries, result.Content)
	}

	combined := strings.Join(summaries, "\n\n")
	msg := messages.NewHumanMessage(c.ReducePrompt + "\n" + combined)
	result, err := c.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return nil, NewChainError(ErrExecution, "reduce LLM call failed", err)
	}

	return ChainResult{c.OutputKey: result.Content}, nil
}

// MapRerankDocumentsChain maps over documents and reranks results.
type MapRerankDocumentsChain struct {
	llm          lm.BaseChatModel
	MapPrompt    string
	InputKey     string
	OutputKey    string
	TopK         int
	name         string
}

func NewMapRerankDocumentsChain(llm lm.BaseChatModel, mapPrompt string) *MapRerankDocumentsChain {
	return &MapRerankDocumentsChain{
		llm:       llm,
		MapPrompt: mapPrompt,
		InputKey:  "input_documents",
		OutputKey: "output_text",
		TopK:      1,
		name:      "map_rerank_documents_chain",
	}
}

func (c *MapRerankDocumentsChain) InputKeys() []string  { return []string{c.InputKey} }
func (c *MapRerankDocumentsChain) OutputKeys() []string { return []string{c.OutputKey} }
func (c *MapRerankDocumentsChain) Name() string { return c.name }
func (c *MapRerankDocumentsChain) ValidateInputs(inputs map[string]any) error {
	if _, ok := inputs[c.InputKey]; !ok {
		return NewChainError(ErrMissingInput, "missing: "+c.InputKey, nil)
	}
	return nil
}

func (c *MapRerankDocumentsChain) Invoke(inputs map[string]any) (ChainResult, error) {
	if err := c.ValidateInputs(inputs); err != nil {
		return nil, err
	}
	docs, ok := inputs[c.InputKey].([]string)
	if !ok || len(docs) == 0 {
		return nil, NewChainError(ErrExecution, "input_documents must be non-empty []string", nil)
	}

	type scored struct {
		text  string
		score float64
	}
	var scoredResults []scored

	for _, doc := range docs {
		prompt := fmt.Sprintf("%s\n\nDocument: %s\nScore (0-10):", c.MapPrompt, doc)
		msg := messages.NewHumanMessage(prompt)
		result, err := c.llm.Chat(nil, []messages.Message{msg})
		if err != nil {
			continue
		}
		// Crude score extraction from LLM response
		score := 5.0
		fmt.Sscanf(result.Content, "%f", &score)
		scoredResults = append(scoredResults, scored{text: doc, score: score})
	}

	sort.Slice(scoredResults, func(i, j int) bool {
		return scoredResults[i].score > scoredResults[j].score
	})

	if c.TopK > len(scoredResults) {
		c.TopK = len(scoredResults)
	}
	scoredResults = scoredResults[:c.TopK]

	var output string
	for _, s := range scoredResults {
		output += fmt.Sprintf("Score %.1f: %s\n", s.score, s.text)
	}

	return ChainResult{c.OutputKey: output}, nil
}
