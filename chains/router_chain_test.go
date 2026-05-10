package chains

import (
	"testing"
)

// 娴嬭瘯 RouterChain 鍏抽敭璇嶅尮閰嶈矾鐢?func TestRouterChainKeywordMatch(t *testing.T) {
	mathChain := &mockChain{inputKey: "input", outputKey: "output", transform: func(s string) string {
		return "math: " + s
	}}
	codeChain := &mockChain{inputKey: "input", outputKey: "output", transform: func(s string) string {
		return "code: " + s
	}}
	defaultChain := &mockChain{inputKey: "input", outputKey: "output", transform: func(s string) string {
		return "default: " + s
	}}

	router := NewRouterChain().
		AddRouteWithKeywords("math", "鏁板", mathChain, []string{"璁＄畻", "鏁板"}).
		AddRouteWithKeywords("code", "缂栫▼", codeChain, []string{"浠ｇ爜", "Rust"}).
		WithDefault(defaultChain)

	t.Run("match math keyword", func(t *testing.T) {
		result, err := router.Invoke(map[string]any{"input": "甯垜璁＄畻涓€涓?})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["output"] != "math: 甯垜璁＄畻涓€涓? {
			t.Errorf("unexpected: %v", result["output"])
		}
	})

	t.Run("match code keyword", func(t *testing.T) {
		result, err := router.Invoke(map[string]any{"input": "Rust浠ｇ爜"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["output"] != "code: Rust浠ｇ爜" {
			t.Errorf("unexpected: %v", result["output"])
		}
	})

	t.Run("use default when no match", func(t *testing.T) {
		result, err := router.Invoke(map[string]any{"input": "浣犲ソ"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result["output"] != "default: 浣犲ソ" {
			t.Errorf("unexpected: %v", result["output"])
		}
	})
}

// 娴嬭瘯鏃犲尮閰嶄笖鏃犻粯璁?Chain 鏃惰繑鍥為敊璇?func TestRouterChainNoMatchNoDefault(t *testing.T) {
	router := NewRouterChain().
		AddRouteWithKeywords("math", "鏁板", &mockChain{inputKey: "input", outputKey: "output", transform: func(s string) string { return s }}, []string{"璁＄畻"})

	_, err := router.Invoke(map[string]any{"input": "hello"})
	if err == nil {
		t.Error("expected error when no match and no default")
	}
}

// 娴嬭瘯 RouteDestination 鍒涘缓
func TestRouteDestination(t *testing.T) {
	dest := NewRouteDestination("math", "鏁板澶勭悊", &mockChain{inputKey: "input", outputKey: "output", transform: func(s string) string { return s }}).
		WithKeywords("璁＄畻", "鏁板")

	if dest.Name != "math" || dest.Description != "鏁板澶勭悊" {
		t.Errorf("destination fields mismatch")
	}
	if len(dest.Keywords) != 2 {
		t.Errorf("expected 2 keywords, got %d", len(dest.Keywords))
	}
}
