package langgraph

// SubgraphNode embeds a compiled graph as a node in a parent graph.
type SubgraphNode struct {
	name     string
	subgraph *CompiledGraph
}

func NewSubgraphNode(name string, subgraph *CompiledGraph) *SubgraphNode {
	return &SubgraphNode{name: name, subgraph: subgraph}
}

func (n *SubgraphNode) Name() string { return n.name }

func (n *SubgraphNode) Execute(state StateSchema) (StateUpdate, error) {
	result, err := n.subgraph.Invoke(state)
	if err != nil {
		return StateUpdate{}, NewGraphError(ErrExecution, "subgraph execution failed: "+err.Error(), err)
	}
	return NewStateUpdate(result.FinalState), nil
}
