package bm25

import (
	"strings"
	"unicode"
)

// Tokenizer splits text into tokens, supporting English and Chinese.
type Tokenizer struct {
	keepStopwords bool
	stopwords     map[string]bool
}

func NewTokenizer() *Tokenizer {
	t := &Tokenizer{stopwords: make(map[string]bool)}
	t.initStopwords()
	return t
}

func NewTokenizerWithStopwords() *Tokenizer {
	return &Tokenizer{keepStopwords: true}
}

func (t *Tokenizer) initStopwords() {
	en := []string{"a", "an", "the", "is", "are", "was", "were", "be", "been",
		"have", "has", "had", "do", "does", "did", "will", "would", "could",
		"should", "may", "might", "shall", "can", "need", "dare", "ought",
		"used", "to", "of", "in", "for", "on", "with", "at", "by", "from",
		"as", "into", "through", "during", "before", "after", "above",
		"below", "between", "out", "off", "over", "under", "again",
		"further", "then", "once", "here", "there", "when", "where", "why",
		"how", "all", "each", "every", "both", "few", "more", "most",
		"other", "some", "such", "no", "nor", "not", "only", "own", "same",
		"so", "than", "too", "very", "just", "because", "as", "until",
		"while", "about", "between", "and", "but", "or", "if", "while",
		"that", "this", "it", "its", "you", "your", "he", "him", "his",
		"she", "her", "hers", "we", "us", "our", "they", "them", "their",
		"i", "me", "my", "myself", "am"}
	for _, w := range en {
		t.stopwords[w] = true
	}
}

func (t *Tokenizer) isStopword(token string) bool {
	if t.keepStopwords {
		return false
	}
	return t.stopwords[token]
}

// Tokenize splits text into tokens.
func (t *Tokenizer) Tokenize(text string) []string {
	var tokens []string
	var current []rune

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			current = append(current, r)
		} else {
			if len(current) > 0 {
				token := strings.ToLower(string(current))
				if !t.isStopword(token) {
					tokens = append(tokens, token)
				}
				current = current[:0]
			}
			if isCJK(r) {
				token := string(r)
				if !t.isStopword(token) {
					tokens = append(tokens, token)
				}
			}
		}
	}
	if len(current) > 0 {
		token := strings.ToLower(string(current))
		if !t.isStopword(token) {
			tokens = append(tokens, token)
		}
	}

	return tokens
}

func isCJK(r rune) bool {
	return unicode.In(r, unicode.Han, unicode.Hangul, unicode.Hiragana, unicode.Katakana)
}
