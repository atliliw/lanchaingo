package embeddings

import "testing"

// 娴嬭瘯 CosineSimilarity 鐩稿悓鍚戦噺
func TestCosineSimilaritySame(t *testing.T) {
	a := []float32{1.0, 0.0, 0.0}
	b := []float32{1.0, 0.0, 0.0}
	sim := CosineSimilarity(a, b)
	if (sim - 1.0) > 0.0001 {
		t.Errorf("expected ~1.0, got %f", sim)
	}
}

// 娴嬭瘯 CosineSimilarity 姝ｄ氦鍚戦噺
func TestCosineSimilarityOrthogonal(t *testing.T) {
	a := []float32{1.0, 0.0, 0.0}
	b := []float32{0.0, 1.0, 0.0}
	sim := CosineSimilarity(a, b)
	if (sim - 0.0) > 0.0001 {
		t.Errorf("expected ~0.0, got %f", sim)
	}
}

// 娴嬭瘯 CosineSimilarity 涓嶅悓闀垮害杩斿洖 0
func TestCosineSimilarityDifferentLengths(t *testing.T) {
	a := []float32{1.0, 0.0}
	b := []float32{1.0, 0.0, 0.0}
	if CosineSimilarity(a, b) != 0.0 {
		t.Error("expected 0 for different lengths")
	}
}
