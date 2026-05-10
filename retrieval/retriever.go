package retrieval

import (
	vs "github.com/atliliw/lanchaingo/vector_stores"
	"github.com/atliliw/lanchaingo/embeddings"
)

// RetrieverTrait defines the interface for document retrieval.
type RetrieverTrait interface {
	Retrieve(query string, k int) ([]vs.Document, error)
	RetrieveWithScores(query string, k int) ([]vs.SearchResult, error)
	AddDocuments(documents []vs.Document) error
}

// SimilarityRetriever retrieves documents by embedding similarity.
type SimilarityRetriever struct {
	store      vs.VectorStore
	embeddings embeddings.Embeddings
}

func NewSimilarityRetriever(store vs.VectorStore, emb embeddings.Embeddings) *SimilarityRetriever {
	return &SimilarityRetriever{store: store, embeddings: emb}
}

func (r *SimilarityRetriever) Retrieve(query string, k int) ([]vs.Document, error) {
	results, err := r.RetrieveWithScores(query, k)
	if err != nil {
		return nil, err
	}
	docs := make([]vs.Document, len(results))
	for i, r := range results {
		docs[i] = r.Document
	}
	return docs, nil
}

func (r *SimilarityRetriever) RetrieveWithScores(query string, k int) ([]vs.SearchResult, error) {
	queryEmb, err := r.embeddings.EmbedQuery(query)
	if err != nil {
		return nil, NewRetrieverError(ErrEmbedding, "embed query failed", err)
	}

	results, err := r.store.SimilaritySearch(queryEmb, k)
	if err != nil {
		return nil, NewRetrieverError(ErrStore, "search failed", err)
	}
	if len(results) == 0 {
		return nil, NewRetrieverError(ErrNoResults, "no documents found", nil)
	}

	return results, nil
}

func (r *SimilarityRetriever) AddDocuments(documents []vs.Document) error {
	texts := make([]string, len(documents))
	for i, d := range documents {
		texts[i] = d.Content
	}

	embeddings, err := r.embeddings.EmbedDocuments(texts)
	if err != nil {
		return NewRetrieverError(ErrEmbedding, "embed documents failed", err)
	}

	_, err = r.store.AddDocuments(documents, embeddings)
	if err != nil {
		return NewRetrieverError(ErrStore, "add documents failed", err)
	}
	return nil
}

// Retriever is a type alias for SimilarityRetriever.
type Retriever = SimilarityRetriever
