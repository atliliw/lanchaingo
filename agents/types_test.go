package agents

import "testing"

// 娴嬭瘯 ToolInput 瀛楃涓叉ā寮?func TestStringToolInput(t *testing.T) {
	ti := StringToolInput("hello")
	if ti.String() != "hello" {
		t.Errorf("expected hello, got %s", ti.String())
	}
	if ti.IsJSON {
		t.Error("expected IsJSON false")
	}
}

// 娴嬭瘯 ToolInput JSON 妯″紡
func TestJSONToolInput(t *testing.T) {
	ti := JSONToolInput(map[string]any{"key": "val"})
	if !ti.IsJSON {
		t.Error("expected IsJSON true")
	}
}

// 娴嬭瘯 AgentAction 鍒涘缓
func TestAgentActionCreation(t *testing.T) {
	action := AgentAction{
		Tool:      "calculator",
		ToolInput: StringToolInput("2+2"),
		Log:       "thought...",
	}
	if action.Tool != "calculator" || action.ToolInput.String() != "2+2" {
		t.Errorf("unexpected action: %+v", action)
	}
}

// 娴嬭瘯 AgentFinish 鍒涘缓
func TestAgentFinishCreation(t *testing.T) {
	f := AgentFinish{Output: "42", Log: "done"}
	if f.Output != "42" {
		t.Errorf("expected 42, got %s", f.Output)
	}
}

// 娴嬭瘯 AgentOutput 鎺ュ彛
func TestAgentOutputTypes(t *testing.T) {
	action := &AgentActionOutput{Action: AgentAction{Tool: "calc"}}
	if !action.IsAction() || action.IsFinish() {
		t.Error("Action output type check failed")
	}

	finish := &AgentFinishOutput{Finish: AgentFinish{Output: "done"}}
	if !finish.IsFinish() || finish.IsAction() {
		t.Error("Finish output type check failed")
	}
}

// 娴嬭瘯 AgentStep
func TestAgentStep(t *testing.T) {
	step := AgentStep{
		Action:      AgentAction{Tool: "search", ToolInput: StringToolInput("query")},
		Observation: "result",
	}
	if step.Action.Tool != "search" || step.Observation != "result" {
		t.Errorf("unexpected step: %+v", step)
	}
}
