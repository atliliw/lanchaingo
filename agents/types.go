package agents

import "encoding/json"

// ToolInput represents input to a tool, either a string or JSON object.
type ToolInput struct {
	Str    string
	Object map[string]any
	IsJSON bool
}

func StringToolInput(s string) ToolInput {
	return ToolInput{Str: s, IsJSON: false}
}

func JSONToolInput(obj map[string]any) ToolInput {
	return ToolInput{Object: obj, IsJSON: true}
}

func (ti ToolInput) String() string {
	if ti.IsJSON {
		b, _ := json.Marshal(ti.Object)
		return string(b)
	}
	return ti.Str
}

// AgentAction represents a decision to execute a tool.
type AgentAction struct {
	Tool      string    `json:"tool"`
	ToolInput ToolInput `json:"tool_input"`
	Log       string    `json:"log"`
}

// AgentFinish represents the final answer from an agent.
type AgentFinish struct {
	Output string `json:"output"`
	Log    string `json:"log"`
}

// AgentStep is one executed step in the agent loop.
type AgentStep struct {
	Action      AgentAction `json:"action"`
	Observation string      `json:"observation"`
}

// AgentOutput is what an agent's plan method returns.
type AgentOutput interface {
	IsFinish() bool
	IsAction() bool
}

type AgentActionOutput struct {
	Action AgentAction
}

func (o *AgentActionOutput) IsFinish() bool { return false }
func (o *AgentActionOutput) IsAction() bool { return true }

type AgentFinishOutput struct {
	Finish AgentFinish
}

func (o *AgentFinishOutput) IsFinish() bool { return true }
func (o *AgentFinishOutput) IsAction() bool { return false }
