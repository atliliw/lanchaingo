package langgraph

import (
	"fmt"
	"sync"
)

// SubgraphNode embeds a compiled graph as a node in a parent graph.
// It supports state mapping via InputMap and OutputMap for flexible
// parent-child state transformations.
//
// SubgraphNode implements GraphNode and can be used anywhere a regular
// node is expected: linear edges, conditional branches, fan-out points.
type SubgraphNode struct {
	name     string
	subgraph *CompiledGraph

	// InputMap transforms parent state into subgraph input state.
	// When non-nil, this is called before subgraph execution, allowing
	// the parent to pass a subset of state or a differently-typed state
	// to the subgraph.
	// When nil, the parent state is passed through directly, requiring
	// the parent and subgraph to share the same StateSchema type.
	inputMap func(StateSchema) StateSchema

	// OutputMap transforms subgraph output state back into parent state.
	// The function receives (subgraphFinalState, parentOriginalState) and
	// should return the merged parent state.
	// When nil, the subgraph's final state replaces the parent state
	// directly (requires same StateSchema type).
	outputMap func(StateSchema, StateSchema) StateSchema

	mu sync.Mutex
}

// SubgraphOption configures a SubgraphNode.
type SubgraphOption func(*SubgraphNode)

// WithInputMap sets a function that transforms the parent's state into the
// subgraph's input state. This enables type-safe state mapping between
// parent and child graphs with different state schemas.
func WithInputMap(fn func(StateSchema) StateSchema) SubgraphOption {
	return func(n *SubgraphNode) {
		n.inputMap = fn
	}
}

// WithOutputMap sets a function that transforms the subgraph's final state
// back into the parent's state. The function receives the subgraph's final
// state and the parent's original state (before subgraph execution), and
// should return the merged parent state.
func WithOutputMap(fn func(StateSchema, StateSchema) StateSchema) SubgraphOption {
	return func(n *SubgraphNode) {
		n.outputMap = fn
	}
}

// NewSubgraphNode creates a new SubgraphNode.
// The subgraph must be a CompiledGraph ready for execution.
// Nil subgraph panics.
func NewSubgraphNode(name string, subgraph *CompiledGraph, opts ...SubgraphOption) *SubgraphNode {
	if subgraph == nil {
		panic("langgraph: subgraph must not be nil")
	}
	n := &SubgraphNode{
		name:     name,
		subgraph: subgraph,
	}
	for _, opt := range opts {
		opt(n)
	}
	return n
}

func (n *SubgraphNode) Name() string { return n.name }

// SubgraphDepth returns the nesting depth of the subgraph tree rooted at
// this node. A leaf subgraph (containing no further SubgraphNodes) has
// depth 1. Each level of nesting adds 1.
func (n *SubgraphNode) SubgraphDepth() int {
	maxDepth := 0
	for _, node := range n.subgraph.nodes {
		if sn, ok := node.(*SubgraphNode); ok {
			if d := sn.SubgraphDepth(); d > maxDepth {
				maxDepth = d
			}
		}
	}
	return maxDepth + 1
}

// Execute runs the embedded subgraph with the given state.
//
// State mapping:
//  1. If InputMap is set, parent state is transformed before passing to subgraph.
//  2. The subgraph is invoked with the (possibly transformed) state.
//  3. If the subgraph completes successfully, OutputMap (if set) transforms
//     the result back into parent state, then the state is merged via the
//     parent graph's reducer.
//  4. If the subgraph is interrupted, a SubgraphInterruptError is returned
//     containing the partial state. The parent graph's execution loop
//     recognizes this error and converts it to an interrupt result.
func (n *SubgraphNode) Execute(state StateSchema) (StateUpdate, error) {
	n.mu.Lock()
	defer n.mu.Unlock()

	input := state
	if n.inputMap != nil {
		input = n.inputMap(state)
	}

	result, err := n.subgraph.Invoke(input)
	if err != nil {
		return StateUpdate{}, fmt.Errorf("subgraph '%s' failed: %w", n.name, err)
	}

	if result.Interrupted {
		partialState := result.FinalState
		if partialState != nil && n.outputMap != nil {
			partialState = n.outputMap(partialState, state)
		}
		return StateUpdate{}, &SubgraphInterruptError{
			SubgraphName:  n.name,
			InterruptedAt: result.InterruptedAt,
			PartialState:  partialState,
		}
	}

	output := result.FinalState
	if n.outputMap != nil {
		output = n.outputMap(result.FinalState, state)
	}

	return NewStateUpdate(output), nil
}
