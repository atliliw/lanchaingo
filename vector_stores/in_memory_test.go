package vector_stores

import "testing"

// 娴嬭瘯 InMemoryVectorStore 娣诲姞鍜屾悳绱?func TestAddAndSearch(t *testing.T) {
	s := NewInMemoryVectorStore()

	docs := []Document{
		NewDocument("Rust is a systems programming language"),
		NewDocument("Python is a scripting language"),
		NewDocument("JavaScript is for web"),
	}
	embeddings := [][]float32{
		{1.0, 0.0, 0.0},
		{0.0, 1.0, 0.0},
		{0.0, 0.0, 1.0},
	}

	ids, err := s.AddDocuments(docs, embeddings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 3 {
		t.Fatalf("expected 3 ids, got %d", len(ids))
	}
	if s.Count() != 3 {
		t.Errorf("expected count 3, got %d", s.Count())
	}

	results, err := s.SimilaritySearch([]float32{0.9, 0.1, 0.0}, 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if !contains(results[0].Document.Content, "Rust") {
		t.Errorf("expected Rust first, got %s", results[0].Document.Content)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// 娴嬭瘯 GetDocument 鍜?DeleteDocument
func TestGetAndDelete(t *testing.T) {
	s := NewInMemoryVectorStore()

	doc := NewDocument("test").WithID("id-1")
	s.AddDocuments([]Document{doc}, [][]float32{{1.0, 0.0, 0.0}})

	got, err := s.GetDocument("id-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Content != "test" {
		t.Errorf("expected test, got %s", got.Content)
	}

	s.DeleteDocument("id-1")
	if s.Count() != 0 {
		t.Errorf("expected 0 after delete, got %d", s.Count())
	}

	_, err = s.GetDocument("id-1")
	if err == nil {
		t.Error("expected error for deleted document")
	}
}

// 娴嬭瘯 Clear
func TestClear(t *testing.T) {
	s := NewInMemoryVectorStore()
	s.AddDocuments(
		[]Document{NewDocument("a"), NewDocument("b")},
		[][]float32{{1.0}, {0.0}},
	)
	s.Clear()
	if s.Count() != 0 {
		t.Errorf("expected 0 after clear, got %d", s.Count())
	}
}

// 娴嬭瘯 Document 鍒涘缓
func TestDocumentCreation(t *testing.T) {
	doc := NewDocument("content").WithID("d1").WithMetadata("key", "val")
	if doc.Content != "content" || doc.ID != "d1" || doc.Metadata["key"] != "val" {
		t.Errorf("unexpected document: %+v", doc)
	}
}
