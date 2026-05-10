package langgraph

// EdgeTarget specifies where an edge leads.
type EdgeTarget struct {
	Type     EdgeType
	Target   string // for fixed edges
	Router   string // for conditional edges
}

type EdgeType int

const (
	EdgeFixed EdgeType = iota
	EdgeConditional
	EdgeFanOut
	EdgeFanIn
)

// GraphEdge defines a transition between nodes.
type GraphEdge struct {
	Source        string
	Target        string // fixed target or fan-in target
	RouterName    string // conditional router name
	RouteTargets  map[string]string
	DefaultTarget *string
	FanOutTargets []string
	FanInSources  []string
	Type          EdgeType
}

func NewFixedEdge(source, target string) GraphEdge {
	return GraphEdge{
		Source: source,
		Target: target,
		Type:   EdgeFixed,
	}
}

func NewConditionalEdge(source, routerName string, targets map[string]string, defaultTarget *string) GraphEdge {
	return GraphEdge{
		Source:        source,
		RouterName:    routerName,
		RouteTargets:  targets,
		DefaultTarget: defaultTarget,
		Type:          EdgeConditional,
	}
}

func NewFanOutEdge(source string, targets []string) GraphEdge {
	return GraphEdge{
		Source:        source,
		FanOutTargets: targets,
		Type:          EdgeFanOut,
	}
}

func NewFanInEdge(sources []string, target string) GraphEdge {
	return GraphEdge{
		Source:       "__fanin__",
		FanInSources: sources,
		Target:       target,
		Type:         EdgeFanIn,
	}
}

// ConditionalRouter routes to the next node based on state.
type ConditionalRouter func(state StateSchema) string
