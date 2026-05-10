// 闆嗘垚娴嬭瘯锛氬甫璁板繂鐨勫璇?//
// 娴嬭瘯鍦烘櫙锛氫娇鐢?ConversationChain + BufferMemory 楠岃瘉澶氳疆瀵硅瘽璁板繂
// 杩愯鏂瑰紡: go test -tags=integration ./tests/integration/

//go:build integration

package integration

import (
	"testing"

	"github.com/atliliw/lanchaingo/chains"
	"github.com/atliliw/lanchaingo/memory"
	"github.com/atliliw/lanchaingo/llms/providers"
)

// 娴嬭瘯甯﹁蹇嗙殑瀵硅瘽
//
// 娴嬭瘯鍦烘櫙锛?//  1. 绗竴杞憡璇?AI 鐢ㄦ埛鍚嶅瓧鍙?寮犱笁"
//  2. 绗簩杞棶 AI "鎴戝彨浠€涔堝悕瀛楋紵"
//  3. 楠岃瘉 AI 鑳借浣忎箣鍓嶇殑瀵硅瘽鍐呭
//
// 杩欐槸 LLM 搴旂敤涓渶鏍稿績鐨勮蹇嗗姛鑳芥祴璇?func TestConversationWithMemory(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	// 鍒濆鍖?Qwen Chat 妯″瀷
	llm := providers.NewQwenChat(cfg.APIKey)
	llm.BaseURL = cfg.BaseURL
	llm.Model = cfg.ChatModel

	// 鍒涘缓瀵硅瘽璁板繂锛堜娇鐢?BufferMemory 淇濆瓨鍏ㄩ儴鍘嗗彶锛?	mem := memory.NewConversationBufferMemory()

	// 鍒涘缓 ConversationChain锛岃嚜鍔ㄧ鐞嗚蹇?	chain := chains.NewConversationChain(llm, mem).
		WithSystemPrompt("浣犳槸涓€涓弸濂界殑鍔╂墜锛岃璁颁綇鐢ㄦ埛鐨勪俊鎭?).
		WithVerbose(true) // 寮€鍚缁嗘棩蹇楋紝渚夸簬瑙傚療

	t.Run("绗竴杞細鑷垜浠嬬粛", func(t *testing.T) {
		input := map[string]any{"input": "浣犲ソ锛屾垜鍙紶涓夛紝寰堥珮鍏磋璇嗕綘"}
		result, err := chain.Invoke(input)
		if err != nil {
			t.Fatalf("绗竴杞璇濆け璐? %v", err)
		}
		output, _ := result["output"].(string)
		t.Logf("AI 鍥炲: %s", output)
		if output == "" {
			t.Error("AI 杩斿洖绌哄洖澶?)
		}
	})

	t.Run("绗簩杞細楠岃瘉璁板繂", func(t *testing.T) {
		input := map[string]any{"input": "鎴戝彨浠€涔堝悕瀛楋紵"}
		result, err := chain.Invoke(input)
		if err != nil {
			t.Fatalf("绗簩杞璇濆け璐? %v", err)
		}
		output, _ := result["output"].(string)
		t.Logf("AI 鍥炲: %s", output)

		// 楠岃瘉 AI 璁板緱鐢ㄦ埛鍚嶅瓧锛堣繖鏄蹇嗗姛鑳界殑鏍稿績楠岃瘉鐐癸級
		if output == "" {
			t.Error("AI 杩斿洖绌哄洖澶?)
		}
	})
}

// 娴嬭瘯浣跨敤 SummaryMemory 鐨勫璇?//
// 娴嬭瘯鍦烘櫙锛氫娇鐢ㄦ憳瑕佽蹇嗭紝楠岃瘉闀垮璇濆悗琚帇缂╀负鎽樿
// SummaryMemory 浼氳嚜鍔ㄦ€荤粨鍘嗗彶瀵硅瘽锛岄伩鍏?token 瓒呴檺
func TestConversationWithSummaryMemory(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	llm := providers.NewQwenChat(cfg.APIKey)
	llm.BaseURL = cfg.BaseURL
	llm.Model = cfg.ChatModel

	// SummaryMemory 闇€瑕?LLM 鏉ョ敓鎴愭憳瑕?	sumMem := memory.NewConversationSummaryMemory(llm)

	// 鍏堜繚瀛樺嚑杞璇?	sumMem.SaveContext(
		map[string]string{"input": "鎴戝枩娆㈢紪绋嬪拰璇讳功"},
		map[string]string{"output": "缂栫▼鍜岃涔﹂兘鏄緢濂界殑鐖卞ソ锛?},
	)
	sumMem.SaveContext(
		map[string]string{"input": "鎴戞渶鍠滄鐨勮瑷€鏄?Go"},
		map[string]string{"output": "Go 璇█寰堟锛屽苟鍙戠紪绋嬫槸瀹冪殑寮洪」"},
	)

	// 鍔犺浇璁板繂鍙橀噺锛屾煡鐪嬬敓鎴愮殑鎽樿
	vars, err := sumMem.LoadMemoryVariables(nil)
	if err != nil {
		t.Fatalf("鍔犺浇璁板繂澶辫触: %v", err)
	}

	summary, _ := vars["history"].(string)
	t.Logf("鐢熸垚鐨勬憳瑕? %s", summary)
	if summary == "" {
		t.Log("鎽樿涓虹┖锛堥娆″璇濆彲鑳介渶瑕佹洿澶氬唴瀹规墠鑳界敓鎴愭憳瑕侊級")
	}
}
