package openai

import (
	"strings"
	"testing"
)

func TestParseSSE(t *testing.T) {
	input := "data: {\"key\":\"value\"}\n\ndata: {\"num\":42}\n\n"
	events := ParseSSE(strings.NewReader(input))
	var c []StreamEvent
	for e := range events { c = append(c, e) }
	if len(c) != 2 { t.Fatalf("expected 2 events, got %d", len(c)) }
}
func TestParseSSEWithEventType(t *testing.T) {
	input := "event: ping\ndata: {}\n\n"
	e := <-ParseSSE(strings.NewReader(input))
	if e.Event != "ping" { t.Errorf("expected ping, got %s", e.Event) }
}
func TestParseSSEWithID(t *testing.T) {
	input := "id: 123\ndata: hello\n\n"
	e := <-ParseSSE(strings.NewReader(input))
	if e.ID != "123" { t.Errorf("expected 123, got %s", e.ID) }
}
func TestParseSSEEmptyStream(t *testing.T) {
	count := 0
	for range ParseSSE(strings.NewReader("")) { count++ }
	if count != 0 { t.Errorf("expected 0 events, got %d", count) }
}
func TestIsDoneEvent(t *testing.T) {
	if !isDoneEvent("[DONE]") { t.Error("expected [DONE] to be done") }
	if isDoneEvent("not done") { t.Error("expected not done") }
}
func TestDefaultOpenAIConfig(t *testing.T) {
	cfg := DefaultOpenAIConfig()
	if cfg.Model != defaultModel { t.Errorf("expected %s", defaultModel) }
}
func TestOpenAIConfigChaining(t *testing.T) {
	cfg := DefaultOpenAIConfig().WithAPIKey("sk-key").WithModel("gpt-4")
	if cfg.APIKey != "sk-key" || cfg.Model != "gpt-4" { t.Errorf("config fields mismatch") }
}
