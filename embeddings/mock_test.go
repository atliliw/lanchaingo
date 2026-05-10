package embeddings

import "testing"

// 娴嬭瘯 MockEmbeddings 缁村害
func TestMockEmbeddingsDimension(t *testing.T) {
	m := NewMockEmbeddings(128)
	emb, err := m.EmbedQuery("hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(emb) != 128 {
		t.Errorf("expected 128 dims, got %d", len(emb))
	}
}

// 娴嬭瘯 MockEmbeddings 涓€鑷存€э細鐩稿悓鏂囨湰浜х敓鐩稿悓鍚戦噺
func TestMockEmbeddingsConsistency(t *testing.T) {
	m := NewMockEmbeddings(64)
	e1, _ := m.EmbedQuery("test")
	e2, _ := m.EmbedQuery("test")
	for i := range e1 {
		if e1[i] != e2[i] {
			t.Errorf("expected consistent at index %d: %f vs %f", i, e1[i], e2[i])
			break
		}
	}
}

// 娴嬭瘯 MockEmbeddings 绌鸿緭鍏?func TestMockEmbeddingsEmpty(t *testing.T) {
	m := NewMockEmbeddings(64)
	_, err := m.EmbedQuery("")
	if err == nil {
		t.Error("expected error for empty input")
	}
}

// 娴嬭瘯 MockEmbeddings EmbedDocuments
func TestMockEmbedDocuments(t *testing.T) {
	m := NewMockEmbeddings(32)
	results, err := m.EmbedDocuments([]string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results, got %d", len(results))
	}
}

// 娴嬭瘯 MockEmbeddings ModelName
func TestMockEmbeddingsModelName(t *testing.T) {
	m := NewMockEmbeddings(64)
	if m.ModelName() != "mock-embeddings" {
		t.Errorf("unexpected model name: %s", m.ModelName())
	}
}
