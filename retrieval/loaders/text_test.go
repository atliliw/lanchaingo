package loaders

import "testing"

func TestTextLoaderFromString(t *testing.T) {
	docs := TextLoaderFromString("hello")
	if len(docs) != 1 { t.Fatalf("expected 1, got %d", len(docs)) }
	if docs[0].Content != "hello" { t.Errorf("unexpected content") }
}
