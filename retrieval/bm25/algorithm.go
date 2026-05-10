package bm25

import (
	"math"
)

// BM25Params holds BM25 scoring parameters.
type BM25Params struct {
	K1 float64
	B  float64
}

func DefaultBM25Params() BM25Params {
	return BM25Params{K1: 1.5, B: 0.75}
}

func NewBM25Params(k1, b float64) BM25Params {
	return BM25Params{K1: k1, B: b}
}

// ComputeIDF calculates IDF for a term.
// IDF(qi) = log((N - n(qi) + 0.5) / (n(qi) + 0.5) + 1)
func ComputeIDF(n, totalDocs int) float64 {
	if n == 0 || totalDocs == 0 {
		return 0
	}
	numerator := float64(totalDocs-n) + 0.5
	denominator := float64(n) + 0.5
	return math.Log(numerator/denominator + 1.0)
}

// BM25Score calculates BM25 score for a document against query terms.
func BM25Score(
	queryTerms []string,
	docTermFreqs map[string]int,
	docLength int,
	avgdl float64,
	idfValues map[string]float64,
	params BM25Params,
) float64 {
	if avgdl == 0 || docLength == 0 {
		return 0
	}

	var score float64
	for _, term := range queryTerms {
		idf := idfValues[term]
		if idf == 0 {
			continue
		}
		tf := docTermFreqs[term]
		if tf == 0 {
			continue
		}

		dlRatio := float64(docLength) / avgdl
		tfComponent := float64(tf)*(params.K1+1.0) / (float64(tf) + params.K1*(1.0-params.B+params.B*dlRatio))
		score += idf * tfComponent
	}
	return score
}
