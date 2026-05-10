package retrieval

import (
	"testing"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// 娴嬭瘯 RecursiveCharacterSplitter 鍩烘湰鍒嗗壊
func TestSplitterBasic(t *testing.T) {
	s := NewRecursiveCharacterSplitter(50, 0)
	chunks := s.SplitText("Hello world. This is a test.")
	if len(chunks) == 0 {
		t.Error("expected at least 1 chunk")
	}
}

// 娴嬭瘯 SplitDocument 鍏冩暟鎹紶鎾?func TestSplitterDocument(t *testing.T) {
	s := NewRecursiveCharacterSplitter(100, 0)
	doc := vs.NewDocument("Part one.\n\nPart two.\n\nPart three.").
		WithMetadata("source", "test")

	docs := s.SplitDocument(doc)
	if len(docs) < 2 {
		t.Error("expected at least 2 chunks")
	}

	for _, d := range docs {
		if d.Metadata["source"] != "test" {
			t.Errorf("expected source=test in all chunks, got %v", d.Metadata)
		}
	}
}
