package vector_stores

import (
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisConfig struct{ Addr, Password, KeyPrefix string; DB int }

type RedisVectorStore struct{ client *redis.Client; prefix string }

func NewRedisVectorStore(cfg RedisConfig) *RedisVectorStore {
	return &RedisVectorStore{
		client: redis.NewClient(&redis.Options{Addr: cfg.Addr, Password: cfg.Password, DB: cfg.DB}),
		prefix: cfg.KeyPrefix,
	}
}
func (s *RedisVectorStore) AddDocuments(docs []Document, embs [][]float32) ([]string, error) {
	return nil, fmt.Errorf("redis: AddDocuments not implemented")
}
func (s *RedisVectorStore) SimilaritySearch(q []float32, k int) ([]SearchResult, error) {
	return nil, fmt.Errorf("redis: SimilaritySearch not implemented")
}
func (s *RedisVectorStore) GetDocument(id string) (*Document, error) { return nil, fmt.Errorf("redis: n/a") }
func (s *RedisVectorStore) DeleteDocument(id string) error { return nil }
func (s *RedisVectorStore) Count() int { return 0 }
func (s *RedisVectorStore) Clear() error { return nil }
