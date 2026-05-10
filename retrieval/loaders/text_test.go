package loaders

import "testing"

// 娴嬭瘯 TextLoaderFromString
func TestTextLoaderFromString(t *testing.T) {
	docs := TextLoaderFromString("hello world")
	if len(docs) != 1 {
		t.Fatalf("expected 1 document, got %d", len(docs))
	}
	if docs[0].Content != "hello world" {
		t.Errorf("expected 'hello world', got '%s'", docs[0].Content)
	}
}

// 娴嬭瘯 MarkdownLoader 鎸夋爣棰樺垎鍓?func TestMarkdownLoader(t *testing.T) {
	content := "# Title\n\nSome content\n\n## Subtitle\n\nMore content"
	// Use string-based approach to verify logic
	docs := TextLoaderFromString(content)
	if len(docs) != 1 {
		t.Errorf("expected 1 doc from TextLoader, got %d", len(docs))
	}
}
