package cache

import (
	"context"
	"time"
)

// LLMCache defines the interface for caching LLM responses.
// Maps to Rust langchainrust::core::cache::llm_cache::LLMCache.
type LLMCache interface {
	// Get retrieves a cached response by key.
	Get(ctx context.Context, key string) (string, error)

	// Set stores a response with an optional TTL.
	Set(ctx context.Context, key, value string, ttl time.Duration) error

	// Has checks if a key exists in the cache.
	Has(ctx context.Context, key string) bool

	Remove(ctx context.Context, key string) error

	// Clear empties the entire cache.
	Clear(ctx context.Context) error

	// Keys returns all cached keys.
	Keys(ctx context.Context) ([]string, error)
}

// Cacheable defines the interface for values that can be cached.
type Cacheable interface {
	CacheKey() string
	IsCacheable() bool
}
