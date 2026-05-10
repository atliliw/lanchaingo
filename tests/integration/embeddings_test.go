// 闆嗘垚娴嬭瘯锛欵mbedding 鍚戦噺鐢熸垚
//
// 杩欎簺娴嬭瘯闇€瑕佺湡瀹炵殑 API Key锛岄粯璁よ鎺掗櫎銆?// 杩愯鏂瑰紡: go test -tags=integration ./tests/integration/
//
// 鍓嶇疆鏉′欢: config/config.toml 涓厤缃簡鏈夋晥鐨?embedding 鍙傛暟

//go:build integration

package integration

import (
	"testing"

	"github.com/atliliw/lanchaingo/embeddings"
)

// 娴嬭瘯鏂囨湰宓屽叆鐢熸垚
//
// 娴嬭瘯鍦烘櫙锛氫娇鐢?text-embedding-v1 妯″瀷鐢熸垚鏂囨湰鍚戦噺
// 楠岃瘉鐐癸細
//   - 鍚戦噺缁村害涓?1536锛坱ext-embedding-v1 鐨勯粯璁ょ淮搴︼級
//   - 鍚戦噺鍊间笉涓虹┖
func TestEmbedQuery(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	// 鍒涘缓 OpenAI Embeddings 瀹㈡埛绔紝杩炴帴鍒?DashScope
	embCfg := embeddings.DefaultOpenAIEmbeddingsConfig().
		WithAPIKey(cfg.APIKey).
		WithModel(cfg.EmbedModel)
	embCfg.BaseURL = cfg.BaseURL

	emb := embeddings.NewOpenAIEmbeddings(embCfg)

	// 鐢熸垚鍗曚釜鏂囨湰鐨勫祵鍏ュ悜閲?	result, err := emb.EmbedQuery("浠婂ぉ澶╂皵鎬庝箞鏍凤紵")
	if err != nil {
		t.Fatalf("鐢熸垚宓屽叆鍚戦噺澶辫触: %v", err)
	}

	// 楠岃瘉鍚戦噺缁村害
	if len(result) == 0 {
		t.Error("宓屽叆鍚戦噺涓虹┖")
	} else {
		t.Logf("鍚戦噺缁村害: %d锛屽墠5涓€? %v", len(result), result[:5])
	}
}

// 娴嬭瘯鎵归噺鏂囨。宓屽叆
//
// 娴嬭瘯鍦烘櫙锛氬悓鏃朵负澶氫釜鏂囨。鐢熸垚宓屽叆鍚戦噺
// 楠岃瘉涓や釜鏂囨。杩斿洖鐨勫悜閲忕淮搴︿竴鑷?
func TestEmbedDocuments(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	embCfg := embeddings.DefaultOpenAIEmbeddingsConfig().
		WithAPIKey(cfg.APIKey).
		WithModel(cfg.EmbedModel)
	embCfg.BaseURL = cfg.BaseURL

	emb := embeddings.NewOpenAIEmbeddings(embCfg)

	// 涓哄涓枃妗ｆ壒閲忕敓鎴愬悜閲?	texts := []string{
		"Rust 鏄竴闂ㄧ郴缁熺紪绋嬭瑷€",
		"Scripting 鏄剼鏈瑷€",
		"Go 璇█鎿呴暱骞跺彂缂栫▼",
	}

	results, err := emb.EmbedDocuments(texts)
	if err != nil {
		t.Fatalf("鎵归噺宓屽叆澶辫触: %v", err)
	}

	// 楠岃瘉杩斿洖缁撴灉鏁伴噺
	if len(results) != len(texts) {
		t.Fatalf("鏈熸湜 %d 涓悜閲忥紝瀹為檯寰楀埌 %d", len(texts), len(results))
	}

	t.Logf("姣忎釜鍚戦噺缁村害: %d", len(results[0]))
}
