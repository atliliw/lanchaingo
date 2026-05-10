package vector_stores

import "testing"

func TestAddAndSearch(t *testing.T) {
	s := NewInMemoryVectorStore()
	docs := []Document{NewDocument("Rust programming"), NewDocument("Scripting scripting")}
	embeddings := [][]float32{{1.0, 0.0, 0.0}, {0.0, 1.0, 0.0}}
	ids, err := s.AddDocuments(docs, embeddings)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if len(ids) != 2 { t.Fatalf("expected 2 ids") }
	if s.Count() != 2 { t.Errorf("expected 2") }
	results, err := s.SimilaritySearch([]float32{0.9, 0.1, 0.0}, 1)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if len(results) != 1 { t.Fatalf("expected 1 result") }
}
func TestGetAndDelete(t *testing.T) {
	s := NewInMemoryVectorStore()
	s.AddDocuments([]Document{NewDocument("test").WithID("id-1")}, [][]float32{{1.0}})
	got, err := s.GetDocument("id-1")
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if got.Content != "test" { t.Errorf("expected test") }
	s.DeleteDocument("id-1")
	if s.Count() != 0 { t.Errorf("expected 0") }
}
func TestClear(t *testing.T) {
	s := NewInMemoryVectorStore()
	s.AddDocuments([]Document{NewDocument("a"), NewDocument("b")}, [][]float32{{1.0},{0.0}})
	s.Clear()
	if s.Count() != 0 { t.Errorf("expected 0") }
}
func TestDocumentCreation(t *testing.T) {
	doc := NewDocument("content").WithID("d1").WithMetadata("key","val")
	if doc.Content != "content" || doc.ID != "d1" { t.Errorf("unexpected document") }
}
