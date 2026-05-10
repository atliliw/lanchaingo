package vector_stores

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct{ URI, Database, Collection string }

type MongoVectorStore struct {
	coll *mongo.Collection
}

func NewMongoVectorStore(ctx context.Context, cfg MongoConfig) (*MongoVectorStore, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil { return nil, fmt.Errorf("mongo: %w", err) }
	return &MongoVectorStore{coll: client.Database(cfg.Database).Collection(cfg.Collection)}, nil
}
func (s *MongoVectorStore) AddDocuments(docs []Document, embs [][]float32) ([]string, error) {
	return nil, fmt.Errorf("mongo: AddDocuments not implemented")
}
func (s *MongoVectorStore) SimilaritySearch(q []float32, k int) ([]SearchResult, error) {
	return nil, fmt.Errorf("mongo: SimilaritySearch not implemented")
}
func (s *MongoVectorStore) GetDocument(id string) (*Document, error) { return nil, fmt.Errorf("mongo: n/a") }
func (s *MongoVectorStore) DeleteDocument(id string) error { return nil }
func (s *MongoVectorStore) Count() int { return 0 }
func (s *MongoVectorStore) Clear() error { return nil }
