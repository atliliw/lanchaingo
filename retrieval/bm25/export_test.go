package bm25

import vs "github.com/atliliw/lanchaingo/vector_stores"

func newDoc(content string) vs.Document { return vs.NewDocument(content) }
