package loaders

import (
	"os"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// PDFLoader loads PDF files as text.
func PDFLoader(filePath string) ([]vs.Document, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, NewLoaderError("read file failed", err)
	}
	return []vs.Document{vs.NewDocument(string(data))}, nil
}
