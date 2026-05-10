package langgraph

import (
	"fmt"
	"sync"
)

// START and END sentinel node identifiers.
const (
	START = "__start__"
	END   = "__end__"
)

// StateGraph manages nodes, edges, and compiles into an executable graph.
type StateGraph struct {
	mu        sync.RWMutex
	nodes     map[string]GraphNode
	edges     []GraphEdge
	entry     string
	routers   map[string]ConditionalRouter
	reducer   Reducer
}

func NewStateGraph() *StateGraph {
	return &StateGraph{
		nodes:   make(map[string]GraphNode),
		edges:   make([]GraphEdge, 0),
		routers: make(map[string]ConditionalRouter),
		reducer: ReplaceReducer,
	}
}

func (g *StateGraph) AddNode(name string, node GraphNode) *StateGraph {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.nodes[name] = node
	return g
}

func (g *StateGraph) AddNodeFn(name string, fn func(StateSchema) (StateUpdate, error)) *StateGraph {
	return g.AddNode(name, NewSyncNode(name, fn))
}

func (g *StateGraph) AddEdge(source, target string) *StateGraph {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges = append(g.edges, NewFixedEdge(source, target))
	return g
}

func (g *StateGraph) AddConditionalEdges(source, routerName string, targets map[string]string, defaultTarget *string) *StateGraph {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges = append(g.edges, NewConditionalEdge(source, routerName, targets, defaultTarget))
	return g
}

func (g *StateGraph) AddFanOut(source string, targets []string) *StateGraph {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges = append(g.edges, NewFanOutEdge(source, targets))
	return g
}

func (g *StateGraph) AddFanIn(sources []string, target string) *StateGraph {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.edges = append(g.edges, NewFanInEdge(sources, target))
	return g
}

func (g *StateGraph) SetConditionalRouter(name string, router ConditionalRouter) *StateGraph {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.routers[name] = router
	return g
}

func (g *StateGraph) SetEntryPoint(node string) *StateGraph {
	g.entry = node
	return g
}

func (g *StateGraph) SetReducer(r Reducer) *StateGraph {
	g.reducer = r
	return g
}

func (g *StateGraph) findEntry() (string, error) {
	if g.entry != "" {
		return g.entry, nil
	}
	for _, e := range g.edges {
		if e.Source == START && e.Type == EdgeFixed {
			return e.Target, nil
		}
	}
	return "", fmt.Errorf("no entry point found")
}

func (g *StateGraph) Compile() (*CompiledGraph, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if len(g.nodes) == 0 {
		return nil, NewGraphError(ErrValidation, "graph has no nodes", nil)
	}

	entry, err := g.findEntry()
	if err != nil {
		return nil, NewGraphError(ErrValidation, err.Error(), nil)
	}

	if _, ok := g.nodes[entry]; !ok && entry != START {
		return nil, NewGraphError(ErrValidation, fmt.Sprintf("entry point '%s' not found", entry), nil)
	}

	cg := &CompiledGraph{
		nodes:   copyNodes(g.nodes),
		edges:   copyEdges(g.edges),
		entry:   entry,
		reducer: g.reducer,
		routers: copyRouters(g.routers),
		recursionLimit: 25,
	}

	if err := cg.validate(); err != nil {
		return nil, err
	}

	return cg, nil
}

func copyNodes(src map[string]GraphNode) map[string]GraphNode {
	dst := make(map[string]GraphNode, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

func copyEdges(src []GraphEdge) []GraphEdge {
	dst := make([]GraphEdge, len(src))
	copy(dst, src)
	return dst
}

func copyRouters(src map[string]ConditionalRouter) map[string]ConditionalRouter {
	dst := make(map[string]ConditionalRouter, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
