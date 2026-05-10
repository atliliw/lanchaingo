package bm25

import "testing"

// 娴嬭瘯 ComputeIDF锛氱█鏈夎瘝 IDF 搴旈珮浜庡父瑙佽瘝
func TestComputeIDF(t *testing.T) {
	common := ComputeIDF(100, 100)
	rare := ComputeIDF(1, 100)
	if rare <= common {
		t.Error("rare word should have higher IDF")
	}
	if ComputeIDF(0, 100) != 0 {
		t.Error("zero-doc term should have IDF 0")
	}
}

// 娴嬭瘯 BM25Params 榛樿鍊?func TestBM25Params(t *testing.T) {
	p := DefaultBM25Params()
	if p.K1 != 1.5 || p.B != 0.75 {
		t.Errorf("unexpected defaults: %+v", p)
	}
}

// 娴嬭瘯 BM25Score 鍩虹璁＄畻
func TestBM25Score(t *testing.T) {
	params := DefaultBM25Params()
	queryTerms := []string{"rust", "programming"}
	docFreqs := map[string]int{"rust": 2, "programming": 1}
	idfValues := map[string]float64{"rust": 2.0, "programming": 1.5}

	score := BM25Score(queryTerms, docFreqs, 10, 15.0, idfValues, params)
	if score <= 0 {
		t.Error("expected positive score")
	}
}

// 娴嬭瘯 BM25Score 楂樿瘝棰戞枃妗ｅ緱鍒嗘洿楂?func TestBM25ScoreHighTF(t *testing.T) {
	params := DefaultBM25Params()
	query := []string{"rust"}
	idf := map[string]float64{"rust": 2.0}

	low := BM25Score(query, map[string]int{"rust": 1}, 10, 15.0, idf, params)
	high := BM25Score(query, map[string]int{"rust": 5}, 10, 15.0, idf, params)

	if high <= low {
		t.Error("high TF should score higher")
	}
}

// 娴嬭瘯 Tokenizer 鑻辨枃鍒嗚瘝
func TestTokenizerEnglish(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("Rust programming language")
	if len(tokens) == 0 {
		t.Fatal("expected tokens")
	}
	expected := []string{"rust", "programming", "language"}
	for i, e := range expected {
		if tokens[i] != e {
			t.Errorf("expected %s at %d, got %s", e, i, tokens[i])
		}
	}
}

// 娴嬭瘯 Tokenizer 杩囨护鍋滅敤璇?func TestTokenizerStopwords(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("the rust programming")
	for _, tok := range tokens {
		if tok == "the" {
			t.Error("should filter 'the'")
		}
	}
}

// 娴嬭瘯 Tokenizer 涓枃鍒嗚瘝
func TestTokenizerChinese(t *testing.T) {
	tok := NewTokenizer()
	tokens := tok.Tokenize("Rust鏄竴闂ㄧ紪绋嬭瑷€")
	if len(tokens) == 0 {
		t.Fatal("expected Chinese tokens")
	}
}

// 娴嬭瘯 BM25Index 鍩虹鍔熻兘
func TestBM25IndexBasic(t *testing.T) {
	idx := NewBM25Index()
	idx.AddDocument(newDoc("Rust programming"), []string{"rust", "programming"})
	if idx.NDocs() != 1 {
		t.Errorf("expected 1 doc, got %d", idx.NDocs())
	}
}

// 娴嬭瘯 BM25Index IDF
func TestBM25IndexIDF(t *testing.T) {
	idx := NewBM25Index()
	idx.AddDocument(newDoc("Rust programming"), []string{"rust", "programming", "language"})
	idx.AddDocument(newDoc("Python scripting"), []string{"python", "scripting", "language"})

	rustIDF := idx.ComputeIDFForTerm("rust")
	langIDF := idx.ComputeIDFForTerm("language")
	if rustIDF <= langIDF {
		t.Error("rare term should have higher IDF")
	}
}

// 娴嬭瘯 BM25Retriever 鎼滅储
func TestBM25RetrieverSearch(t *testing.T) {
	r := NewBM25Retriever()
	r.AddDocument(newDoc("Rust is a systems programming language"))
	r.AddDocument(newDoc("Python is a scripting language"))
	r.AddDocument(newDoc("JavaScript for web development"))

	if r.Len() != 3 {
		t.Fatalf("expected 3 docs, got %d", r.Len())
	}

	results := r.Search("programming language", 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

// 娴嬭瘯 BM25Retriever 绌烘悳绱?func TestBM25RetrieverEmpty(t *testing.T) {
	r := NewBM25Retriever()
	results := r.Search("test", 5)
	if len(results) != 0 {
		t.Error("expected empty results")
	}
}
