package react

import (
	"testing"
	"github.com/atliliw/lanchaingo/agents"
)

func TestParseAction(t *testing.T) {
	p := NewReActOutputParser()
	text := "Thought: need calc\nAction: calculator\nAction Input: 2+2"
	r, err := p.Parse(text)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	a := r.(*agents.AgentActionOutput)
	if a.Action.Tool != "calculator" { t.Errorf("expected calculator") }
}
func TestParseFinalAnswer(t *testing.T) {
	p := NewReActOutputParser()
	text := "Thought: I know\nFinal Answer: The answer is 42"
	r, err := p.Parse(text)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	f := r.(*agents.AgentFinishOutput)
	if f.Finish.Output != "The answer is 42" { t.Errorf("expected answer") }
}
func TestParseInvalid(t *testing.T) {
	p := NewReActOutputParser()
	_, err := p.Parse("invalid output")
	if err == nil { t.Error("expected error") }
}
func TestBuildReActPrompt(t *testing.T) {
	p := BuildReActPrompt("calc tool", []string{"calc"}, "2+2", "")
	if !containsStr(p, "2+2") { t.Error("should contain input") }
	if !containsStr(p, "[calc]") { t.Error("should contain tool names") }
}
func TestFormatScratchpad(t *testing.T) {
	steps := []agents.AgentStep{{Action: agents.AgentAction{Tool:"calc",ToolInput:agents.StringToolInput("2+2")}, Observation:"4"}}
	s := FormatScratchpad(steps)
	if !containsStr(s, "calc") { t.Error("should contain tool name") }
}
