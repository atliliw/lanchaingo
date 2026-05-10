package cache

import (
	"context"
	"testing"
	"time"
)

func TestNewTTLCache(t *testing.T) {
	c := NewTTLCache(time.Minute); defer c.Close()
	if c.Len() != 0 { t.Errorf("expected empty cache, got %d items", c.Len()) }
}
func TestTTLCacheSetGet(t *testing.T) {
	c := NewTTLCache(time.Minute); defer c.Close()
	ctx := context.Background()
	c.Set(ctx, "key1", "value1", time.Minute)
	val, _ := c.Get(ctx, "key1")
	if val != "value1" { t.Errorf("expected value1, got %s", val) }
}
func TestTTLCacheGetMissing(t *testing.T) {
	c := NewTTLCache(time.Minute); defer c.Close()
	val, _ := c.Get(context.Background(), "nonexistent")
	if val != "" { t.Errorf("expected empty string, got %s", val) }
}
func TestTTLCacheHas(t *testing.T) {
	c := NewTTLCache(time.Minute); defer c.Close()
	ctx := context.Background()
	c.Set(ctx, "k", "v", time.Minute)
	if !c.Has(ctx, "k") { t.Error("expected Has to return true") }
}
func TestTTLCacheRemove(t *testing.T) {
	c := NewTTLCache(time.Minute); defer c.Close()
	ctx := context.Background()
	c.Set(ctx, "k", "v", time.Minute)
	c.Remove(ctx, "k")
	if c.Has(ctx, "k") { t.Error("expected key to be removed") }
}
func TestTTLCacheClear(t *testing.T) {
	c := NewTTLCache(time.Minute); defer c.Close()
	ctx := context.Background()
	c.Set(ctx, "a", "1", time.Minute)
	c.Set(ctx, "b", "2", time.Minute)
	c.Clear(ctx)
	if c.Len() != 0 { t.Errorf("expected 0 after Clear, got %d", c.Len()) }
}
func TestTTLCacheKeys(t *testing.T) {
	c := NewTTLCache(time.Minute); defer c.Close()
	ctx := context.Background()
	c.Set(ctx, "k1", "v1", time.Minute)
	c.Set(ctx, "k2", "v2", time.Minute)
	keys, _ := c.Keys(ctx)
	if len(keys) != 2 { t.Errorf("expected 2 keys, got %d", len(keys)) }
}
