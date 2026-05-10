package bm25

import (
	"math"
	"sort"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// AutoMergingConfig controls parent-child document merging.
type AutoMergingConfig struct {
	MergeThreshold  float64
	LeafChunkSize   int
	ParentChunkSize int
	LeavesPerParent int
}

func DefaultAutoMergingConfig() AutoMergingConfig {
	return AutoMergingConfig{
		MergeThreshold:  0.5,
		LeafChunkSize:   400,
		ParentChunkSize: 2000,
		LeavesPerParent: 5,
	}
}

// ChunkedBM25Retriever supports parent-child document BM25 search with auto-merging.
type ChunkedBM25Retriever struct {
	leafIndex    *BM25Index
	parentIndex  *BM25Index
	leafToParent map[int]int // leaf doc id -> parent doc id
	parentToLeaf map[int][]int
	parents      []vs.Document
	leaves       []vs.Document
	autoMergeCfg AutoMergingConfig
}

func NewChunkedBM25Retriever() *ChunkedBM25Retriever {
	return &ChunkedBM25Retriever{
		leafIndex:    NewBM25Index(),
		parentIndex:  NewBM25Index(),
		leafToParent: make(map[int]int),
		parentToLeaf: make(map[int][]int),
		autoMergeCfg: DefaultAutoMergingConfig(),
	}
}

func (r *ChunkedBM25Retriever) AddDocument(parent vs.Document, leaves []vs.Document, leafTerms [][]string) {
	parentID := len(r.parents)
	r.parents = append(r.parents, parent)

	// Tokenize parent
	parentTokenizer := NewTokenizer()
	parentTerms := parentTokenizer.Tokenize(parent.Content)
	r.parentIndex.AddDocument(parent, parentTerms)

	for i, leaf := range leaves {
		leafID := len(r.leaves)
		r.leaves = append(r.leaves, leaf)
		r.leafToParent[leafID] = parentID
		r.parentToLeaf[parentID] = append(r.parentToLeaf[parentID], leafID)

		terms := leafTerms[i]
		r.leafIndex.AddDocument(leaf, terms)
	}
}

func (r *ChunkedBM25Retriever) Search(query string, k int) []vs.SearchResult {
	queryTerms := NewTokenizer().Tokenize(query)
	if len(queryTerms) == 0 {
		return nil
	}

	// Search leaves
	leafIDF := r.leafIndex.ComputeIDFForTerms(queryTerms)
	type scored struct {
		id    int
		score float64
	}
	var leafScores []scored

	for i := 0; i < r.leafIndex.NDocs(); i++ {
		tf := r.leafIndex.GetDocTermFreq(i)
		dl := r.leafIndex.GetDocLength(i)
		avgdl := r.leafIndex.Avgdl()
		params := r.leafIndex.Params()
		score := BM25Score(queryTerms, tf, dl, avgdl, leafIDF, params)
		if score > 0 {
			leafScores = append(leafScores, scored{id: i, score: score})
		}
	}

	sort.Slice(leafScores, func(i, j int) bool {
		return leafScores[i].score > leafScores[j].score
	})

	// Auto-merge: group leaf scores by parent
	type parentScore struct {
		parentID int
		total    float64
		count    int
	}
	parentScores := make(map[int]*parentScore)

	for _, ls := range leafScores {
		pid := r.leafToParent[ls.id]
		if ps, ok := parentScores[pid]; ok {
			ps.total += ls.score
			ps.count++
		} else {
			parentScores[pid] = &parentScore{parentID: pid, total: ls.score, count: 1}
		}
	}

	var merged []parentScore
	for _, ps := range parentScores {
		ratio := float64(ps.count) / float64(r.autoMergeCfg.LeavesPerParent)
		if ratio >= r.autoMergeCfg.MergeThreshold {
			ps.total *= (1 + math.Log2(1+float64(ps.count)))
		}
		merged = append(merged, *ps)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].total > merged[j].total
	})

	if k > len(merged) {
		k = len(merged)
	}
	results := make([]vs.SearchResult, k)
	for i, m := range merged[:k] {
		results[i] = vs.SearchResult{
			Document: r.parents[m.parentID],
			Score:    m.total,
		}
	}
	return results
}

func (r *ChunkedBM25Retriever) Len() int { return len(r.parents) }
func (r *ChunkedBM25Retriever) Parents() []vs.Document { return r.parents }
func (r *ChunkedBM25Retriever) Leaves() []vs.Document { return r.leaves }
