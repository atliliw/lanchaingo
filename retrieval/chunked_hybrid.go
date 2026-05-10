package retrieval

import (
	"github.com/atliliw/lanchaingo/embeddings"
	"github.com/atliliw/lanchaingo/retrieval/bm25"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// ChunkedHybridRetriever combines BM25 and vector search on chunked documents.
type ChunkedHybridRetriever struct {
	bm25     *bm25.ChunkedBM25Retriever
	embeddings embeddings.Embeddings
	store    vs.VectorStore
	bm25K    int
	vectorK  int
	rrfK     int
}

func NewChunkedHybridRetriever(b *bm25.ChunkedBM25Retriever, emb embeddings.Embeddings, store vs.VectorStore) *ChunkedHybridRetriever {
	return &ChunkedHybridRetriever{
		bm25:       b,
		embeddings: emb,
		store:      store,
		bm25K:      10,
		vectorK:    10,
		rrfK:       60,
	}
}

func (r *ChunkedHybridRetriever) Retrieve(query string, k int) ([]RetrievedDocument, error) {
	bm25Docs := r.bm25.Search(query, r.bm25K)
	bm25Raw := make([]vs.Document, len(bm25Docs))
	for i, d := range bm25Docs {
		bm25Raw[i] = d.Document
	}

	qEmb, err := r.embeddings.EmbedQuery(query)
	if err != nil {
		return nil, NewRetrieverError(ErrEmbedding, "embed query failed", err)
	}
	vectorResults, err := r.store.SimilaritySearch(qEmb, r.vectorK)
	if err != nil {
		return nil, NewRetrieverError(ErrStore, "vector search failed", err)
	}
	vectorRaw := make([]vs.Document, len(vectorResults))
	for i, d := range vectorResults {
		vectorRaw[i] = d.Document
	}

	fused := ReciprocalRankFusion(bm25Raw, vectorRaw, r.rrfK)
	if k < len(fused) {
		fused = fused[:k]
	}
	return fused, nil
}
