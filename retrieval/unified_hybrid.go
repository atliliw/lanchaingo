package retrieval

import (
	"github.com/atliliw/lanchaingo/embeddings"
	"github.com/atliliw/lanchaingo/retrieval/bm25"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// HybridIndexConfig configures the unified hybrid index.
type HybridIndexConfig struct {
	ChunkSize      int
	ChunkOverlap   int
	BM25K          int
	VectorK        int
	RRFK           int
	MergeThreshold float64
}

func DefaultHybridIndexConfig() HybridIndexConfig {
	return HybridIndexConfig{
		ChunkSize:      500,
		ChunkOverlap:   50,
		BM25K:          10,
		VectorK:        10,
		RRFK:           60,
		MergeThreshold: 0.5,
	}
}

// UnifiedHybridIndex manages BM25 + vector indexes together.
type UnifiedHybridIndex struct {
	bm25       *bm25.ChunkedBM25Retriever
	store      vs.VectorStore
	embeddings embeddings.Embeddings
	config     HybridIndexConfig
}

func NewUnifiedHybridIndex(embeddings embeddings.Embeddings, store vs.VectorStore) *UnifiedHybridIndex {
	return &UnifiedHybridIndex{
		bm25:       bm25.NewChunkedBM25Retriever(),
		store:      store,
		embeddings: embeddings,
		config:     DefaultHybridIndexConfig(),
	}
}

func (idx *UnifiedHybridIndex) AddDocument(parent vs.Document, leaves []vs.Document, leafTerms [][]string) {
	idx.bm25.AddDocument(parent, leaves, leafTerms)
}

func (idx *UnifiedHybridIndex) Search(query string, k int) ([]RetrievedDocument, error) {
	bm25Docs := idx.bm25.Search(query, idx.config.BM25K)
	bm25Raw := make([]vs.Document, len(bm25Docs))
	for i, d := range bm25Docs {
		bm25Raw[i] = d.Document
	}

	qEmb, err := idx.embeddings.EmbedQuery(query)
	if err != nil {
		return nil, err
	}
	vectorResults, err := idx.store.SimilaritySearch(qEmb, idx.config.VectorK)
	if err != nil {
		return nil, err
	}
	vectorRaw := make([]vs.Document, len(vectorResults))
	for i, d := range vectorResults {
		vectorRaw[i] = d.Document
	}

	fused := ReciprocalRankFusion(bm25Raw, vectorRaw, idx.config.RRFK)
	if k < len(fused) {
		fused = fused[:k]
	}
	return fused, nil
}
