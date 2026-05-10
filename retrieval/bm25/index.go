package bm25

import (
	"sync"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// BM25Index stores documents and builds term frequency indices for BM25 scoring.
type BM25Index struct {
	mu           sync.RWMutex
	documents    []vs.Document
	termDocFreqs map[string]map[int]int
	docTermFreqs []map[string]int
	docLengths   []int
	avgdl        float64
	nDocs        int
	idfCache     map[string]float64
	params       BM25Params
}

func NewBM25Index() *BM25Index {
	return NewBM25IndexWithParams(DefaultBM25Params())
}

func NewBM25IndexWithParams(params BM25Params) *BM25Index {
	return &BM25Index{
		termDocFreqs: make(map[string]map[int]int),
		docTermFreqs: make([]map[string]int, 0),
		docLengths:   make([]int, 0),
		idfCache:     make(map[string]float64),
		params:       params,
	}
}

func (idx *BM25Index) AddDocument(doc vs.Document, terms []string) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	docID := idx.nDocs

	termFreq := make(map[string]int)
	for _, term := range terms {
		termFreq[term]++
	}

	for term, freq := range termFreq {
		if idx.termDocFreqs[term] == nil {
			idx.termDocFreqs[term] = make(map[int]int)
		}
		idx.termDocFreqs[term][docID] = freq
	}

	idx.documents = append(idx.documents, doc)
	idx.docTermFreqs = append(idx.docTermFreqs, termFreq)
	idx.docLengths = append(idx.docLengths, len(terms))
	idx.nDocs++

	totalLength := 0
	for _, l := range idx.docLengths {
		totalLength += l
	}
	idx.avgdl = float64(totalLength) / float64(idx.nDocs)

	idx.idfCache = make(map[string]float64)
}

func (idx *BM25Index) AddDocuments(docs []vs.Document, termsList [][]string) {
	for i, doc := range docs {
		if i < len(termsList) {
			idx.AddDocument(doc, termsList[i])
		}
	}
}

func (idx *BM25Index) ComputeIDFForTerm(term string) float64 {
	idx.mu.RLock()
	if idf, ok := idx.idfCache[term]; ok {
		idx.mu.RUnlock()
		return idf
	}
	idx.mu.RUnlock()

	idx.mu.Lock()
	defer idx.mu.Unlock()

	if idf, ok := idx.idfCache[term]; ok {
		return idf
	}

	n := 0
	if docMap, ok := idx.termDocFreqs[term]; ok {
		n = len(docMap)
	}

	idf := ComputeIDF(n, idx.nDocs)
	idx.idfCache[term] = idf
	return idf
}

func (idx *BM25Index) ComputeIDFForTerms(terms []string) map[string]float64 {
	result := make(map[string]float64, len(terms))
	for _, term := range terms {
		result[term] = idx.ComputeIDFForTerm(term)
	}
	return result
}

func (idx *BM25Index) NDocs() int                    { idx.mu.RLock(); defer idx.mu.RUnlock(); return idx.nDocs }
func (idx *BM25Index) Avgdl() float64                 { idx.mu.RLock(); defer idx.mu.RUnlock(); return idx.avgdl }
func (idx *BM25Index) Params() BM25Params             { return idx.params }
func (idx *BM25Index) GetDocument(i int) vs.Document  { idx.mu.RLock(); defer idx.mu.RUnlock(); return idx.documents[i] }

func (idx *BM25Index) GetDocTermFreq(i int) map[string]int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if i < len(idx.docTermFreqs) {
		return idx.docTermFreqs[i]
	}
	return nil
}

func (idx *BM25Index) GetDocLength(i int) int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	if i < len(idx.docLengths) {
		return idx.docLengths[i]
	}
	return 0
}

func (idx *BM25Index) Clear() {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.documents = nil
	idx.termDocFreqs = make(map[string]map[int]int)
	idx.docTermFreqs = nil
	idx.docLengths = nil
	idx.avgdl = 0
	idx.nDocs = 0
	idx.idfCache = make(map[string]float64)
}
