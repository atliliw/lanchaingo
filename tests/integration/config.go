// Package integration provides config loading and test utilities
package integration

import (
	"os"
	"strings"
)

// TestConfig stores test configuration loaded from config.toml
type TestConfig struct {
	APIKey      string
	BaseURL     string
	ChatModel   string
	EmbedModel  string

	MongoURI    string
	MongoDB     string
	ParentColl  string
	ChunkColl   string

	SQLitePath  string

	QdrantURL   string
	QdrantColl  string
	VectorSize  int
}

// LoadConfig loads configuration from config.toml
func LoadConfig() *TestConfig {
	cfg := &TestConfig{}

	data, err := os.ReadFile("../../config/config.toml")
	if err != nil {
		return cfg
	}

	content := string(data)
	section := ""

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSuffix(strings.TrimPrefix(line, "["), "]")
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
		if value == "" {
			continue
		}

		switch section {
		case "openai":
			switch key {
			case "api_key":
				cfg.APIKey = value
			case "base_url":
				cfg.BaseURL = value
			case "chat_model":
				cfg.ChatModel = value
			case "embedding_model":
				cfg.EmbedModel = value
			}
		case "embedding":
			if key == "api_key" && cfg.APIKey == "" {
				cfg.APIKey = value
			}
			if key == "base_url" && cfg.BaseURL == "" {
				cfg.BaseURL = value
			}
			if key == "model" && cfg.EmbedModel == "" {
				cfg.EmbedModel = value
			}
		case "mongodb":
			switch key {
			case "uri":
				cfg.MongoURI = value
			case "database":
				cfg.MongoDB = value
			case "parent_collection":
				cfg.ParentColl = value
			case "chunk_collection":
				cfg.ChunkColl = value
			}
		case "sqlite":
			if key == "db_path" {
				cfg.SQLitePath = value
			}
		case "qdrant":
			switch key {
			case "url":
				cfg.QdrantURL = value
			case "collection_name":
				cfg.QdrantColl = value
			}
		}
	}

	return cfg
}

// HasAPIKey checks whether a valid API Key was configured
func (c *TestConfig) HasAPIKey() bool {
	return c.APIKey != "" && c.BaseURL != ""
}
