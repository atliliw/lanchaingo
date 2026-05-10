package loaders

import (
	"os"
	"strings"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// TextLoader loads plain text files.
func TextLoader(filePath string) ([]vs.Document, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, NewLoaderError("read file failed", err)
	}
	return []vs.Document{vs.NewDocument(string(data))}, nil
}

// TextLoaderFromString loads text from a string.
func TextLoaderFromString(content string) []vs.Document {
	return []vs.Document{vs.NewDocument(strings.TrimSpace(content))}
}
