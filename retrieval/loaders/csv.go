package loaders

import (
	"encoding/csv"
	"fmt"
	"os"

	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// CSVLoader loads CSV files, each row becomes a document.
func CSVLoader(filePath string) ([]vs.Document, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, NewLoaderError("open file failed", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, NewLoaderError("read CSV failed", err)
	}

	if len(rows) == 0 {
		return nil, nil
	}

	headers := rows[0]
	var docs []vs.Document

	for _, row := range rows[1:] {
		content := ""
		meta := make(map[string]string)
		for i, val := range row {
			if i < len(headers) {
				if content != "" {
					content += " "
				}
				content += fmt.Sprintf("%s: %s", headers[i], val)
				meta[headers[i]] = val
			}
		}
		doc := vs.NewDocument(content)
		doc.Metadata = meta
		docs = append(docs, doc)
	}

	return docs, nil
}
