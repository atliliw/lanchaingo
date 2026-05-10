package core

import (
	"testing"
	"time"
)

// 娴嬭瘯 NewRunnableConfig 鍒涘缓榛樿閰嶇疆锛歍ags 涓虹┖鍒囩墖锛孧etadata 涓洪潪 nil 鐨?map
func TestNewRunnableConfig(t *testing.T) {
	cfg := NewRunnableConfig()
	if cfg == nil {
		t.Fatal("expected non-nil config")
	}
	if len(cfg.Tags) != 0 {
		t.Errorf("expected empty tags, got %v", cfg.Tags)
	}
	if cfg.Metadata == nil {
		t.Error("expected non-nil Metadata map")
	}
}

// 娴嬭瘯 WithTag 閾惧紡璋冪敤鑳芥纭坊鍔?tag
func TestRunnableConfigWithTag(t *testing.T) {
	cfg := NewRunnableConfig().WithTag("test-tag")
	if len(cfg.Tags) != 1 || cfg.Tags[0] != "test-tag" {
		t.Errorf("expected [test-tag], got %v", cfg.Tags)
	}
}

// 娴嬭瘯 WithRunName 閾惧紡璋冪敤鑳芥纭缃?run name
func TestRunnableConfigWithRunName(t *testing.T) {
	cfg := NewRunnableConfig().WithRunName("my-run")
	if cfg.RunName != "my-run" {
		t.Errorf("expected my-run, got %s", cfg.RunName)
	}
}

// 娴嬭瘯 TokenUsage 缁撴瀯浣撳瓧娈佃祴鍊兼纭?func TestTokenUsage(t *testing.T) {
	tu := TokenUsage{PromptTokens: 10, CompletionTokens: 20, TotalTokens: 30}
	if tu.PromptTokens != 10 || tu.CompletionTokens != 20 || tu.TotalTokens != 30 {
		t.Errorf("TokenUsage fields mismatch: %+v", tu)
	}
}

// 娴嬭瘯 LLMResult 缁撴瀯浣擄細Content/Model 瀛楃涓?+ 宓屽 TokenUsage 鎸囬拡
func TestLLMResult(t *testing.T) {
	r := LLMResult{
		Content: "hello",
		Model:   "gpt-4",
		TokenUsage: &TokenUsage{
			PromptTokens: 5, CompletionTokens: 10, TotalTokens: 15,
		},
	}
	if r.Content != "hello" || r.Model != "gpt-4" {
		t.Errorf("LLMResult fields mismatch: %+v", r)
	}
	if r.TokenUsage == nil || r.TokenUsage.TotalTokens != 15 {
		t.Error("TokenUsage not properly set")
	}
}

// 娴嬭瘯 NewDocument 鍒涘缓鏂囨。锛歅ageContent 璧嬪€兼纭紝Metadata map 琚垵濮嬪寲
func TestNewDocument(t *testing.T) {
	doc := NewDocument("test content")
	if doc.PageContent != "test content" {
		t.Errorf("expected test content, got %s", doc.PageContent)
	}
	if doc.Metadata == nil {
		t.Error("expected non-nil Metadata")
	}
}

// 娴嬭瘯 CacheEntry.IsExpired 涓夌鎯呭喌锛?// 1. 鏃犺繃鏈熸椂闂?鈫?涓嶈繃鏈?// 2. 鏈潵杩囨湡鏃堕棿 鈫?涓嶈繃鏈?// 3. 杩囧幓杩囨湡鏃堕棿 鈫?宸茶繃鏈?func TestCacheEntryExpired(t *testing.T) {
	t.Run("no expiry", func(t *testing.T) {
		e := CacheEntry{Value: "hello"}
		if e.IsExpired() {
			t.Error("entry without expiry should not be expired")
		}
	})

	t.Run("future expiry", func(t *testing.T) {
		e := CacheEntry{Value: "hello", ExpiresAt: timePtr(time.Now().Add(time.Hour))}
		if e.IsExpired() {
			t.Error("entry with future expiry should not be expired")
		}
	})

	t.Run("past expiry", func(t *testing.T) {
		e := CacheEntry{Value: "hello", ExpiresAt: timePtr(time.Now().Add(-time.Hour))}
		if !e.IsExpired() {
			t.Error("entry with past expiry should be expired")
		}
	})
}

// timePtr 杈呭姪鍑芥暟锛氬皢 time.Time 杞负 *time.Time锛岀敤浜?CacheEntry.ExpiresAt
func timePtr(t time.Time) *time.Time {
	return &t
}

// 娴嬭瘯 ToolCall 缁撴瀯浣擄細ID/Name/Arguments 涓変釜瀛楁璧嬪€兼纭?func TestToolCallRoundTrip(t *testing.T) {
	tc := ToolCall{ID: "call_123", Name: "calculator", Arguments: `{"a":1}`}
	if tc.ID != "call_123" || tc.Name != "calculator" || tc.Arguments != `{"a":1}` {
		t.Errorf("ToolCall fields mismatch: %+v", tc)
	}
}

// 娴嬭瘯 AgentAction 缁撴瀯浣擄細Tool/ToolInput/Log 瀛楁璧嬪€兼纭?func TestAgentAction(t *testing.T) {
	aa := AgentAction{Tool: "search", ToolInput: "query", Log: "thinking..."}
	if aa.Tool != "search" || aa.ToolInput != "query" {
		t.Errorf("AgentAction fields mismatch: %+v", aa)
	}
}

// 娴嬭瘯 AgentFinish 缁撴瀯浣擄細Output/Log 瀛楁璧嬪€兼纭?func TestAgentFinish(t *testing.T) {
	af := AgentFinish{Output: "done", Log: "steps"}
	if af.Output != "done" {
		t.Errorf("expected done, got %s", af.Output)
	}
}
