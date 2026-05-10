package vector_stores

import "fmt"

type ChromaDBConfig struct{ Host string; Port int }
type ChromaDBVectorStore struct{ cfg ChromaDBConfig }

func NewChromaDBVectorStore(cfg ChromaDBConfig) *ChromaDBVectorStore {
	return &ChromaDBVectorStore{cfg: cfg}
}
func (s *ChromaDBVectorStore) AddDocuments(docs []Document, embs [][]float32) ([]string, error) {
	return nil, fmt.Errorf("chromadb: requires github.com/amikos-tech/chroma-go")
}
func (s *ChromaDBVectorStore) SimilaritySearch(q []float32, k int) ([]SearchResult, error) {
	return nil, fmt.Errorf("chromadb: not enabled")
}
func (s *ChromaDBVectorStore) GetDocument(id string) (*Document, error) { return nil, fmt.Errorf("n/a") }
func (s *ChromaDBVectorStore) DeleteDocument(id string) error { return nil }
func (s *ChromaDBVectorStore) Count() int { return 0 }
func (s *ChromaDBVectorStore) Clear() error { return nil }
