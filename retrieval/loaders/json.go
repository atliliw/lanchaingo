package loaders

import (
	"encoding/json"
	"fmt"
	"os"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// JSONLoader loads JSON files as documents.
func JSONLoader(filePath string) ([]vs.Document, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, NewLoaderError("read file failed", err)
	}

	var jsonData any
	if err := json.Unmarshal(data, &jsonData); err != nil {
		return nil, NewLoaderError("parse JSON failed", err)
	}

	switch v := jsonData.(type) {
	case []any:
		var docs []vs.Document
		for i, item := range v {
			content := fmt.Sprintf("%v", item)
			doc := vs.NewDocument(content)
			doc.Metadata = map[string]string{"index": fmt.Sprintf("%d", i)}
			docs = append(docs, doc)
		}
		return docs, nil
	default:
		return []vs.Document{vs.NewDocument(string(data))}, nil
	}
}
