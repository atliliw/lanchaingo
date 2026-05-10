//go:build qdrant
package vector_stores

// QdrantVectorStore integrates with Qdrant vector database.
// Build with: go build -tags qdrant
type QdrantVectorStore struct{}

func NewQdrantVectorStore(url, apiKey string) *QdrantVectorStore {
	return &QdrantVectorStore{}
}

func (s *QdrantVectorStore) AddDocuments(documents []Document, embeddings [][]float32) ([]string, error) {
	return nil, NewVectorStoreError(ErrConnection, "Qdrant not enabled (build tag: qdrant)", nil)
}

func (s *QdrantVectorStore) SimilaritySearch(queryEmbedding []float32, k int) ([]SearchResult, error) {
	return nil, NewVectorStoreError(ErrConnection, "Qdrant not enabled", nil)
}

func (s *QdrantVectorStore) GetDocument(id string) (*Document, error) {
	return nil, NewVectorStoreError(ErrConnection, "Qdrant not enabled", nil)
}

func (s *QdrantVectorStore) DeleteDocument(id string) error {
	return NewVectorStoreError(ErrConnection, "Qdrant not enabled", nil)
}

func (s *QdrantVectorStore) Count() int     { return 0 }
func (s *QdrantVectorStore) Clear() error   { return nil }
