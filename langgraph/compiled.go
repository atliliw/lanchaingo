package langgraph

import (
	"fmt"
	"sort"
)

// CompiledGraph is a ready-to-execute graph.
type CompiledGraph struct {
	nodes          map[string]GraphNode
	edges          []GraphEdge
	entry          string
	reducer        Reducer
	routers        map[string]ConditionalRouter
	recursionLimit int
}

func (cg *CompiledGraph) WithRecursionLimit(limit int) *CompiledGraph {
	cg.recursionLimit = limit
	return cg
}

func (cg *CompiledGraph) validate() error {
	// Collect all reachable nodes
	reachable := cg.computeReachable()
	endReachable := cg.computeEndReachable()

	for name := range cg.nodes {
		if !reachable[name] {
			return NewGraphError(ErrOrphanNode, fmt.Sprintf("unreachable node: %s", name), nil)
		}
		if !endReachable[name] {
			return NewGraphError(ErrInfiniteCycle, fmt.Sprintf("node '%s' has no path to END", name), nil)
		}
	}

	return nil
}

func (cg *CompiledGraph) computeReachable() map[string]bool {
	reachable := make(map[string]bool)
	queue := []string{cg.entry}
	visited := make(map[string]bool)

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] || current == END {
			continue
		}
		visited[current] = true
		reachable[current] = true

		for _, e := range cg.edges {
			if e.Source == current {
				switch e.Type {
				case EdgeFixed:
					if e.Target != END && !visited[e.Target] {
						queue = append(queue, e.Target)
					}
				case EdgeConditional:
					for _, t := range e.RouteTargets {
						if t != END && !visited[t] {
							queue = append(queue, t)
						}
					}
					if e.DefaultTarget != nil && *e.DefaultTarget != END && !visited[*e.DefaultTarget] {
						queue = append(queue, *e.DefaultTarget)
					}
				case EdgeFanOut:
					for _, t := range e.FanOutTargets {
						if t != END && !visited[t] {
							queue = append(queue, t)
						}
					}
				case EdgeFanIn:
					if e.Target != END && !visited[e.Target] {
						queue = append(queue, e.Target)
					}
				}
			}
		}
	}
	return reachable
}

func (cg *CompiledGraph) computeEndReachable() map[string]bool {
	reachable := map[string]bool{END: true}
	changed := true

	for changed {
		changed = false
		for _, e := range cg.edges {
			switch e.Type {
			case EdgeFixed:
				if reachable[e.Target] && !reachable[e.Source] {
					reachable[e.Source] = true
					changed = true
				}
			case EdgeConditional:
				anyTarget := false
				for _, t := range e.RouteTargets {
					if reachable[t] {
						anyTarget = true
						break
					}
				}
				if !anyTarget && e.DefaultTarget != nil {
					anyTarget = reachable[*e.DefaultTarget]
				}
				if anyTarget && !reachable[e.Source] {
					reachable[e.Source] = true
					changed = true
				}
			case EdgeFanOut:
				all := len(e.FanOutTargets) > 0
				for _, t := range e.FanOutTargets {
					if !reachable[t] {
						all = false
						break
					}
				}
				if all && !reachable[e.Source] {
					reachable[e.Source] = true
					changed = true
				}
			case EdgeFanIn:
				if reachable[e.Target] {
					for _, s := range e.FanInSources {
						if !reachable[s] {
							reachable[s] = true
							changed = true
						}
					}
				}
			}
		}
	}
	return reachable
}

func (cg *CompiledGraph) findNext(current string, state StateSchema) (string, error) {
	for _, e := range cg.edges {
		if e.Source == current {
			switch e.Type {
			case EdgeFixed:
				return e.Target, nil
			case EdgeConditional:
				router, ok := cg.routers[e.RouterName]
				if !ok {
					return "", NewGraphError(ErrRouting, fmt.Sprintf("router '%s' not found", e.RouterName), nil)
				}
				route := router(state)
				if target, ok := e.RouteTargets[route]; ok {
					return target, nil
				}
				if e.DefaultTarget != nil {
					return *e.DefaultTarget, nil
				}
				return "", NewGraphError(ErrRouting, fmt.Sprintf("no target for route '%s'", route), nil)
			case EdgeFanOut:
				if len(e.FanOutTargets) > 0 {
					return e.FanOutTargets[0], nil
				}
				return END, nil
			case EdgeFanIn:
				// Skip fan-in edges when looking for next
				continue
			}
		}
	}

	if current == cg.entry && len(cg.nodes) == 1 {
		return END, nil
	}

	return "", NewGraphError(ErrRouting, fmt.Sprintf("no outgoing edge from '%s'", current), nil)
}

// GraphInvocation is the result of graph execution.
type GraphInvocation struct {
	FinalState     StateSchema
	Steps          []ExecutionStep
	RecursionCount int
}

// ExecutionStep records a single step in graph execution.
type ExecutionStep struct {
	Name     string
	Metadata map[string]any
}

// Invoke executes the graph from the entry point.
func (cg *CompiledGraph) Invoke(input StateSchema) (*GraphInvocation, error) {
	state := input.CloneState()
	current := cg.entry
	recursions := 0

	for current != END {
		if recursions >= cg.recursionLimit {
			return nil, NewGraphError(ErrRecursionLimit, fmt.Sprintf("recursion limit %d reached", cg.recursionLimit), nil)
		}
		recursions++

		node, ok := cg.nodes[current]
		if !ok {
			return nil, NewGraphError(ErrExecution, fmt.Sprintf("node '%s' not found", current), nil)
		}

		update, err := node.Execute(state)
		if err != nil {
			return nil, NewGraphError(ErrExecution, fmt.Sprintf("node '%s' failed: %v", current, err), err)
		}

		if update.Update != nil {
			state = cg.reducer(state, update.Update)
		}

		next, err := cg.findNext(current, state)
		if err != nil {
			return nil, err
		}
		current = next
	}

	return &GraphInvocation{
		FinalState:     state,
		Steps:          nil,
		RecursionCount: recursions,
	}, nil
}

// NodeNames returns the names of all nodes.
func (cg *CompiledGraph) NodeNames() []string {
	names := make([]string, 0, len(cg.nodes))
	for n := range cg.nodes {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// EntryPoint returns the entry point node name.
func (cg *CompiledGraph) EntryPoint() string { return cg.entry }

// RecursionLimit returns the recursion limit.
func (cg *CompiledGraph) RecursionLimit() int { return cg.recursionLimit }
