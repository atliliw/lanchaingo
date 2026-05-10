package langgraph

import "testing"

// 娴嬭瘯鍒涘缓绌哄浘
func TestNewStateGraph(t *testing.T) {
	g := NewStateGraph()
	if g == nil {
		t.Fatal("expected non-nil graph")
	}
}

// 娴嬭瘯娣诲姞鑺傜偣鍜岃竟锛岀紪璇戞垚鍔?func TestLinearGraph(t *testing.T) {
	g := NewStateGraph()
	g.AddNodeFn("step1", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	g.AddNodeFn("step2", func(s StateSchema) (StateUpdate, error) {
		as := s.(AgentState)
		as.Output = strPtr("done")
		return NewStateUpdate(as), nil
	})
	g.AddEdge(START, "step1")
	g.AddEdge("step1", "step2")
	g.AddEdge("step2", END)

	cg, err := g.Compile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	input := AgentState{Input: "test"}
	result, err := cg.Invoke(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	fs := result.FinalState.(AgentState)
	if fs.Output == nil || *fs.Output != "done" {
		t.Errorf("expected done, got %v", fs.Output)
	}
}

// 娴嬭瘯缂栬瘧绌哄浘杩斿洖閿欒
func TestEmptyGraphCompileError(t *testing.T) {
	g := NewStateGraph()
	_, err := g.Compile()
	if err == nil {
		t.Error("expected error for empty graph")
	}
}

// 娴嬭瘯甯︽潯浠惰矾鐢辩殑鍥?func TestConditionalGraph(t *testing.T) {
	g := NewStateGraph()
	g.AddNodeFn("decide", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	g.AddNodeFn("process_a", func(s StateSchema) (StateUpdate, error) {
		as := s.(AgentState)
		as.Output = strPtr("A")
		return NewStateUpdate(as), nil
	})
	g.AddNodeFn("process_b", func(s StateSchema) (StateUpdate, error) {
		as := s.(AgentState)
		as.Output = strPtr("B")
		return NewStateUpdate(as), nil
	})
	g.AddNodeFn("end_node", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})

	targets := map[string]string{"a": "process_a", "b": "process_b"}
	g.AddConditionalEdges("decide", "router", targets, nil)
	g.SetConditionalRouter("router", func(s StateSchema) string {
		as := s.(AgentState)
		if as.Input == "a" {
			return "a"
		}
		return "b"
	})
	g.AddEdge(START, "decide")
	g.AddEdge("process_a", "end_node")
	g.AddEdge("process_b", "end_node")
	g.AddEdge("end_node", END)

	cg, err := g.Compile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Run("route to A", func(t *testing.T) {
		result, err := cg.Invoke(AgentState{Input: "a"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		fs := result.FinalState.(AgentState)
		if fs.Output == nil || *fs.Output != "A" {
			t.Errorf("expected A, got %v", fs.Output)
		}
	})

	t.Run("route to B", func(t *testing.T) {
		result, err := cg.Invoke(AgentState{Input: "b"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		fs := result.FinalState.(AgentState)
		if fs.Output == nil || *fs.Output != "B" {
			t.Errorf("expected B, got %v", fs.Output)
		}
	})
}

// 娴嬭瘯閫掑綊闄愬埗
func TestRecursionLimit(t *testing.T) {
	g := NewStateGraph()
	count := 0
	g.AddNodeFn("counter", func(s StateSchema) (StateUpdate, error) {
		count++
		return NewStateUpdate(s.CloneState()), nil
	})
	g.AddNodeFn("end_node", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	g.AddEdge(START, "counter")

	cycleTargets := map[string]string{"continue": "counter", "end": "end_node"}
	g.AddConditionalEdges("counter", "router", cycleTargets, nil)
	g.SetConditionalRouter("router", func(s StateSchema) string {
		as := s.(AgentState)
		as.Steps++
		if as.Steps >= 10 {
			return "end"
		}
		return "continue"
	})
	g.AddEdge("end_node", END)

	cg, err := g.Compile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	cg.WithRecursionLimit(5)

	_, err = cg.Invoke(AgentState{Input: "test"})
	if err == nil {
		t.Error("expected recursion limit error")
	}
}

// 娴嬭瘯澶氬垎鏀浘
func TestMultiBranchGraph(t *testing.T) {
	g := NewStateGraph()
	g.AddNodeFn("preprocess", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	g.AddNodeFn("branch_a", func(s StateSchema) (StateUpdate, error) {
		as := s.(AgentState)
		as.Messages = append(as.Messages, "A")
		return NewStateUpdate(as), nil
	})
	g.AddNodeFn("branch_b", func(s StateSchema) (StateUpdate, error) {
		as := s.(AgentState)
		as.Messages = append(as.Messages, "B")
		return NewStateUpdate(as), nil
	})
	g.AddNodeFn("merge", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	g.AddEdge(START, "preprocess")
	g.AddEdge("preprocess", "branch_a")
	g.AddEdge("preprocess", "branch_b")
	g.AddEdge("branch_a", "merge")
	g.AddEdge("branch_b", "merge")
	g.AddEdge("merge", END)

	cg, err := g.Compile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := cg.Invoke(AgentState{Input: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.RecursionCount == 0 {
		t.Error("expected some recursion")
	}
}

// 娴嬭瘯鐘舵€佷笉鍙樿妭鐐?func TestStateUpdateUnchanged(t *testing.T) {
	upd := NoStateUpdate()
	if upd.Update != nil {
		t.Error("expected nil update for unchanged")
	}
}

// 娴嬭瘯 AgentState 鍏嬮殕
func TestAgentStateClone(t *testing.T) {
	as := AgentState{Input: "hello", Steps: 5, Messages: []string{"hi"}}
	clone := as.CloneState().(AgentState)
	if clone.Input != "hello" || clone.Steps != 5 {
		t.Errorf("clone fields mismatch")
	}
	if len(clone.Messages) != 1 || clone.Messages[0] != "hi" {
		t.Errorf("clone messages mismatch")
	}
	clone.Messages[0] = "modified"
	if as.Messages[0] == "modified" {
		t.Error("clone should be independent")
	}
}

func strPtr(s string) *string { return &s }

// 娴嬭瘯 MemoryCheckpointer
func TestMemoryCheckpointer(t *testing.T) {
	mc := NewMemoryCheckpointer()
	state := AgentState{Input: "save me"}
	id, err := mc.Save(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id == "" {
		t.Error("expected non-empty checkpoint id")
	}

	loaded, err := mc.Load(id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if loaded.(AgentState).Input != "save me" {
		t.Errorf("state mismatch")
	}

	ids, err := mc.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ids) != 1 {
		t.Errorf("expected 1 checkpoint, got %d", len(ids))
	}
}
