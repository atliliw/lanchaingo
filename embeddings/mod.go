package embeddings

import "math"

// Embeddings defines the interface for text embedding models.
type Embeddings interface {
	EmbedQuery(text string) ([]float32, error)
	EmbedDocuments(texts []string) ([][]float32, error)
	Dimension() int
	ModelName() string
}

// CosineSimilarity computes cosine similarity between two vectors.
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}

	var dot, normA, normB float64
	for i := range a {
		ai := float64(a[i])
		bi := float64(b[i])
		dot += ai * bi
		normA += ai * ai
		normB += bi * bi
	}

	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
