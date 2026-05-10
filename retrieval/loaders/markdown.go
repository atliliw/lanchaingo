package loaders

import (
	"os"
	"strings"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// MarkdownLoader loads markdown files, splitting by headings.
func MarkdownLoader(filePath string) ([]vs.Document, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, NewLoaderError("read file failed", err)
	}

	content := string(data)
	var docs []vs.Document

	lines := strings.Split(content, "\n")
	var currentSection strings.Builder
	var currentHeading string

	for _, line := range lines {
		if strings.HasPrefix(line, "#") {
			if currentSection.Len() > 0 {
				doc := vs.NewDocument(strings.TrimSpace(currentSection.String()))
				if currentHeading != "" {
					doc.Metadata = map[string]string{"heading": currentHeading}
				}
				docs = append(docs, doc)
				currentSection.Reset()
			}
			currentHeading = strings.TrimLeft(line, "# ")
		} else {
			currentSection.WriteString(line + "\n")
		}
	}

	if currentSection.Len() > 0 {
		doc := vs.NewDocument(strings.TrimSpace(currentSection.String()))
		if currentHeading != "" {
			doc.Metadata = map[string]string{"heading": currentHeading}
		}
		docs = append(docs, doc)
	}

	return docs, nil
}
