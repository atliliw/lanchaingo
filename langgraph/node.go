package langgraph

// GraphNode is the execution unit in a graph.
type GraphNode interface {
	Name() string
	Execute(state StateSchema) (StateUpdate, error)
}

// SyncNode wraps a synchronous function as a GraphNode.
type SyncNode struct {
	name string
	fn   func(StateSchema) (StateUpdate, error)
}

func NewSyncNode(name string, fn func(StateSchema) (StateUpdate, error)) *SyncNode {
	return &SyncNode{name: name, fn: fn}
}

func (n *SyncNode) Name() string { return n.name }

func (n *SyncNode) Execute(state StateSchema) (StateUpdate, error) {
	return n.fn(state)
}
