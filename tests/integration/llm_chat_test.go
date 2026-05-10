// 闆嗘垚娴嬭瘯锛歀LM 瀵硅瘽鑳藉姏
//
// 杩欎簺娴嬭瘯闇€瑕佺湡瀹炵殑 API Key锛岄粯璁よ鎺掗櫎銆?// 杩愯鏂瑰紡: go test -tags=integration ./tests/integration/
//
// 鍓嶇疆鏉′欢: config/config.toml 涓厤缃簡鏈夋晥鐨?openai.api_key 鍜?base_url

//go:build integration

package integration

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/atliliw/lanchaingo/schema/messages"

	// 浣跨敤 OpenAI 鍏煎鐨?providers 鍖?	"github.com/atliliw/lanchaingo/llms/providers"
)

// 娴嬭瘯 Qwen-Plus 鍩虹瀵硅瘽鍔熻兘
//
// 娴嬭瘯鍦烘櫙锛氬彂閫佷竴鏉＄畝鍗曠殑浜虹被娑堟伅锛岄獙璇?LLM 鑳借繑鍥為潪绌哄搷搴?// 浣跨敤妯″瀷锛歲wen-plus锛堥樋閲岄€氫箟鍗冮棶锛?func TestQwenChatBasic(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	// 鍒涘缓 QwenChat 瀹炰緥锛屼娇鐢ㄩ厤缃腑鐨?API Key 鍜?BaseURL
	llm := providers.NewQwenChat(cfg.APIKey)
	llm.BaseURL = cfg.BaseURL
	llm.Model = cfg.ChatModel

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 鍙戦€佹祴璇曟秷鎭?	msg := messages.NewHumanMessage("浣犲ソ锛岃鐢ㄤ竴鍙ヨ瘽浠嬬粛浣犺嚜宸?)
	result, err := llm.Chat(ctx, []messages.Message{msg})
	if err != nil {
		t.Fatalf("LLM 瀵硅瘽澶辫触: %v", err)
	}

	// 楠岃瘉杩斿洖缁撴灉涓嶄负绌?	if result.Content == "" {
		t.Error("LLM 杩斿洖鍐呭涓虹┖")
	}
	t.Logf("妯″瀷: %s\n鍝嶅簲: %s\nToken 鐢ㄩ噺: %+v", result.Model, result.Content, result.TokenUsage)
}

// 娴嬭瘯娴佸紡瀵硅瘽鍔熻兘
//
// 娴嬭瘯鍦烘櫙锛氫娇鐢?StreamChat 鑾峰彇娴佸紡鍝嶅簲锛岄獙璇佽兘鏀跺埌鑷冲皯涓€涓?chunk
// 娉ㄦ剰锛氭祦寮忔ā寮忎笅閫愬潡杩斿洖鍐呭
func TestQwenChatStreaming(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	// 鍒涘缓 QwenChat 瀹炰緥
	llm := providers.NewQwenChat(cfg.APIKey)
	llm.BaseURL = cfg.BaseURL
	llm.Model = cfg.ChatModel

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 鍙戦€佹祦寮忓璇濊姹?	msg := messages.NewHumanMessage("缁欐垜璁蹭釜绗戣瘽")
	ch, err := llm.StreamChat(ctx, []messages.Message{msg})
	if err != nil {
		t.Fatalf("娴佸紡瀵硅瘽澶辫触: %v", err)
	}

	// 鏀堕泦娴佸紡 chunk锛岄獙璇佽嚦灏戞敹鍒颁竴涓?	received := false
	for chunk := range ch {
		received = true
		fmt.Print(chunk) // 瀹炴椂鎵撳嵃娴佸紡杈撳嚭锛屼究浜庤瀵?	}
	if !received {
		t.Error("娴佸紡鍝嶅簲鏈敹鍒颁换浣曟暟鎹?)
	}
}

// 娴嬭瘯澶氳疆瀵硅瘽
//
// 娴嬭瘯鍦烘櫙锛氳繛缁彂閫佷袱鏉℃秷鎭紝楠岃瘉 LLM 鑳芥纭悊瑙ｅ拰鍝嶅簲
// 绗竴鏉★細鑷垜浠嬬粛
// 绗簩鏉★細鍩轰簬鍓嶆枃鎻愰棶
func TestMultiTurnChat(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	llm := providers.NewQwenChat(cfg.APIKey)
	llm.BaseURL = cfg.BaseURL
	llm.Model = cfg.ChatModel

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 绗竴杞璇?	msg1 := messages.NewHumanMessage("鎴戠殑鍚嶅瓧鏄紶涓?)
	result1, err := llm.Chat(ctx, []messages.Message{msg1})
	if err != nil {
		t.Fatalf("绗竴杞璇濆け璐? %v", err)
	}
	t.Logf("绗竴杞? %s", result1.Content)

	// 绗簩杞璇濓紙甯︿笂鍘嗗彶娑堟伅锛?	msg2 := messages.NewHumanMessage("鎴戝彨浠€涔堝悕瀛楋紵")
	history := []messages.Message{
		msg1,
		messages.NewAIMessage(result1.Content),
		msg2,
	}
	result2, err := llm.Chat(ctx, history)
	if err != nil {
		t.Fatalf("绗簩杞璇濆け璐? %v", err)
	}
	t.Logf("绗簩杞? %s", result2.Content)
}
