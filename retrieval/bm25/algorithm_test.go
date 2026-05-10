package bm25

import "testing"

func TestComputeIDF(t *testing.T) {
	if ComputeIDF(100, 100) >= ComputeIDF(1, 100) { t.Error("rare word should have higher IDF") }
	if ComputeIDF(0, 100) != 0 { t.Error("zero-doc term should have 0 IDF") }
}
func TestBM25Params(t *testing.T) {
	p := DefaultBM25Params()
	if p.K1 != 1.5 || p.B != 0.75 { t.Errorf("unexpected defaults: %+v", p) }
}
func TestBM25Score(t *testing.T) {
	p := DefaultBM25Params()
	score := BM25Score([]string{"rust"}, map[string]int{"rust":2}, 10, 15.0, map[string]float64{"rust":2.0}, p)
	if score <= 0 { t.Error("expected positive score") }
}
func TestTokenizer(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("Rust programming")
	if len(tokens) == 0 { t.Fatal("expected tokens") }
}
func TestTokenizerStopwords(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("the rust")
	for _, tok := range tokens { if tok == "the" { t.Error("should filter 'the'") } }
}
func TestBM25IndexBasic(t *testing.T) {
	idx := NewBM25Index()
	idx.AddDocument(newDoc("test"), []string{"test"})
	if idx.NDocs() != 1 { t.Errorf("expected 1") }
}
func TestBM25IndexIDF(t *testing.T) {
	idx := NewBM25Index()
	idx.AddDocument(newDoc("Rust programming"), []string{"rust","programming","language"})
	idx.AddDocument(newDoc("Scripting scripting"), []string{"Scripting","scripting","language"})
	if idx.ComputeIDFForTerm("rust") <= idx.ComputeIDFForTerm("language") { t.Error("rare should have higher IDF") }
}
func TestBM25Retriever(t *testing.T) {
	r := NewBM25Retriever()
	r.AddDocument(newDoc("Rust programming language"))
	r.AddDocument(newDoc("Scripting scripting"))
	if r.Len() != 2 { t.Errorf("expected 2") }
	results := r.Search("programming", 1)
	if len(results) == 0 { t.Error("expected results") }
}
func TestBM25RetrieverEmpty(t *testing.T) {
	r := NewBM25Retriever()
	results := r.Search("test", 5)
	if len(results) != 0 { t.Error("expected empty") }
}
