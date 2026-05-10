package langgraph

import "testing"

func TestNewStateGraph(t *testing.T) {
	g := NewStateGraph()
	if g == nil { t.Fatal("expected non-nil graph") }
}
func TestLinearGraph(t *testing.T) {
	g := NewStateGraph()
	g.AddNodeFn("step", func(s StateSchema) (StateUpdate, error) { return NewStateUpdate(s.CloneState()), nil })
	g.AddEdge(START, "step")
	g.AddEdge("step", END)
	cg, err := g.Compile()
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	_, err = cg.Invoke(AgentState{Input: "test"})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
}
func TestEmptyGraphCompileError(t *testing.T) {
	_, err := NewStateGraph().Compile()
	if err == nil { t.Error("expected error for empty graph") }
}
func TestConditionalGraph(t *testing.T) {
	g := NewStateGraph()
	g.AddNodeFn("decide", func(s StateSchema) (StateUpdate, error) { return NewStateUpdate(s.CloneState()), nil })
	g.AddNodeFn("end", func(s StateSchema) (StateUpdate, error) { return NewStateUpdate(s.CloneState()), nil })
	targets := map[string]string{"a": "end"}
	g.AddConditionalEdges("decide", "router", targets, nil)
	g.SetConditionalRouter("router", func(s StateSchema) string { return "a" })
	g.AddEdge(START, "decide")
	g.AddEdge("end", END)
	_, err := g.Compile()
	if err != nil { t.Fatalf("unexpected error: %v", err) }
}
func TestAgentStateClone(t *testing.T) {
	as := AgentState{Input: "hello", Steps: 5, Messages: []string{"hi"}}
	clone := as.CloneState().(AgentState)
	if clone.Input != "hello" { t.Errorf("clone fields mismatch") }
}
func TestMemoryCheckpointer(t *testing.T) {
	mc := NewMemoryCheckpointer()
	id, _ := mc.Save(AgentState{Input: "save me"})
	if id == "" { t.Error("expected non-empty id") }
	ids, _ := mc.List()
	if len(ids) != 1 { t.Errorf("expected 1 checkpoint, got %d", len(ids)) }
}
