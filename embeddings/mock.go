package embeddings

// MockEmbeddings generates deterministic pseudo-random embeddings for testing.
type MockEmbeddings struct {
	dimension int
}

func NewMockEmbeddings(dimension int) *MockEmbeddings {
	if dimension <= 0 {
		dimension = 1536
	}
	return &MockEmbeddings{dimension: dimension}
}

func hashText(text string) uint64 {
	var hash uint64
	for i, c := range text {
		hash = hash + uint64(c)*uint64(i+1)
	}
	return hash
}

func (m *MockEmbeddings) embed(text string) []float32 {
	hash := hashText(text)
	vec := make([]float32, m.dimension)
	for i := range vec {
		vec[i] = float32(float64(hash+uint64(i))/1000.0 - 0.5)
	}
	var norm float64
	for _, v := range vec {
		norm += float64(v) * float64(v)
	}
	norm = mathSqrt(norm)
	if norm > 0 {
		for i := range vec {
			vec[i] = float32(float64(vec[i]) / norm)
		}
	}
	return vec
}

func mathSqrt(f float64) float64 {
	if f <= 0 {
		return 0
	}
	// Newton's method
	x := f
	for i := 0; i < 10; i++ {
		x = (x + f/x) / 2
	}
	return x
}

func (m *MockEmbeddings) EmbedQuery(text string) ([]float32, error) {
	if text == "" {
		return nil, NewEmbeddingError(ErrEmptyInput, "empty input", nil)
	}
	return m.embed(text), nil
}

func (m *MockEmbeddings) EmbedDocuments(texts []string) ([][]float32, error) {
	result := make([][]float32, len(texts))
	for i, t := range texts {
		emb, err := m.EmbedQuery(t)
		if err != nil {
			return nil, err
		}
		result[i] = emb
	}
	return result, nil
}

func (m *MockEmbeddings) Dimension() int      { return m.dimension }
func (m *MockEmbeddings) ModelName() string   { return "mock-embeddings" }
