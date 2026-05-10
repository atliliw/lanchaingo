package retrieval

import (
	"math"
	"sort"
	"strings"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// Reranker scores documents for reranking.
type Reranker interface {
	Score(query string, documents []vs.Document) ([]float64, error)
}

// RerankingConfig holds reranking parameters.
type RerankingConfig struct {
	TopN                  int
	MinScore              float64
	PreserveOriginalScore bool
}

func DefaultRerankingConfig() RerankingConfig {
	return RerankingConfig{
		TopN:                  5,
		MinScore:              0,
		PreserveOriginalScore: true,
	}
}

// KeywordReranker scores based on keyword overlap with query.
type KeywordReranker struct {
	config RerankingConfig
}

func NewKeywordReranker() *KeywordReranker {
	return &KeywordReranker{config: DefaultRerankingConfig()}
}

func (r *KeywordReranker) Score(query string, documents []vs.Document) ([]float64, error) {
	queryTerms := strings.Fields(strings.ToLower(query))
	scores := make([]float64, len(documents))

	for i, doc := range documents {
		docLower := strings.ToLower(doc.Content)
		var matchCount float64
		for _, term := range queryTerms {
			if strings.Contains(docLower, term) {
				matchCount++
			}
		}
		scores[i] = matchCount / float64(len(queryTerms))
	}
	return scores, nil
}

// BM25Reranker uses BM25-like scoring for reranking.
type BM25Reranker struct {
	config RerankingConfig
}

func NewBM25Reranker() *BM25Reranker {
	return &BM25Reranker{config: DefaultRerankingConfig()}
}

func (r *BM25Reranker) Score(query string, documents []vs.Document) ([]float64, error) {
	queryTerms := strings.Fields(strings.ToLower(query))
	scores := make([]float64, len(documents))

	for i, doc := range documents {
		docLower := strings.ToLower(doc.Content)
		terms := strings.Fields(docLower)
		docLen := len(terms)
		if docLen == 0 {
			continue
		}

		var totalLen float64
		var avgdl float64
		for _, d := range documents {
			totalLen += float64(len(strings.Fields(strings.ToLower(d.Content))))
		}
		avgdl = totalLen / float64(len(documents))

		for _, qt := range queryTerms {
			tf := 0
			for _, t := range terms {
				if t == qt {
					tf++
				}
			}
			if tf == 0 {
				continue
			}
			n := 0
			for _, d := range documents {
				if strings.Contains(strings.ToLower(d.Content), qt) {
					n++
				}
			}
			idf := math.Log(float64(len(documents)-n+1) / float64(n+1))
			k1 := 1.5
			b := 0.75
			scores[i] += idf * (float64(tf) * (k1 + 1)) / (float64(tf) + k1*(1-b+b*float64(docLen)/avgdl))
		}
	}
	return scores, nil
}

// RerankingExecutor applies a reranker to search results.
type RerankingExecutor struct {
	reranker Reranker
	config   RerankingConfig
}

func NewRerankingExecutor(reranker Reranker) *RerankingExecutor {
	return &RerankingExecutor{
		reranker: reranker,
		config:   DefaultRerankingConfig(),
	}
}

func (e *RerankingExecutor) Rerank(query string, results []vs.SearchResult) ([]vs.SearchResult, error) {
	docs := make([]vs.Document, len(results))
	origScores := make([]float64, len(results))
	for i, r := range results {
		docs[i] = r.Document
		origScores[i] = r.Score
	}

	scores, err := e.reranker.Score(query, docs)
	if err != nil {
		return nil, err
	}

	type scoredDoc struct {
		result vs.SearchResult
		score  float64
	}

	scored := make([]scoredDoc, len(results))
	for i := range results {
		s := scores[i]
		if e.config.PreserveOriginalScore {
			s = (s + origScores[i]) / 2
		}
		scored[i] = scoredDoc{result: results[i], score: s}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	final := make([]vs.SearchResult, 0, len(scored))
	for _, s := range scored {
		if s.score >= e.config.MinScore {
			s.result.Score = s.score
			final = append(final, s.result)
		}
	}

	if len(final) > e.config.TopN {
		final = final[:e.config.TopN]
	}

	return final, nil
}
