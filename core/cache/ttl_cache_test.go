package cache

import (
	"context"
	"testing"
	"time"
)

// 娴嬭瘯 NewTTLCache锛氭柊寤虹紦瀛樹负绌猴紝Len() 杩斿洖 0
func TestNewTTLCache(t *testing.T) {
	c := NewTTLCache(time.Minute)
	if c == nil {
		t.Fatal("expected non-nil cache")
	}
	defer c.Close()

	if c.Len() != 0 {
		t.Errorf("expected empty cache, got %d items", c.Len())
	}
}

// 娴嬭瘯 Set + Get锛氬啓鍏ュ悗鑳芥纭鍙栵紝鍊间竴鑷?func TestTTLCacheSetGet(t *testing.T) {
	c := NewTTLCache(time.Minute)
	defer c.Close()

	ctx := context.Background()

	err := c.Set(ctx, "key1", "value1", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	val, err := c.Get(ctx, "key1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "value1" {
		t.Errorf("expected value1, got %s", val)
	}
}

// 娴嬭瘯 Get 涓嶅瓨鍦ㄧ殑 key锛氳繑鍥炵┖瀛楃涓诧紝涓嶆姤閿?func TestTTLCacheGetMissing(t *testing.T) {
	c := NewTTLCache(time.Minute)
	defer c.Close()

	val, err := c.Get(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty string, got %s", val)
	}
}

// 娴嬭瘯 Has锛氬瓨鍦ㄧ殑 key 杩斿洖 true锛屼笉瀛樺湪鐨勮繑鍥?false
func TestTTLCacheHas(t *testing.T) {
	c := NewTTLCache(time.Minute)
	defer c.Close()

	ctx := context.Background()
	c.Set(ctx, "k", "v", time.Minute)

	if !c.Has(ctx, "k") {
		t.Error("expected Has to return true")
	}
	if c.Has(ctx, "missing") {
		t.Error("expected Has to return false for missing key")
	}
}

// 娴嬭瘯 TTL 杩囨湡鏈哄埗锛氱敤 1ns 鐨勮秴鐭?TTL 鍐欏叆锛宻leep 1ms 鍚?Get 搴旇繑鍥炵┖
// 娉ㄦ剰锛氱撼绉掔骇 TTL 鍦?Go timer 绮惧害涓嬪彲鑳界◢鏈夊亸宸紝浣?1ms 鐨?sleep 瓒冲璁?entry 杩囨湡
func TestTTLCacheExpiredEntry(t *testing.T) {
	c := NewTTLCache(time.Minute)
	defer c.Close()

	ctx := context.Background()
	c.Set(ctx, "fast", "value", time.Nanosecond)

	// 绛夊緟纭繚 entry 杩囨湡锛?ns 杩滃皬浜?Go 鐨勮皟搴︾簿搴︼紝浣?1ms sleep 瓒冲锛?	time.Sleep(time.Millisecond)

	val, err := c.Get(ctx, "fast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "" {
		t.Errorf("expected empty for expired entry, got %s", val)
	}
}

// 娴嬭瘯 Remove锛氬垹闄ゅ悗 Has 杩斿洖 false
func TestTTLCacheRemove(t *testing.T) {
	c := NewTTLCache(time.Minute)
	defer c.Close()

	ctx := context.Background()
	c.Set(ctx, "k", "v", time.Minute)
	c.Remove(ctx, "k")

	if c.Has(ctx, "k") {
		t.Error("expected key to be removed")
	}
}

// 娴嬭瘯 Clear锛氭竻绌哄悗鎵€鏈?key 娑堝け锛孡en() 涓?0
func TestTTLCacheClear(t *testing.T) {
	c := NewTTLCache(time.Minute)
	defer c.Close()

	ctx := context.Background()
	c.Set(ctx, "a", "1", time.Minute)
	c.Set(ctx, "b", "2", time.Minute)
	c.Clear(ctx)

	if c.Len() != 0 {
		t.Errorf("expected empty cache after Clear, got %d items", c.Len())
	}
}

// 娴嬭瘯 Keys锛氳繑鍥炴墍鏈夋湭杩囨湡鐨?key锛屾暟閲忔纭?func TestTTLCacheKeys(t *testing.T) {
	c := NewTTLCache(time.Minute)
	defer c.Close()

	ctx := context.Background()
	c.Set(ctx, "k1", "v1", time.Minute)
	c.Set(ctx, "k2", "v2", time.Minute)

	keys, err := c.Keys(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
}

// 娴嬭瘯榛樿 TTL锛歋et 鏃朵紶 0 搴旇浣跨敤鏋勯€犳椂璁剧疆鐨?defaultTTL
func TestTTLCacheDefaultTTL(t *testing.T) {
	c := NewTTLCache(30 * time.Minute)
	defer c.Close()

	ctx := context.Background()
	c.Set(ctx, "k", "v", 0)

	if !c.Has(ctx, "k") {
		t.Error("expected key to exist with default TTL")
	}
}
