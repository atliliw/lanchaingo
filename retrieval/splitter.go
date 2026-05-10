package retrieval

import (
	"strconv"
	"strings"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// TextSplitter defines how to split text into chunks.
type TextSplitter interface {
	SplitText(text string) []string
	SplitDocument(doc vs.Document) []vs.Document
}

// RecursiveCharacterSplitter splits text recursively by separators.
type RecursiveCharacterSplitter struct {
	ChunkSize    int
	ChunkOverlap int
	Separators   []string
}

func NewRecursiveCharacterSplitter(chunkSize, chunkOverlap int) *RecursiveCharacterSplitter {
	return &RecursiveCharacterSplitter{
		ChunkSize:    chunkSize,
		ChunkOverlap: chunkOverlap,
		Separators:   []string{"\n\n", "\n", ".", " ", ""},
	}
}

func (s *RecursiveCharacterSplitter) SplitText(text string) []string {
	return s.splitText(text, s.Separators)
}

func (s *RecursiveCharacterSplitter) splitText(text string, separators []string) []string {
	var chunks []string
	if len(separators) == 0 {
		return []string{text}
	}

	separator := separators[0]
	restSeparators := separators[1:]

	if separator == "" {
		// Split by character
		runes := []rune(text)
		for i := 0; i < len(runes); i += s.ChunkSize - s.ChunkOverlap {
			end := i + s.ChunkSize
			if end > len(runes) {
				end = len(runes)
			}
			chunks = append(chunks, string(runes[i:end]))
		}
		return chunks
	}

	splits := strings.Split(text, separator)

	var merged []string
	for _, split := range splits {
		split = strings.TrimSpace(split)
		if split == "" {
			continue
		}
		merged = append(merged, split)
	}

	if len(merged) == 1 {
		return s.splitText(merged[0], restSeparators)
	}

	for _, chunk := range merged {
		subChunks := s.splitText(chunk, restSeparators)
		chunks = append(chunks, subChunks...)
	}

	return chunks
}

func (s *RecursiveCharacterSplitter) SplitDocument(doc vs.Document) []vs.Document {
	chunks := s.SplitText(doc.Content)
	result := make([]vs.Document, len(chunks))
	for i, chunk := range chunks {
		meta := make(map[string]string)
		for k, v := range doc.Metadata {
			meta[k] = v
		}
		meta["chunk"] = strconv.Itoa(i)
		result[i] = vs.NewDocument(chunk)
		result[i].Metadata = meta
	}
	return result
}
