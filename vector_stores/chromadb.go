package vector_stores

// ChromaDBVectorStore integrates with ChromaDB.
type ChromaDBVectorStore struct{}

func NewChromaDBVectorStore(host string, port int) *ChromaDBVectorStore {
	return &ChromaDBVectorStore{}
}

func (s *ChromaDBVectorStore) AddDocuments(documents []Document, embeddings [][]float32) ([]string, error) {
	return nil, NewVectorStoreError(ErrConnection, "ChromaDB requires github.com/amikos-tech/chroma-go", nil)
}

func (s *ChromaDBVectorStore) SimilaritySearch(queryEmbedding []float32, k int) ([]SearchResult, error) {
	return nil, NewVectorStoreError(ErrConnection, "ChromaDB not enabled", nil)
}

func (s *ChromaDBVectorStore) GetDocument(id string) (*Document, error) {
	return nil, NewVectorStoreError(ErrConnection, "ChromaDB not enabled", nil)
}

func (s *ChromaDBVectorStore) DeleteDocument(id string) error {
	return NewVectorStoreError(ErrConnection, "ChromaDB not enabled", nil)
}

func (s *ChromaDBVectorStore) Count() int   { return 0 }
func (s *ChromaDBVectorStore) Clear() error { return nil }
