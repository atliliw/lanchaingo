package openai

import (
	"strings"
	"testing"
)

// 娴嬭瘯 ParseSSE 鍩虹鍔熻兘锛氫袱涓?data 浜嬩欢琚纭В鏋?func TestParseSSE(t *testing.T) {
	input := "data: {\"key\":\"value\"}\n\ndata: {\"num\":42}\n\n"
	reader := strings.NewReader(input)

	events := ParseSSE(reader)
	var collected []StreamEvent
	for e := range events {
		collected = append(collected, e)
	}

	if len(collected) != 2 {
		t.Fatalf("expected 2 events, got %d", len(collected))
	}
	if collected[0].Data != `{"key":"value"}` {
		t.Errorf("unexpected first event data: %s", collected[0].Data)
	}
	if collected[1].Data != `{"num":42}` {
		t.Errorf("unexpected second event data: %s", collected[1].Data)
	}
}

// 娴嬭瘯 ParseSSE 瑙ｆ瀽 event 瀛楁锛氶潪 data 鐨?event type 涔熻兘姝ｅ父璇嗗埆
func TestParseSSEWithEventType(t *testing.T) {
	input := "event: ping\ndata: {}\n\n"
	reader := strings.NewReader(input)

	events := ParseSSE(reader)
	e := <-events

	if e.Event != "ping" {
		t.Errorf("expected ping event, got %s", e.Event)
	}
}

// 娴嬭瘯 ParseSSE 瑙ｆ瀽 id 瀛楁锛歋SE 鐨?ID 灞炴€ф纭彁鍙?func TestParseSSEWithID(t *testing.T) {
	input := "id: 123\ndata: hello\n\n"
	reader := strings.NewReader(input)

	events := ParseSSE(reader)
	e := <-events

	if e.ID != "123" {
		t.Errorf("expected id 123, got %s", e.ID)
	}
}

// 娴嬭瘯 ParseSSE 璺宠繃娉ㄩ噴琛岋細浠?":" 寮€澶寸殑琛屽簲琚拷鐣?func TestParseSSEWithComments(t *testing.T) {
	input := ": comment line\ndata: actual\n\n"
	reader := strings.NewReader(input)

	events := ParseSSE(reader)
	e := <-events

	if e.Data != "actual" {
		t.Errorf("expected actual, got %s", e.Data)
	}
}

// 娴嬭瘯 ParseSSE 澶氳 data锛氬悓涓€涓簨浠剁殑 data 璺ㄨ秺澶氳鏃剁敤 \n 鎷兼帴
func TestParseSSEMultilineData(t *testing.T) {
	input := "data: line1\ndata: line2\n\n"
	reader := strings.NewReader(input)

	events := ParseSSE(reader)
	e := <-events

	if e.Data != "line1\nline2" {
		t.Errorf("expected multiline data, got %q", e.Data)
	}
}

// 娴嬭瘯 ParseSSE 绌烘祦锛氫笉浜х敓浠讳綍浜嬩欢
func TestParseSSEEmptyStream(t *testing.T) {
	reader := strings.NewReader("")
	events := ParseSSE(reader)

	count := 0
	for range events {
		count++
	}
	if count != 0 {
		t.Errorf("expected 0 events from empty stream, got %d", count)
	}
}

// 娴嬭瘯 ParseSSE 灏鹃儴鏃犵┖琛岋細鏈熬鏁版嵁娌℃湁 \n\n 缁撳熬涔熷簲琚彂灏?// SSE 瑙勮寖鍏佽鏈€鍚庝竴鏉′簨浠舵病鏈夊熬閮ㄧ┖琛?func TestParseSSETrailingData(t *testing.T) {
	input := "data: trailing"
	reader := strings.NewReader(input)

	events := ParseSSE(reader)
	e := <-events

	if e.Data != "trailing" {
		t.Errorf("expected trailing, got %s", e.Data)
	}
}

// 娴嬭瘯 isDoneEvent锛氳瘑鍒?OpenAI 娴佸紡缁撴潫鏍囪 "[DONE]"
// 鍚屾椂楠岃瘉鍓嶅悗绌烘牸琚?trim銆佸叾浠栧瓧绗︿覆鍜岄潪绌哄瓧绗︿覆涓嶈璇垽
func TestIsDoneEvent(t *testing.T) {
	if !isDoneEvent("[DONE]") {
		t.Error("expected [DONE] to be done")
	}
	if !isDoneEvent("  [DONE]  ") {
		t.Error("expected trimmed [DONE] to be done")
	}
	if isDoneEvent("not done") {
		t.Error("expected not done")
	}
	if isDoneEvent("") {
		t.Error("expected empty string not done")
	}
}
