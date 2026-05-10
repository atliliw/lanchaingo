package retrieval

import (
	"fmt"

	emb "github.com/atliliw/lanchaingo/embeddings"
	lm "github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/schema/messages"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

const defaultHydePrompt = `Please write a passage to answer the following question:
Question: {question}
Passage:`

// HyDERetriever uses LLM-generated hypothetical documents for retrieval.
type HyDERetriever struct {
	llm                   lm.BaseChatModel
	embeddings            emb.Embeddings
	store                 vs.VectorStore
	Prompt                string
	K                     int
	IncludeOriginalQuery  bool
}

func NewHyDERetriever(llm lm.BaseChatModel, embeddings emb.Embeddings, store vs.VectorStore) *HyDERetriever {
	return &HyDERetriever{
		llm:                  llm,
		embeddings:           embeddings,
		store:                store,
		Prompt:               defaultHydePrompt,
		K:                    5,
		IncludeOriginalQuery: true,
	}
}

func (h *HyDERetriever) generateHypotheticalDoc(query string) (string, error) {
	prompt := fmt.Sprintf(h.Prompt, query)
	msg := messages.NewHumanMessage(prompt)
	result, err := h.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return "", fmt.Errorf("hyde: llm failed: %w", err)
	}
	return result.Content, nil
}

func (h *HyDERetriever) Retrieve(query string, k int) ([]vs.Document, error) {
	return h.search(query, k)
}

func (h *HyDERetriever) RetrieveWithScores(query string, k int) ([]vs.SearchResult, error) {
	hdoc, err := h.generateHypotheticalDoc(query)
	if err != nil {
		return nil, err
	}

	hdocEmb, err := h.embeddings.EmbedQuery(hdoc)
	if err != nil {
		return nil, fmt.Errorf("hyde: embed failed: %w", err)
	}

	results, err := h.store.SimilaritySearch(hdocEmb, k)
	if err != nil {
		return nil, err
	}

	if h.IncludeOriginalQuery {
		qEmb, err := h.embeddings.EmbedQuery(query)
		if err == nil {
			qResults, err := h.store.SimilaritySearch(qEmb, k)
			if err == nil {
				seen := make(map[string]bool)
				for _, r := range results {
					seen[r.Document.ID] = true
				}
				for _, r := range qResults {
					if !seen[r.Document.ID] {
						results = append(results, r)
						seen[r.Document.ID] = true
					}
				}
			}
		}
	}

	return results, nil
}

func (h *HyDERetriever) AddDocuments(docs []vs.Document) error {
	texts := make([]string, len(docs))
	for i, d := range docs {
		texts[i] = d.Content
	}
	vecs, err := h.embeddings.EmbedDocuments(texts)
	if err != nil {
		return err
	}
	_, err = h.store.AddDocuments(docs, vecs)
	return err
}

func (h *HyDERetriever) search(query string, k int) ([]vs.Document, error) {
	results, err := h.RetrieveWithScores(query, k)
	if err != nil {
		return nil, err
	}
	docs := make([]vs.Document, len(results))
	for i, r := range results {
		docs[i] = r.Document
	}
	return docs, nil
}
