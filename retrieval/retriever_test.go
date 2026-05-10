package retrieval

import (
	"testing"
	"github.com/atliliw/lanchaingo/embeddings"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

func TestSimilarityRetriever(t *testing.T) {
	store := vs.NewInMemoryVectorStore()
	emb := embeddings.NewMockEmbeddings(64)
	r := NewSimilarityRetriever(store, emb)
	r.AddDocuments([]vs.Document{vs.NewDocument("test doc")})
	if store.Count() != 1 { t.Errorf("expected 1 doc") }
	results, err := r.Retrieve("test", 2)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if len(results) == 0 { t.Error("expected at least 1 result") }
}
func TestRetrieveWithScores(t *testing.T) {
	store := vs.NewInMemoryVectorStore()
	emb := embeddings.NewMockEmbeddings(64)
	r := NewSimilarityRetriever(store, emb)
	r.AddDocuments([]vs.Document{vs.NewDocument("doc a"), vs.NewDocument("doc b")})
	results, err := r.RetrieveWithScores("query", 2)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if len(results) != 2 { t.Fatalf("expected 2 results") }
}
func TestReciprocalRankFusion(t *testing.T) {
	bm25 := []vs.Document{vs.NewDocument("d1").WithID("1"), vs.NewDocument("d2").WithID("2")}
	vector := []vs.Document{vs.NewDocument("d3").WithID("3"), vs.NewDocument("d1").WithID("1")}
	results := ReciprocalRankFusion(bm25, vector, RRFK)
	if len(results) == 0 { t.Fatal("expected non-empty results") }
}
func TestReciprocalRankFusionEmpty(t *testing.T) {
	results := ReciprocalRankFusion(nil, nil, RRFK)
	if len(results) != 0 { t.Errorf("expected empty, got %d", len(results)) }
}
func TestSplitterBasic(t *testing.T) {
	s := NewRecursiveCharacterSplitter(50, 0)
	chunks := s.SplitText("Hello world. This is a test.")
	if len(chunks) == 0 { t.Error("expected at least 1 chunk") }
}
func TestSplitterDocument(t *testing.T) {
	s := NewRecursiveCharacterSplitter(100, 0)
	doc := vs.NewDocument("Part one.\n\nPart two.").WithMetadata("source", "test")
	docs := s.SplitDocument(doc)
	if len(docs) < 2 { t.Error("expected at least 2 chunks") }
	for _, d := range docs {
		if d.Metadata["source"] != "test" { t.Errorf("expected source=test") }
	}
}
