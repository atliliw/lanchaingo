package memory

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoMemoryConfig struct {
	URI, Database, Collection string
}

type MongoPersistentMemory struct {
	ConversationBufferMemory
	coll    *mongo.Collection
	session string
}

func NewMongoPersistentMemory(ctx context.Context, cfg MongoMemoryConfig) (*MongoPersistentMemory, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, err
	}
	coll := client.Database(cfg.Database).Collection(cfg.Collection)
	return &MongoPersistentMemory{
		ConversationBufferMemory: *NewConversationBufferMemory(),
		coll:                     coll,
	}, nil
}

func (m *MongoPersistentMemory) SetSessionID(id string) { m.session = id }

func (m *MongoPersistentMemory) LoadFromStore(sessionID string) error {
	var result struct {
		Messages []struct{ Role, Content string } `bson:"messages"`
	}
	err := m.coll.FindOne(context.Background(), bson.M{"_id": sessionID}).Decode(&result)
	if err != nil {
		return err
	}
	m.ChatMemory.Clear()
	for _, msg := range result.Messages {
		if msg.Role == "human" {
			m.ChatMemory.AddUserMessage(msg.Content)
		}
		if msg.Role == "ai" {
			m.ChatMemory.AddAIMessage(msg.Content)
		}
	}
	return nil
}

func (m *MongoPersistentMemory) SaveToStore(sessionID string) error {
	id := sessionID
	if id == "" {
		id = m.session
	}
	if id == "" {
		return nil
	}

	msgs := m.ChatMemory.Messages()
	var bsonMsgs []bson.M
	for _, msg := range msgs {
		role := "ai"
		if msg.MessageType == 1 {
			role = "human"
		}
		bsonMsgs = append(bsonMsgs, bson.M{"role": role, "content": msg.Content})
	}

	_, err := m.coll.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"messages": bsonMsgs}},
		options.Update().SetUpsert(true),
	)
	return err
}
