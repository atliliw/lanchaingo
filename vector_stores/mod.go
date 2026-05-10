package vector_stores

// Document represents a document in the vector store.
type Document struct {
	ID       string
	Content  string
	Metadata map[string]string
}

func NewDocument(content string) Document {
	return Document{
		Content:  content,
		Metadata: make(map[string]string),
	}
}

func (d Document) WithID(id string) Document {
	d.ID = id
	return d
}

func (d Document) WithMetadata(key, value string) Document {
	if d.Metadata == nil {
		d.Metadata = make(map[string]string)
	}
	d.Metadata[key] = value
	return d
}

// SearchResult is a document with a similarity score.
type SearchResult struct {
	Document Document
	Score    float64
}

// VectorStore defines the interface for storing and searching document vectors.
type VectorStore interface {
	AddDocuments(documents []Document, embeddings [][]float32) ([]string, error)
	SimilaritySearch(queryEmbedding []float32, k int) ([]SearchResult, error)
	GetDocument(id string) (*Document, error)
	DeleteDocument(id string) error
	Count() int
	Clear() error
}
