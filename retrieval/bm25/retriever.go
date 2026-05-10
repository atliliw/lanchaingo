package bm25

import (
	"sort"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// BM25Retriever performs keyword-based retrieval using BM25 scoring.
type BM25Retriever struct {
	index     *BM25Index
	tokenizer *Tokenizer
}

func NewBM25Retriever() *BM25Retriever {
	return &BM25Retriever{
		index:     NewBM25Index(),
		tokenizer: NewTokenizer(),
	}
}

func NewBM25RetrieverWithParams(k1, b float64) *BM25Retriever {
	return &BM25Retriever{
		index:     NewBM25IndexWithParams(NewBM25Params(k1, b)),
		tokenizer: NewTokenizer(),
	}
}

func (r *BM25Retriever) AddDocument(doc vs.Document) {
	terms := r.tokenizer.Tokenize(doc.Content)
	r.index.AddDocument(doc, terms)
}

func (r *BM25Retriever) AddDocuments(docs []vs.Document) {
	for _, doc := range docs {
		r.AddDocument(doc)
	}
}

func (r *BM25Retriever) Search(query string, k int) []vs.SearchResult {
	if r.index.NDocs() == 0 {
		return nil
	}

	queryTerms := r.tokenizer.Tokenize(query)
	if len(queryTerms) == 0 {
		return nil
	}

	idfValues := r.index.ComputeIDFForTerms(queryTerms)

	type scoredDoc struct {
		id    int
		score float64
	}

	var scored []scoredDoc
	for docID := 0; docID < r.index.NDocs(); docID++ {
		termFreqs := r.index.GetDocTermFreq(docID)
		docLength := r.index.GetDocLength(docID)
		avgdl := r.index.Avgdl()
		params := r.index.Params()

		score := BM25Score(queryTerms, termFreqs, docLength, avgdl, idfValues, params)
		if score > 0 {
			scored = append(scored, scoredDoc{id: docID, score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score > scored[j].score
	})

	if k > len(scored) {
		k = len(scored)
	}

	results := make([]vs.SearchResult, k)
	for i, sd := range scored[:k] {
		doc := r.index.GetDocument(sd.id)
		results[i] = vs.SearchResult{
			Document: doc,
			Score:    sd.score,
		}
	}
	return results
}

func (r *BM25Retriever) Len() int    { return r.index.NDocs() }
func (r *BM25Retriever) Clear()      { r.index.Clear() }
func (r *BM25Retriever) Index() *BM25Index { return r.index }
