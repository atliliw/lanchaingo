package retrieval

import (
	"fmt"

	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

const defaultMultiQueryPrompt = `You are an AI assistant. Generate {num_queries} different versions of the given question to retrieve relevant documents. Provide each query on a new line.

Original question: {question}`

// MultiQueryRetriever generates multiple query variants to improve recall.
type MultiQueryRetriever struct {
	llm           lm.BaseChatModel
	retriever     RetrieverTrait
	NumQueries    int
	KPerQuery     int
	FinalK        int
	Prompt        string
}

func NewMultiQueryRetriever(llm lm.BaseChatModel, retriever RetrieverTrait) *MultiQueryRetriever {
	return &MultiQueryRetriever{
		llm:        llm,
		retriever:  retriever,
		NumQueries: 3,
		KPerQuery:  5,
		FinalK:     10,
		Prompt:     defaultMultiQueryPrompt,
	}
}

func (m *MultiQueryRetriever) generateQueries(query string) ([]string, error) {
	prompt := fmt.Sprintf(m.Prompt+"\n\nQuestion: %s", query)
	prompt = prompt + query
	msg := messages.NewHumanMessage(prompt)
	result, err := m.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return nil, fmt.Errorf("multi_query: llm failed: %w", err)
	}
	// Simple heuristic: one query per line, skip empty
	var queries []string
	current := ""
	for _, c := range result.Content {
		if c == '\n' {
			if current != "" {
				queries = append(queries, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		queries = append(queries, current)
	}
	if len(queries) == 0 {
		queries = []string{query}
	}
	if len(queries) > m.NumQueries {
		queries = queries[:m.NumQueries]
	}
	return queries, nil
}

func (m *MultiQueryRetriever) Retrieve(query string, k int) ([]vs.Document, error) {
	return m.retriever.Retrieve(query, k)
}

func (m *MultiQueryRetriever) RetrieveWithScores(query string, k int) ([]vs.SearchResult, error) {
	return m.retriever.RetrieveWithScores(query, k)
}

func (m *MultiQueryRetriever) AddDocuments(docs []vs.Document) error {
	return m.retriever.AddDocuments(docs)
}

// RetrieveWithMultiQuery generates variants and merges results by unique ID.
func (m *MultiQueryRetriever) RetrieveWithMultiQuery(query string) ([]vs.Document, error) {
	queries, err := m.generateQueries(query)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var all []vs.Document
	for _, q := range queries {
		docs, err := m.retriever.Retrieve(q, m.KPerQuery)
		if err != nil {
			continue
		}
		for _, d := range docs {
			id := d.ID
			if id == "" {
				id = d.Content
			}
			if !seen[id] {
				seen[id] = true
				all = append(all, d)
			}
		}
	}
	if len(all) > m.FinalK {
		all = all[:m.FinalK]
	}
	return all, nil
}
