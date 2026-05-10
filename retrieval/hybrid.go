package retrieval

import (
	"sort"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

const RRFK = 60

// RetrievedDocument is a document with its retrieval source.
type RetrievedDocument struct {
	Document vs.Document
	Score    float64
	Source   RetrievalSource
}

// RetrievalSource indicates where a document came from.
type RetrievalSource int

const (
	SourceBM25   RetrievalSource = iota
	SourceVector
	SourceHybrid
)

// ReciprocalRankFusion combines BM25 and vector search results using RRF.
func ReciprocalRankFusion(bm25Results, vectorResults []vs.Document, k int) []RetrievedDocument {
	rrfScores := make(map[string]*struct {
		score float64
		doc   vs.Document
	})

	for rank, doc := range bm25Results {
		id := doc.ID
		if id == "" {
			id = "bm25"
		}
		contrib := 1.0 / float64(k+rank+1)
		if entry, ok := rrfScores[id]; ok {
			entry.score += contrib
		} else {
			rrfScores[id] = &struct {
				score float64
				doc   vs.Document
			}{score: contrib, doc: doc}
		}
	}

	for rank, doc := range vectorResults {
		id := doc.ID
		if id == "" {
			id = "vector"
		}
		contrib := 1.0 / float64(k+rank+1)
		if entry, ok := rrfScores[id]; ok {
			entry.score += contrib
		} else {
			rrfScores[id] = &struct {
				score float64
				doc   vs.Document
			}{score: contrib, doc: doc}
		}
	}

	results := make([]RetrievedDocument, 0, len(rrfScores))
	for _, entry := range rrfScores {
		results = append(results, RetrievedDocument{
			Document: entry.doc,
			Score:    entry.score,
			Source:   SourceHybrid,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}


