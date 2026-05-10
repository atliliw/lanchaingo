// 闆嗘垚娴嬭瘯锛歊AG 妫€绱㈠寮虹敓鎴?//
// 娴嬭瘯瀹屾暣鐨?RAG 绠＄嚎锛氭枃妗ｅ祵鍏?鈫?鍚戦噺瀛樺偍 鈫?鐩镐技搴︽绱?鈫?LLM 鐢熸垚
// 杩愯鏂瑰紡: go test -tags=integration ./tests/integration/

//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/atliliw/lanchaingo/embeddings"
	"github.com/atliliw/lanchaingo/llms/providers"
	"github.com/atliliw/lanchaingo/retrieval"
	"github.com/atliliw/lanchaingo/schema/messages"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

// 娴嬭瘯瀹屾暣鐨?RAG 绠＄嚎
//
// 娴嬭瘯鍦烘櫙锛?//  1. 鍑嗗娴嬭瘯鏂囨。锛堝叧浜?Rust銆丳ython銆丟o 鐨勪粙缁嶏級
//  2. 浣跨敤 text-embedding-v1 妯″瀷鐢熸垚鏂囨。鍚戦噺
//  3. 灏嗘枃妗?鍚戦噺瀛樺叆 InMemoryVectorStore
//  4. 浣跨敤 SimilarityRetriever 妫€绱㈢浉鍏虫枃妗?//  5. 楠岃瘉妫€绱㈢粨鏋滃噯纭?//
// 杩欐槸 RAG 搴旂敤鏈€鏍稿績鐨勭鍒扮娴佺▼
func TestRAGPipeline(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	// 姝ラ1锛氬噯澶囨枃妗ｂ€斺€斾笁涓富棰樺悇寮傜殑鏂囨。鐢ㄤ簬楠岃瘉妫€绱㈠噯纭€?	docs := []vs.Document{
		vs.NewDocument("Rust 鏄竴闂ㄧ郴缁熺紪绋嬭瑷€锛屾敞閲嶅唴瀛樺畨鍏ㄥ拰骞跺彂鎬ц兘锛岀敱 Mozilla 寮€鍙?).
			WithMetadata("source", "doc1"),
		vs.NewDocument("Python 鏄竴绉嶈В閲婂瀷楂樼骇缂栫▼璇█锛岃娉曠畝娲侊紝閫傚悎蹇€熷紑鍙戝拰鏁版嵁绉戝").
			WithMetadata("source", "doc2"),
		vs.NewDocument("Go 璇█鐢?Google 寮€鍙戯紝浠ョ畝娲佺殑璇硶鍜屽己澶х殑骞跺彂妯″瀷钁楃О").
			WithMetadata("source", "doc3"),
	}

	// 姝ラ2锛氱敓鎴愭枃妗ｅ祵鍏ュ悜閲?	embCfg := embeddings.DefaultOpenAIEmbeddingsConfig().
		WithAPIKey(cfg.APIKey).
		WithModel(cfg.EmbedModel)
	embCfg.BaseURL = cfg.BaseURL
	emb := embeddings.NewOpenAIEmbeddings(embCfg)

	texts := make([]string, len(docs))
	for i, d := range docs {
		texts[i] = d.Content
	}
	vecs, err := emb.EmbedDocuments(texts)
	if err != nil {
		t.Fatalf("鐢熸垚鏂囨。宓屽叆鍚戦噺澶辫触: %v", err)
	}
	t.Logf("鏂囨。宓屽叆鍚戦噺鐢熸垚鎴愬姛锛岀淮搴? %d", len(vecs[0]))

	// 姝ラ3锛氬瓨鍏ュ悜閲忓瓨鍌?	store := vs.NewInMemoryVectorStore()
	ids, err := store.AddDocuments(docs, vecs)
	if err != nil {
		t.Fatalf("娣诲姞鏂囨。鍒板悜閲忓瓨鍌ㄥけ璐? %v", err)
	}
	t.Logf("鏂囨。宸插瓨鍏ュ悜閲忓瓨鍌紝IDs: %v", ids)

	// 姝ラ4+5锛氬垱寤烘绱㈠櫒骞舵绱?	retriever := retrieval.NewSimilarityRetriever(store, emb)

	t.Run("妫€绱?Rust 鐩稿叧鏂囨。", func(t *testing.T) {
		results, err := retriever.RetrieveWithScores("绯荤粺缂栫▼ 鍐呭瓨瀹夊叏", 2)
		if err != nil {
			t.Fatalf("妫€绱㈠け璐? %v", err)
		}
		if len(results) == 0 {
			t.Fatal("妫€绱㈢粨鏋滀负绌?)
		}
		t.Logf("妫€绱㈠埌 %d 涓粨鏋?", len(results))
		for i, r := range results {
			t.Logf("  [%d] 鍒嗘暟=%.4f: %s", i+1, r.Score, r.Document.Content)
		}
	})

	t.Run("妫€绱?Python 鐩稿叧鏂囨。", func(t *testing.T) {
		results, err := retriever.RetrieveWithScores("鏁版嵁绉戝 鑴氭湰璇█", 2)
		if err != nil {
			t.Fatalf("妫€绱㈠け璐? %v", err)
		}
		if len(results) == 0 {
			t.Fatal("妫€绱㈢粨鏋滀负绌?)
		}
		t.Logf("妫€绱㈠埌 %d 涓粨鏋?", len(results))
		for i, r := range results {
			t.Logf("  [%d] 鍒嗘暟=%.4f: %s", i+1, r.Score, r.Document.Content)
		}
	})
}

// 娴嬭瘯 RAG + LLM 闂瓟
//
// 娴嬭瘯鍦烘櫙锛?//  1. 鍏堢敤鍚戦噺妫€绱㈡壘鍒扮浉鍏虫枃妗?//  2. 灏嗘绱㈢粨鏋滀綔涓轰笂涓嬫枃杈撳叆 LLM
//  3. LLM 鍩轰簬涓婁笅鏂囧洖绛旈棶棰?//
// 杩欐槸鐢熶骇鐜涓?RAG 搴旂敤鐨勫畬鏁村舰鎬?
func TestRAGWithQA(t *testing.T) {
	cfg := LoadConfig()
	if !cfg.HasAPIKey() {
		t.Skip("璺宠繃闆嗘垚娴嬭瘯锛氭湭閰嶇疆 API Key")
	}

	// 鍑嗗 LLM 鍜屽悜閲忕粍浠?	llm := providers.NewQwenChat(cfg.APIKey)
	llm.BaseURL = cfg.BaseURL
	llm.Model = cfg.ChatModel

	docs := []vs.Document{
		vs.NewDocument("Rust 鐢?Mozilla 鐮斿彂锛?015 骞村彂甯?1.0 鐗堟湰锛屼互鍏舵墍鏈夋潈绯荤粺鍜屽€熺敤妫€鏌ュ櫒闂诲悕"),
		vs.NewDocument("Go 鐢?Google 鐨?Robert Griesemer銆丷ob Pike 鍜?Ken Thompson 璁捐锛?009 骞撮娆″叕寮€"),
	}

	embCfg := embeddings.DefaultOpenAIEmbeddingsConfig().
		WithAPIKey(cfg.APIKey).
		WithModel(cfg.EmbedModel)
	embCfg.BaseURL = cfg.BaseURL
	emb := embeddings.NewOpenAIEmbeddings(embCfg)

	store := vs.NewInMemoryVectorStore()
	texts := []string{docs[0].Content, docs[1].Content}
	vecs, _ := emb.EmbedDocuments(texts)
	store.AddDocuments(docs, vecs)

	retriever := retrieval.NewSimilarityRetriever(store, emb)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 妫€绱㈢浉鍏虫枃妗?	query := "Rust 璇█鏄皝寮€鍙戠殑锛?
	results, err := retriever.Retrieve(query, 1)
	if err != nil {
		t.Fatalf("妫€绱㈠け璐? %v", err)
	}

	if len(results) > 0 {
		contextStr := results[0].Content
		prompt := "鏍规嵁浠ヤ笅淇℃伅鍥炵瓟闂銆俓n\n宸茬煡淇℃伅: " + contextStr + "\n\n闂: " + query

		// 灏嗘绱㈢粨鏋滀綔涓轰笂涓嬫枃锛岃 LLM 鐢熸垚绛旀
		msg := messages.NewHumanMessage(prompt)
		result, err := llm.Chat(ctx, []messages.Message{msg})
		if err != nil {
			t.Fatalf("LLM 鍥炵瓟澶辫触: %v", err)
		}
		t.Logf("闂: %s", query)
		t.Logf("绛旀: %s", result.Content)
	}
}
