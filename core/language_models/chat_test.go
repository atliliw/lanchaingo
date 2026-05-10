package language_models

import (
	"context"
	"testing"
	"time"
)

// 娴嬭瘯 NewChatStream锛氶獙璇?Data 鍜?Err 涓や釜 channel 琚纭垵濮嬪寲
func TestNewChatStream(t *testing.T) {
	cs := NewChatStream(0)
	if cs == nil {
		t.Fatal("expected non-nil ChatStream")
	}
	if cs.Data == nil {
		t.Error("expected non-nil Data channel")
	}
	if cs.Err == nil {
		t.Error("expected non-nil Err channel")
	}
}

// 娴嬭瘯 ChatStream 鑳介€氳繃 Data channel 姝ｅ父鏀跺彂娑堟伅
func TestChatStreamSend(t *testing.T) {
	cs := NewChatStream(10)
	go func() {
		cs.Data <- "hello"
		cs.Close()
	}()

	msg := <-cs.Data
	if msg != "hello" {
		t.Errorf("expected hello, got %s", msg)
	}
}

// 娴嬭瘯 Close 鐨勫箓绛夋€э細澶氭璋冪敤涓嶄細 panic
// 锛堝唴閮ㄧ敤 recover 澶勭悊閲嶅鍏抽棴 channel 鐨?panic锛?func TestChatStreamClose(t *testing.T) {
	cs := NewChatStream(10)
	cs.Close()
	cs.Close()
}

// 娴嬭瘯 ModelConfig.Validate 涓夌鍦烘櫙锛?// 1. APIKey 鍜?Model 閮藉～ 鈫?閫氳繃
// 2. 缂?APIKey 鈫?鎶ラ敊 ErrMissingAPIKey
// 3. 缂?Model 鈫?鎶ラ敊 ErrMissingModel
func TestModelConfigValidate(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		cfg := ModelConfig{APIKey: "sk-xxx", Model: "gpt-4"}
		if err := cfg.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("missing API key", func(t *testing.T) {
		cfg := ModelConfig{Model: "gpt-4"}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing API key")
		}
	})

	t.Run("missing model", func(t *testing.T) {
		cfg := ModelConfig{APIKey: "sk-xxx"}
		if err := cfg.Validate(); err == nil {
			t.Error("expected error for missing model")
		}
	})
}

// 娴嬭瘯 LLMResult 绫诲瀷鍒悕锛氶獙璇佷粠 core 鍖呴噸瀵煎嚭鐨勭被鍨嬭兘姝ｅ父浣跨敤
func TestLLMResultTypeAlias(t *testing.T) {
	r := LLMResult{Content: "test", Model: "gpt-4"}
	if r.Content != "test" || r.Model != "gpt-4" {
		t.Errorf("LLMResult fields mismatch: %+v", r)
	}
}

// 娴嬭瘯 ChatStream 鍦ㄦ棤鏁版嵁鍙戦€佹椂锛岄€氳繃 context 瓒呮椂鑳芥纭€€鍑?// 鍦烘櫙锛氬垱寤?ChatStream锛屼笉鍙戦€佷换浣曟暟鎹紝璁剧疆 1ms 瓒呮椂鐨?context
func TestChatStreamTimeout(t *testing.T) {
	cs := NewChatStream(10)
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	select {
	case <-cs.Data:
		t.Error("unexpected data")
	case <-ctx.Done():
	}
}
