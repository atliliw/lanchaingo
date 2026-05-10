package vector_stores

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
)

type SQLiteConfig struct{ Path string }

type SQLiteVectorStore struct{ db *sql.DB }

func NewSQLiteVectorStore(cfg SQLiteConfig) (*SQLiteVectorStore, error) {
	db, err := sql.Open("sqlite", cfg.Path)
	if err != nil { return nil, fmt.Errorf("sqlite: %w", err) }
	db.Exec("CREATE TABLE IF NOT EXISTS documents (id TEXT, content TEXT, metadata TEXT)")
	return &SQLiteVectorStore{db: db}, nil
}
func (s *SQLiteVectorStore) AddDocuments(docs []Document, embs [][]float32) ([]string, error) {
	return nil, fmt.Errorf("sqlite: vector search requires extension, use as doc store only")
}
func (s *SQLiteVectorStore) SimilaritySearch(q []float32, k int) ([]SearchResult, error) {
	return nil, fmt.Errorf("sqlite: vector search not implemented")
}
func (s *SQLiteVectorStore) GetDocument(id string) (*Document, error) { return nil, fmt.Errorf("sqlite: n/a") }
func (s *SQLiteVectorStore) DeleteDocument(id string) error { return nil }
func (s *SQLiteVectorStore) Count() int {
	var n int; s.db.QueryRow("SELECT COUNT(*) FROM documents").Scan(&n); return n
}
func (s *SQLiteVectorStore) Clear() error {
	_, err := s.db.Exec("DELETE FROM documents"); return err
}
