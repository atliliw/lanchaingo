package memory

import (
	"github.com/atliliw/lanchaingo/schema/messages"
)

// PersistentMemory extends BaseMemory with persistence capabilities.
type PersistentMemory interface {
	BaseMemory
	LoadFromStore(sessionID string) error
	SaveToStore(sessionID string) error
	DeleteSession(sessionID string) error
	SessionExists(sessionID string) (bool, error)
	CurrentSessionID() string
	SetSessionID(sessionID string)
}

// PersistenceConfig holds configuration for persistent memory.
type PersistenceConfig struct {
	AutoSave     bool
	AutoLoad     bool
	MaxMessages  int
	TokenLimit   int
}

func DefaultPersistenceConfig() PersistenceConfig {
	return PersistenceConfig{
		AutoSave:    true,
		AutoLoad:    true,
		MaxMessages: 100,
		TokenLimit:  4000,
	}
}

func (c PersistenceConfig) WithAutoSave(v bool) PersistenceConfig {
	c.AutoSave = v
	return c
}

func (c PersistenceConfig) WithAutoLoad(v bool) PersistenceConfig {
	c.AutoLoad = v
	return c
}

func (c PersistenceConfig) WithMaxMessages(v int) PersistenceConfig {
	c.MaxMessages = v
	return c
}

// MemoryData represents serializable memory state.
type MemoryData struct {
	SessionID string
	Messages  []messages.Message
	Summary   string
	Metadata  map[string]string
}

func NewMemoryData(sessionID string) *MemoryData {
	return &MemoryData{
		SessionID: sessionID,
		Messages:  make([]messages.Message, 0),
		Metadata:  make(map[string]string),
	}
}

func (d *MemoryData) AddMessage(msg messages.Message) {
	d.Messages = append(d.Messages, msg)
}
