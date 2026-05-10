package vector_stores

import (
	"context"
	"fmt"
	pb "github.com/qdrant/go-client/qdrant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QdrantConfig struct {
	URL, Collection string
	VectorSize      int
}
type QdrantVectorStore struct {
	points     pb.PointsClient
	collection string
	conn       *grpc.ClientConn
}

func NewQdrantVectorStore(ctx context.Context, cfg QdrantConfig) (*QdrantVectorStore, error) {
	conn, err := grpc.DialContext(ctx, cfg.URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil { return nil, fmt.Errorf("qdrant: %w", err) }
	return &QdrantVectorStore{
		points: pb.NewPointsClient(conn), collection: cfg.Collection, conn: conn,
	}, nil
}
func (s *QdrantVectorStore) AddDocuments(documents []Document, embeddings [][]float32) ([]string, error) {
	ids := make([]string, len(documents))
	var pts []*pb.PointStruct
	for i, doc := range documents {
		id := fmt.Sprintf("d%d", i)
		ids[i] = id
		pts = append(pts, &pb.PointStruct{
			Id: &pb.PointId{PointIdOptions: &pb.PointId_Uuid{Uuid: id}},
			Vectors: &pb.Vectors{VectorsOptions: &pb.Vectors_Vector{Vector: &pb.Vector{Data: embeddings[i]}}},
			Payload: map[string]*pb.Value{"content": {Kind: &pb.Value_StringValue{StringValue: doc.Content}}},
		})
	}
	_, err := s.points.Upsert(context.Background(), &pb.UpsertPoints{CollectionName: s.collection, Points: pts})
	return ids, err
}
func (s *QdrantVectorStore) SimilaritySearch(query []float32, k int) ([]SearchResult, error) {
	resp, err := s.points.Search(context.Background(), &pb.SearchPoints{
		CollectionName: s.collection, Vector: query, Limit: uint64(k),
		WithPayload: &pb.WithPayloadSelector{SelectorOptions: &pb.WithPayloadSelector_Enable{Enable: true}},
	})
	if err != nil { return nil, err }
	var results []SearchResult
	for _, r := range resp.GetResult() {
		c := ""
		if v, ok := r.GetPayload()["content"]; ok { c = v.GetStringValue() }
		results = append(results, SearchResult{Document: Document{Content: c}, Score: float64(r.GetScore())})
	}
	return results, nil
}
func (s *QdrantVectorStore) GetDocument(id string) (*Document, error) {
	return nil, NewVectorStoreError(ErrDocumentNotFound, "", nil)
}
func (s *QdrantVectorStore) DeleteDocument(id string) error { return nil }
func (s *QdrantVectorStore) Count() int { return 0 }
func (s *QdrantVectorStore) Clear() error {
	_, err := s.points.Delete(context.Background(), &pb.DeletePoints{
		CollectionName: s.collection,
		Points: &pb.PointsSelector{
			PointsSelectorOneOf: &pb.PointsSelector_Filter{Filter: &pb.Filter{}},
		},
	})
	return err
}