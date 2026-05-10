package langgraph

// StateSchema is implemented by state types that flow through the graph.
type StateSchema interface {
	CloneState() StateSchema
}

// StateUpdate carries a state change from a node execution.
type StateUpdate struct {
	Update   StateSchema
	Metadata map[string]any
}

func NewStateUpdate(state StateSchema) StateUpdate {
	return StateUpdate{Update: state, Metadata: make(map[string]any)}
}

func NoStateUpdate() StateUpdate {
	return StateUpdate{Metadata: make(map[string]any)}
}

// Reducer defines how state updates are merged into current state.
type Reducer func(current, update StateSchema) StateSchema

// ReplaceReducer replaces the entire state with the update.
func ReplaceReducer(current, update StateSchema) StateSchema {
	return update
}

// AgentState is a simple default state implementation.
type AgentState struct {
	Input    string
	Output   *string
	Messages []string
	Steps    int
}

func (a AgentState) CloneState() StateSchema {
	msgCopy := make([]string, len(a.Messages))
	copy(msgCopy, a.Messages)
	var out *string
	if a.Output != nil {
		s := *a.Output
		out = &s
	}
	return AgentState{
		Input:    a.Input,
		Output:   out,
		Messages: msgCopy,
		Steps:    a.Steps,
	}
}
