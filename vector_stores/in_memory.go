package vector_stores

import (
	"math"
	"sort"
	"sync"

	"github.com/google/uuid"
)

type vectorDocument struct {
	document  Document
	embedding []float32
}

// InMemoryVectorStore stores documents and vectors in memory.
type InMemoryVectorStore struct {
	mu    sync.RWMutex
	docs  map[string]*vectorDocument
}

func NewInMemoryVectorStore() *InMemoryVectorStore {
	return &InMemoryVectorStore{
		docs: make(map[string]*vectorDocument),
	}
}

func cosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		ai, bi := float64(a[i]), float64(b[i])
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

func (s *InMemoryVectorStore) AddDocuments(documents []Document, embeddings [][]float32) ([]string, error) {
	if len(documents) != len(embeddings) {
		return nil, NewVectorStoreError(ErrStorage, "documents and embeddings count mismatch", nil)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	ids := make([]string, len(documents))
	for i, doc := range documents {
		id := doc.ID
		if id == "" {
			id = uuid.New().String()
		}
		s.docs[id] = &vectorDocument{
			document:  doc,
			embedding: embeddings[i],
		}
		ids[i] = id
	}
	return ids, nil
}

func (s *InMemoryVectorStore) SimilaritySearch(queryEmbedding []float32, k int) ([]SearchResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []SearchResult
	for _, vd := range s.docs {
		score := cosineSimilarity(queryEmbedding, vd.embedding)
		results = append(results, SearchResult{
			Document: vd.document,
			Score:    score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if k > len(results) {
		k = len(results)
	}
	return results[:k], nil
}

func (s *InMemoryVectorStore) GetDocument(id string) (*Document, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	vd, ok := s.docs[id]
	if !ok {
		return nil, NewVectorStoreError(ErrDocumentNotFound, "document not found: "+id, nil)
	}
	doc := vd.document
	return &doc, nil
}

func (s *InMemoryVectorStore) DeleteDocument(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.docs, id)
	return nil
}

func (s *InMemoryVectorStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.docs)
}

func (s *InMemoryVectorStore) Clear() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.docs = make(map[string]*vectorDocument)
	return nil
}
