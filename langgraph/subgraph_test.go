package langgraph

import (
	"errors"
	"strings"
	"sync"
	"testing"
)

type MapState map[string]any

func (m MapState) CloneState() StateSchema {
	clone := make(MapState, len(m))
	for k, v := range m {
		clone[k] = v
	}
	return clone
}

func TestNewSubgraphNode(t *testing.T) {
	sub := NewStateGraph()
	sub.AddNodeFn("n", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	sub.AddEdge(START, "n")
	sub.AddEdge("n", END)
	c, err := sub.Compile()
	if err != nil {
		t.Fatal(err)
	}
	sn := NewSubgraphNode("test", c)
	if sn.Name() != "test" {
		t.Errorf("expected name 'test', got %q", sn.Name())
	}
	if sn.SubgraphDepth() != 1 {
		t.Errorf("expected depth 1, got %d", sn.SubgraphDepth())
	}
}

func TestNewSubgraphNodePanicsOnNil(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil subgraph")
		}
	}()
	NewSubgraphNode("bad", nil)
}

func TestSubgraphBasic(t *testing.T) {
	sub := NewStateGraph()
	sub.AddNodeFn("sub_node", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "subgraph_executed")
		return NewStateUpdate(as), nil
	})
	sub.AddEdge(START, "sub_node")
	sub.AddEdge("sub_node", END)
	subCompiled, err := sub.Compile()
	if err != nil {
		t.Fatalf("subgraph compile: %v", err)
	}

	parent := NewStateGraph()
	parent.AddNode("sub", NewSubgraphNode("sub", subCompiled))
	parent.AddEdge(START, "sub")
	parent.AddEdge("sub", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	result, err := parentCompiled.Invoke(AgentState{Input: "test"})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if result.Interrupted {
		t.Fatal("unexpected interrupt")
	}

	fs := result.FinalState.(AgentState)
	if len(fs.Messages) != 1 || fs.Messages[0] != "subgraph_executed" {
		t.Errorf("subgraph was not executed, messages: %v", fs.Messages)
	}
}

func TestSubgraphPreservesParentState(t *testing.T) {
	sub := NewStateGraph()
	sub.AddNodeFn("sub_node", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	sub.AddEdge(START, "sub_node")
	sub.AddEdge("sub_node", END)
	subCompiled, err := sub.Compile()
	if err != nil {
		t.Fatalf("subgraph compile: %v", err)
	}

	parent := NewStateGraph()
	parent.AddNodeFn("before", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "before_sub")
		return NewStateUpdate(as), nil
	})
	parent.AddNode("sub", NewSubgraphNode("sub", subCompiled))
	parent.AddNodeFn("after", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "after_sub")
		return NewStateUpdate(as), nil
	})
	parent.AddEdge(START, "before")
	parent.AddEdge("before", "sub")
	parent.AddEdge("sub", "after")
	parent.AddEdge("after", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	result, err := parentCompiled.Invoke(AgentState{Input: "test", Messages: []string{"initial"}})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}

	fs := result.FinalState.(AgentState)
	if len(fs.Messages) != 3 || fs.Messages[0] != "initial" || fs.Messages[1] != "before_sub" || fs.Messages[2] != "after_sub" {
		t.Errorf("expected parent state preserved, messages: %v", fs.Messages)
	}
}

func TestSubgraphInputMap(t *testing.T) {
	sub := NewStateGraph()
	sub.SetReducer(func(current, update StateSchema) StateSchema {
		cur := current.(MapState)
		upd := update.(MapState)
		for k, v := range upd {
			cur[k] = v
		}
		return cur
	})
	sub.AddNodeFn("sub_node", func(s StateSchema) (StateUpdate, error) {
		ms := s.CloneState().(MapState)
		ms["subgraph_key"] = "from_subgraph"
		return NewStateUpdate(ms), nil
	})
	sub.AddEdge(START, "sub_node")
	sub.AddEdge("sub_node", END)
	subCompiled, err := sub.Compile()
	if err != nil {
		t.Fatalf("subgraph compile: %v", err)
	}

	parent := NewStateGraph()
	parent.AddNode("sub", NewSubgraphNode("sub", subCompiled,
		WithInputMap(func(parentState StateSchema) StateSchema {
			as := parentState.(AgentState)
			return MapState{"input": as.Input}
		}),
		WithOutputMap(func(childState, parentState StateSchema) StateSchema {
			child := childState.(MapState)
			parent := parentState.CloneState().(AgentState)
			parent.Messages = append(parent.Messages, child["subgraph_key"].(string))
			return parent
		}),
	))
	parent.AddEdge(START, "sub")
	parent.AddEdge("sub", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	result, err := parentCompiled.Invoke(AgentState{Input: "hello", Messages: []string{"start"}})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}

	fs := result.FinalState.(AgentState)
	if len(fs.Messages) != 2 || fs.Messages[0] != "start" || fs.Messages[1] != "from_subgraph" {
		t.Errorf("state mapping failed, messages: %v", fs.Messages)
	}
}

func TestSubgraphInterruptPropagation(t *testing.T) {
	sub := NewStateGraph()
	sub.AddNodeFn("sub_node", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "sub_executed")
		return NewStateUpdate(as), nil
	})
	sub.AddEdge(START, "sub_node")
	sub.AddEdge("sub_node", END)
	subCompiled, err := sub.Compile()
	if err != nil {
		t.Fatalf("subgraph compile: %v", err)
	}
	subCompiled.InterruptAfter("sub_node")

	parent := NewStateGraph()
	parent.AddNode("sub", NewSubgraphNode("sub", subCompiled))
	parent.AddNodeFn("after", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "should_not_reach")
		return NewStateUpdate(as), nil
	})
	parent.AddEdge(START, "sub")
	parent.AddEdge("sub", "after")
	parent.AddEdge("after", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	result, err := parentCompiled.Invoke(AgentState{Input: "test"})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if !result.Interrupted {
		t.Fatal("expected interrupt propagation from subgraph")
	}

	fs := result.FinalState.(AgentState)
	if len(fs.Messages) != 1 || fs.Messages[0] != "sub_executed" {
		t.Errorf("expected subgraph to execute but not after node, messages: %v", fs.Messages)
	}
}

func TestSubgraphInterruptWithMapping(t *testing.T) {
	sub := NewStateGraph()
	sub.SetReducer(func(cur, upd StateSchema) StateSchema {
		c := cur.(MapState)
		u := upd.(MapState)
		for k, v := range u {
			c[k] = v
		}
		return c
	})
	sub.AddNodeFn("sub_node", func(s StateSchema) (StateUpdate, error) {
		ms := s.CloneState().(MapState)
		ms["partial"] = "interrupted_data"
		return NewStateUpdate(ms), nil
	})
	sub.AddEdge(START, "sub_node")
	sub.AddEdge("sub_node", END)
	subCompiled, err := sub.Compile()
	if err != nil {
		t.Fatalf("subgraph compile: %v", err)
	}
	subCompiled.InterruptAfter("sub_node")

	parent := NewStateGraph()
	parent.AddNode("sub", NewSubgraphNode("sub", subCompiled,
		WithInputMap(func(s StateSchema) StateSchema {
			as := s.(AgentState)
			return MapState{"input": as.Input}
		}),
		WithOutputMap(func(child, parent StateSchema) StateSchema {
			childMap := child.(MapState)
			as := parent.CloneState().(AgentState)
			if v, ok := childMap["partial"]; ok {
				as.Messages = append(as.Messages, v.(string))
			}
			return as
		}),
	))
	parent.AddEdge(START, "sub")
	parent.AddEdge("sub", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	result, err := parentCompiled.Invoke(AgentState{Input: "test", Messages: []string{"start"}})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}
	if !result.Interrupted {
		t.Fatal("expected interrupt propagation with mapping")
	}

	fs := result.FinalState.(AgentState)
	if len(fs.Messages) != 2 || fs.Messages[0] != "start" || fs.Messages[1] != "interrupted_data" {
		t.Errorf("interrupt with mapping failed, messages: %v", fs.Messages)
	}
}

func TestSubgraphErrorPropagation(t *testing.T) {
	sub := NewStateGraph()
	sub.AddNodeFn("failing", func(s StateSchema) (StateUpdate, error) {
		return StateUpdate{}, errors.New("internal_subgraph_error")
	})
	sub.AddEdge(START, "failing")
	sub.AddEdge("failing", END)
	subCompiled, err := sub.Compile()
	if err != nil {
		t.Fatalf("subgraph compile: %v", err)
	}

	parent := NewStateGraph()
	parent.AddNode("sub", NewSubgraphNode("sub", subCompiled))
	parent.AddEdge(START, "sub")
	parent.AddEdge("sub", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	_, err = parentCompiled.Invoke(AgentState{Input: "test"})
	if err == nil {
		t.Fatal("expected error propagation from subgraph")
	}

	if !strings.Contains(err.Error(), "internal_subgraph_error") {
		t.Errorf("expected wrapped error message, got: %v", err)
	}

	var ge *GraphError
	if !errors.As(err, &ge) {
		t.Errorf("expected GraphError, got %T", err)
	}
}

func TestSubgraphNoCycleForValidNesting(t *testing.T) {
	leaf := NewStateGraph()
	leaf.AddNodeFn("leaf", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	leaf.AddEdge(START, "leaf")
	leaf.AddEdge("leaf", END)
	leafCompiled, err := leaf.Compile()
	if err != nil {
		t.Fatalf("leaf compile: %v", err)
	}

	mid := NewStateGraph()
	mid.AddNode("mid_sub", NewSubgraphNode("leaf_ref", leafCompiled))
	mid.AddEdge(START, "mid_sub")
	mid.AddEdge("mid_sub", END)
	midCompiled, err := mid.Compile()
	if err != nil {
		t.Fatalf("mid compile: %v", err)
	}
	if midCompiled == nil {
		t.Fatal("expected valid compilation for tree-like nesting")
	}
}

func TestSubgraphNested(t *testing.T) {
	leaf := NewStateGraph()
	leaf.AddNodeFn("leaf", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "leaf")
		return NewStateUpdate(as), nil
	})
	leaf.AddEdge(START, "leaf")
	leaf.AddEdge("leaf", END)
	leafCompiled, err := leaf.Compile()
	if err != nil {
		t.Fatalf("leaf compile: %v", err)
	}

	mid := NewStateGraph()
	mid.AddNode("mid_sub", NewSubgraphNode("leaf", leafCompiled))
	mid.AddEdge(START, "mid_sub")
	mid.AddEdge("mid_sub", END)
	midCompiled, err := mid.Compile()
	if err != nil {
		t.Fatalf("mid compile: %v", err)
	}

	parent := NewStateGraph()
	parent.AddNode("top_sub", NewSubgraphNode("mid", midCompiled))
	parent.AddEdge(START, "top_sub")
	parent.AddEdge("top_sub", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	result, err := parentCompiled.Invoke(AgentState{Input: "nested", Messages: []string{}})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}

	fs := result.FinalState.(AgentState)
	if len(fs.Messages) != 1 || fs.Messages[0] != "leaf" {
		t.Errorf("nested subgraph execution failed, messages: %v", fs.Messages)
	}
}

func TestSubgraphConcurrentAccess(t *testing.T) {
	sub := NewStateGraph()
	sub.AddNodeFn("sub_node", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "executed")
		return NewStateUpdate(as), nil
	})
	sub.AddEdge(START, "sub_node")
	sub.AddEdge("sub_node", END)
	subCompiled, err := sub.Compile()
	if err != nil {
		t.Fatalf("subgraph compile: %v", err)
	}

	parent := NewStateGraph()
	parent.AddNode("sub", NewSubgraphNode("sub", subCompiled))
	parent.AddEdge(START, "sub")
	parent.AddEdge("sub", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 20)
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := parentCompiled.Invoke(AgentState{Input: "concurrent"})
			if err != nil {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)

	for err := range errCh {
		t.Errorf("concurrent invoke failed: %v", err)
	}
}

func TestSubgraphDepth(t *testing.T) {
	leaf := NewStateGraph()
	leaf.AddNodeFn("leaf", func(s StateSchema) (StateUpdate, error) {
		return NewStateUpdate(s.CloneState()), nil
	})
	leaf.AddEdge(START, "leaf")
	leaf.AddEdge("leaf", END)
	leafCompiled, _ := leaf.Compile()
	leafNode := NewSubgraphNode("leaf", leafCompiled)
	if leafNode.SubgraphDepth() != 1 {
		t.Errorf("leaf depth: expected 1, got %d", leafNode.SubgraphDepth())
	}

	mid := NewStateGraph()
	mid.AddNode("mid_sub", leafNode)
	mid.AddEdge(START, "mid_sub")
	mid.AddEdge("mid_sub", END)
	midCompiled, _ := mid.Compile()
	midNode := NewSubgraphNode("mid", midCompiled)
	if midNode.SubgraphDepth() != 2 {
		t.Errorf("mid depth: expected 2, got %d", midNode.SubgraphDepth())
	}
}

func TestSubgraphMultipleSequential(t *testing.T) {
	worker := NewStateGraph()
	worker.AddNodeFn("work", func(s StateSchema) (StateUpdate, error) {
		as := s.CloneState().(AgentState)
		as.Messages = append(as.Messages, "work_done")
		return NewStateUpdate(as), nil
	})
	worker.AddEdge(START, "work")
	worker.AddEdge("work", END)
	workerCompiled, err := worker.Compile()
	if err != nil {
		t.Fatalf("worker compile: %v", err)
	}

	parent := NewStateGraph()
	parent.AddNode("a", NewSubgraphNode("worker_a", workerCompiled))
	parent.AddNode("b", NewSubgraphNode("worker_b", workerCompiled))
	parent.AddEdge(START, "a")
	parent.AddEdge("a", "b")
	parent.AddEdge("b", END)
	parentCompiled, err := parent.Compile()
	if err != nil {
		t.Fatalf("parent compile: %v", err)
	}

	result, err := parentCompiled.Invoke(AgentState{Input: "seq", Messages: []string{}})
	if err != nil {
		t.Fatalf("invoke: %v", err)
	}

	fs := result.FinalState.(AgentState)
	count := 0
	for _, m := range fs.Messages {
		if m == "work_done" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 subgraph executions, got %d, messages: %v", count, fs.Messages)
	}
}
