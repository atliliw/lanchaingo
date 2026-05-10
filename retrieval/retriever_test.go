package retrieval

import (
	"testing"

	"github.com/atliliw/lanchaingo/embeddings"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// 娴嬭瘯 SimilarityRetriever 瀹屾暣娴佺▼
func TestSimilarityRetriever(t *testing.T) {
	store := vs.NewInMemoryVectorStore()
	emb := embeddings.NewMockEmbeddings(64)
	retriever := NewSimilarityRetriever(store, emb)

	docs := []vs.Document{
		vs.NewDocument("Rust is a systems programming language"),
		vs.NewDocument("Python is a scripting language"),
	}

	err := retriever.AddDocuments(docs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if store.Count() != 2 {
		t.Errorf("expected 2 docs, got %d", store.Count())
	}

	results, err := retriever.Retrieve("programming", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) == 0 {
		t.Error("expected at least 1 result")
	}
}

// 娴嬭瘯 RetrieveWithScores 杩斿洖鍒嗘暟
func TestRetrieveWithScores(t *testing.T) {
	store := vs.NewInMemoryVectorStore()
	emb := embeddings.NewMockEmbeddings(64)
	retriever := NewSimilarityRetriever(store, emb)

	retriever.AddDocuments([]vs.Document{
		vs.NewDocument("doc a"),
		vs.NewDocument("doc b"),
	})

	results, err := retriever.RetrieveWithScores("query", 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Score < -1 || results[0].Score > 1 {
		t.Errorf("invalid score: %f", results[0].Score)
	}
}
