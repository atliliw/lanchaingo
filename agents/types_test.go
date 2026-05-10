package agents

import (
	"context"
	"testing"
	ct "github.com/atliliw/lanchaingo/core/tools"
)

type mockAgent struct{ finish bool }
func (a *mockAgent) Plan(_ []AgentStep, _ map[string]string) (AgentOutput, error) {
	if a.finish { return &AgentFinishOutput{Finish: AgentFinish{Output:"done"}}, nil }
	a.finish = true
	return &AgentActionOutput{Action: AgentAction{Tool:"tool", ToolInput:StringToolInput("x")}}, nil
}
func (a *mockAgent) InputKeys() []string { return []string{"input"} }
func (a *mockAgent) GetAllowedTools() []string { return nil }
type mockTool struct{}
func (m *mockTool) Name() string { return "tool" }
func (m *mockTool) Description() string { return "test" }
func (m *mockTool) Run(_ context.Context, s string) (string, error) { return "observed:"+s, nil }

func TestAgentExecutorBasic(t *testing.T) {
	r, _ := NewAgentExecutor(&mockAgent{}, nil).Invoke("test")
	if r != "done" { t.Errorf("expected done, got %s", r) }
}
func TestAgentExecutorWithTool(t *testing.T) {
	a := &mockAgent{}
	exec := NewAgentExecutor(a, []ct.Tool{&mockTool{}})
	r, _ := exec.Invoke("test")
	if r != "done" { t.Errorf("expected done, got %s", r) }
}
func TestAgentExecutorMaxIterations(t *testing.T) {
	a := &mockAgent{finish: false}
	exec := NewAgentExecutor(a, []ct.Tool{&mockTool{}}).WithMaxIterations(1)
	r, _ := exec.Invoke("test")
	if r == "" { t.Error("expected some output") }
}
func TestToolInput(t *testing.T) {
	ti := StringToolInput("hello")
	if ti.String() != "hello" { t.Errorf("expected hello") }
	ti2 := JSONToolInput(map[string]any{"key":"val"})
	if !ti2.IsJSON { t.Error("expected IsJSON") }
}
func TestAgentActionCreation(t *testing.T) {
	a := AgentAction{Tool:"calc", ToolInput:StringToolInput("2+2"), Log:"thought"}
	if a.Tool != "calc" { t.Errorf("expected calc") }
}
func TestAgentFinishCreation(t *testing.T) {
	f := AgentFinish{Output:"42", Log:"done"}
	if f.Output != "42" { t.Errorf("expected 42") }
}
