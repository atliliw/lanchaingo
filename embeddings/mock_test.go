package embeddings

import "testing"

func TestCosineSimilaritySame(t *testing.T) {
	a := []float32{1.0, 0.0, 0.0}
	b := []float32{1.0, 0.0, 0.0}
	s := CosineSimilarity(a, b)
	if (s - 1.0) > 0.0001 { t.Errorf("expected ~1.0, got %f", s) }
}
func TestCosineSimilarityOrthogonal(t *testing.T) {
	a := []float32{1.0, 0.0}
	b := []float32{0.0, 1.0}
	s := CosineSimilarity(a, b)
	if (s - 0.0) > 0.0001 { t.Errorf("expected ~0.0, got %f", s) }
}
func TestMockEmbeddings(t *testing.T) {
	m := NewMockEmbeddings(128)
	e, _ := m.EmbedQuery("hello")
	if len(e) != 128 { t.Errorf("expected 128 dims") }
	e2, _ := m.EmbedQuery("hello")
	for i := range e { if e[i] != e2[i] { t.Errorf("inconsistent at %d", i); break } }
}
func TestMockEmbeddingsEmpty(t *testing.T) {
	m := NewMockEmbeddings(64)
	_, err := m.EmbedQuery("")
	if err == nil { t.Error("expected error") }
}
func TestMockEmbedDocuments(t *testing.T) {
	m := NewMockEmbeddings(32)
	r, _ := m.EmbedDocuments([]string{"a","b"})
	if len(r) != 2 { t.Errorf("expected 2, got %d", len(r)) }
}
func TestOpenAIDefaultConfig(t *testing.T) {
	cfg := DefaultOpenAIEmbeddingsConfig()
	if cfg.Model != "text-embedding-ada-002" { t.Errorf("unexpected model: %s", cfg.Model) }
}
func TestOpenAIDimension(t *testing.T) {
	e := NewOpenAIEmbeddings(DefaultOpenAIEmbeddingsConfig())
	if e.Dimension() != 1536 { t.Errorf("expected 1536") }
}
