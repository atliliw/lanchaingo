//go:build e2e

package e2e

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/atliliw/lanchaingo/chains"
	"github.com/atliliw/lanchaingo/embeddings"
	"github.com/atliliw/lanchaingo/llms/providers"
	"github.com/atliliw/lanchaingo/memory"
	"github.com/atliliw/lanchaingo/retrieval"
	"github.com/atliliw/lanchaingo/schema/messages"
	vs "github.com/atliliw/lanchaingo/vector_stores"
)

type e2eConfig struct {
	APIKey     string
	BaseURL    string
	ChatModel  string
	EmbedModel string
}

func loadE2EConfig() *e2eConfig {
	// 浠庣幆澧冨彉閲忔垨閰嶇疆涓鍙?	return &e2eConfig{
		APIKey:     "sk-a827c856f53e459bbfad5b8e1b962fc7",
		BaseURL:    "https://dashscope.aliyuncs.com/compatible-mode/v1",
		ChatModel:  "qwen-plus",
		EmbedModel: "text-embedding-v1",
	}
}

// 绔埌绔祴璇曪細瀹屾暣鐨?RAG 瀵硅瘽搴旂敤
//
// 娴嬭瘯鍦烘櫙锛?//  1. 鍑嗗鐭ヨ瘑搴撴枃妗ｏ紙鍏充簬 Rust 鍜?Go 鐨勪粙缁嶏級
//  2. 鐢熸垚鏂囨。宓屽叆骞跺瓨鍏ュ悜閲忓瓨鍌?//  3. 鐢ㄦ埛鎻愰棶锛岀郴缁熸绱㈢浉鍏虫枃妗?//  4. 灏嗘绱㈢粨鏋滀綔涓轰笂涓嬫枃锛孡LM 鐢熸垚鍥炵瓟
//  5. 楠岃瘉鍥炵瓟鏄惁鍩轰簬鎻愪緵鐨勭煡璇?//
// 杩愯: go test -tags=e2e ./tests/e2e/
func TestFullRAGConversation(t *testing.T) {
	cfg := loadE2EConfig()

	// 姝ラ1锛氬垵濮嬪寲 LLM
	llm := providers.NewQwenChat(cfg.APIKey)
	llm.BaseURL = cfg.BaseURL
	llm.Model = cfg.ChatModel
	ctx := context.Background()

	// 姝ラ2锛氬噯澶囩煡璇嗗簱鏂囨。
	knowledgeDocs := []struct {
		content  string
		category string
	}{
		{"Rust 鐢?Mozilla 鐮斿彂锛?015 骞村彂甯?1.0锛屼互鍐呭瓨瀹夊叏鍜岄浂鎴愭湰鎶借薄闂诲悕", "缂栫▼璇█"},
		{"Go 璇█鐢?Google 璁捐锛?009 骞村紑婧愶紝浠ョ畝鍗曡娉曞拰 goroutine 骞跺彂妯″瀷钁楃О", "缂栫▼璇█"},
		{"Python 鐢?Guido van Rossum 鍒涘缓锛?991 骞村彂甯冿紝骞挎硾鐢ㄤ簬 AI 鍜屾暟鎹瀛?, "缂栫▼璇█"},
	}

	// 姝ラ3锛氱敓鎴愬祵鍏ュ苟瀛樺偍
	embCfg := embeddings.DefaultOpenAIEmbeddingsConfig().
		WithAPIKey(cfg.APIKey).
		WithModel(cfg.EmbedModel)
	embCfg.BaseURL = cfg.BaseURL
	emb := embeddings.NewOpenAIEmbeddings(embCfg)

	store := vs.NewInMemoryVectorStore()
	for _, kd := range knowledgeDocs {
		doc := vs.NewDocument(kd.content).WithMetadata("category", kd.category)
		vec, err := emb.EmbedQuery(kd.content)
		if err != nil {
			t.Fatalf("宓屽叆鐢熸垚澶辫触: %v", err)
		}
		store.AddDocuments([]vs.Document{doc}, [][]float32{vec})
	}

	// 姝ラ4锛氬垱寤烘绱㈠櫒
	retriever := retrieval.NewSimilarityRetriever(store, emb)

	// 姝ラ5锛氬垱寤哄璇濋摼锛堝甫璁板繂锛?	mem := memory.NewConversationBufferMemory()
	conversation := chains.NewConversationChain(llm, mem).
		WithSystemPrompt("浣犳槸涓€涓煡璇嗗姪鎵嬶紝璇峰熀浜庢彁渚涚殑鐭ヨ瘑鍥炵瓟闂")

	// 姝ラ6锛氭ā鎷熺敤鎴峰璇?	t.Run("鎻愰棶 Rust", func(t *testing.T) {
		query := "Rust 璇█鏈変粈涔堢壒鐐癸紵"

		// 妫€绱㈢浉鍏虫枃妗?		docs, err := retriever.Retrieve(query, 1)
		if err != nil {
			t.Fatalf("妫€绱㈠け璐? %v", err)
		}

		// 鏋勯€犲甫涓婁笅鏂囩殑 prompt
		var knowledge string
		if len(docs) > 0 {
			knowledge = docs[0].Content
		}

		prompt := fmt.Sprintf(
			"璇峰熀浜庝互涓嬬煡璇嗗洖绛旈棶棰橈細\n\n鐭ヨ瘑: %s\n\n闂: %s\n\n璇风敤涓枃鍥炵瓟锛屽鏋滅煡璇嗗簱涓病鏈夌浉鍏充俊鎭紝璇峰瀹炲憡鐭ャ€?,
			knowledge, query,
		)

		msg := messages.NewHumanMessage(prompt)
		result, err := llm.Chat(ctx, []messages.Message{msg})
		if err != nil {
			t.Fatalf("LLM 璋冪敤澶辫触: %v", err)
		}
		t.Logf("Q: %s", query)
		t.Logf("A: %s", result.Content)
	})

	t.Run("鎻愰棶 Go 璇█", func(t *testing.T) {
		query := "Go 璇█鏄皝寮€鍙戠殑锛?

		docs, err := retriever.Retrieve(query, 1)
		if err != nil {
			t.Fatalf("妫€绱㈠け璐? %v", err)
		}

		var knowledge string
		if len(docs) > 0 {
			knowledge = docs[0].Content
		}

		time.Sleep(1 * time.Second) // 閬垮厤闄愭祦

		prompt := fmt.Sprintf(
			"闂: %s\n\n鍙傝€冧俊鎭? %s\n\n璇风敤涓枃绠€娲佸洖绛斻€?,
			query, knowledge,
		)

		msg := messages.NewHumanMessage(prompt)
		result, err := llm.Chat(ctx, []messages.Message{msg})
		if err != nil {
			t.Fatalf("LLM 璋冪敤澶辫触: %v", err)
		}
		t.Logf("Q: %s", query)
		t.Logf("A: %s", result.Content)
	})
}
