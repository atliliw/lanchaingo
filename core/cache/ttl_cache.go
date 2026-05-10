package cache

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	value     string
	expiresAt time.Time
}

func (e *entry) isExpired() bool {
	return !e.expiresAt.IsZero() && time.Now().After(e.expiresAt)
}

// TTLCache is an in-memory cache with TTL support.
type TTLCache struct {
	mu       sync.RWMutex
	items    map[string]*entry
	defaultTTL time.Duration
	stopCh   chan struct{}
}

func NewTTLCache(defaultTTL time.Duration) *TTLCache {
	c := &TTLCache{
		items:      make(map[string]*entry),
		defaultTTL: defaultTTL,
		stopCh:     make(chan struct{}),
	}
	go c.evictLoop()
	return c
}

func (c *TTLCache) Get(_ context.Context, key string) (string, error) {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || e.isExpired() {
		if ok {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}
		return "", nil
	}
	return e.value, nil
}

func (c *TTLCache) Set(_ context.Context, key, value string, ttl time.Duration) error {
	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &entry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

func (c *TTLCache) Has(_ context.Context, key string) bool {
	c.mu.RLock()
	e, ok := c.items[key]
	c.mu.RUnlock()

	if !ok || e.isExpired() {
		if ok {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}
		return false
	}
	return true
}

func (c *TTLCache) Remove(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
	return nil
}

func (c *TTLCache) Clear(_ context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*entry)
	return nil
}

func (c *TTLCache) Keys(_ context.Context) ([]string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]string, 0, len(c.items))
	now := time.Now()
	for k, e := range c.items {
		if e.isExpired() {
			continue
		}
		_ = now
		keys = append(keys, k)
	}
	return keys, nil
}

func (c *TTLCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

func (c *TTLCache) evictLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.evictExpired()
		case <-c.stopCh:
			return
		}
	}
}

func (c *TTLCache) evictExpired() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	for k, e := range c.items {
		if !e.expiresAt.IsZero() && now.After(e.expiresAt) {
			delete(c.items, k)
		}
	}
}

// Close stops the eviction loop.
func (c *TTLCache) Close() {
	close(c.stopCh)
}
