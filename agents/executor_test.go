package agents

import (
	"context"
	"testing"

	coretools "github.com/atliliw/lanchaingo/core/tools"
)

// mockAgentForTest 鏄竴涓畝鍗曠殑 BaseAgent 瀹炵幇锛岀敤浜庢祴璇?Executor
type mockAgentForTest struct {
	responses []AgentOutput
	index     int
}

func (a *mockAgentForTest) Plan(_ []AgentStep, _ map[string]string) (AgentOutput, error) {
	if a.index >= len(a.responses) {
		return &AgentFinishOutput{Finish: AgentFinish{Output: "done"}}, nil
	}
	out := a.responses[a.index]
	a.index++
	return out, nil
}

func (a *mockAgentForTest) InputKeys() []string       { return []string{"input"} }
func (a *mockAgentForTest) GetAllowedTools() []string { return nil }

// mockToolForTest 鏄竴涓畝鍗曠殑 coretools.Tool 瀹炵幇
type mockToolForTest struct{}

func (m *mockToolForTest) Name() string { return "test_tool" }
func (m *mockToolForTest) Description() string { return "test" }
func (m *mockToolForTest) Run(_ context.Context, input string) (string, error) { return "observed: " + input, nil }

// 娴嬭瘯 AgentExecutor 鎵ц娴佺▼
func TestAgentExecutorBasic(t *testing.T) {
	agent := &mockAgentForTest{
		responses: []AgentOutput{
			&AgentFinishOutput{Finish: AgentFinish{Output: "answer"}},
		},
	}
	exec := NewAgentExecutor(agent, nil)

	result, err := exec.Invoke("question")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "answer" {
		t.Errorf("expected answer, got %s", result)
	}
}

// 娴嬭瘯 AgentExecutor 澶氭鎵ц锛欰ction 鈫?Observation 鈫?Finish
func TestAgentExecutorWithTool(t *testing.T) {
	agent := &mockAgentForTest{
		responses: []AgentOutput{
			&AgentActionOutput{Action: AgentAction{
				Tool: "test_tool", ToolInput: StringToolInput("hello"),
			}},
			&AgentFinishOutput{Finish: AgentFinish{Output: "done"}},
		},
	}
	tl := []coretools.Tool{&mockToolForTest{}}
	exec := NewAgentExecutor(agent, tl)

	result, err := exec.Invoke("do it")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "done" {
		t.Errorf("expected done, got %s", result)
	}
}

// 娴嬭瘯 max iterations 闄愬埗
func TestAgentExecutorMaxIterations(t *testing.T) {
	agent := &mockAgentForTest{
		responses: []AgentOutput{
			&AgentActionOutput{Action: AgentAction{Tool: "test_tool", ToolInput: StringToolInput("x")}},
			&AgentActionOutput{Action: AgentAction{Tool: "test_tool", ToolInput: StringToolInput("x")}},
			&AgentActionOutput{Action: AgentAction{Tool: "test_tool", ToolInput: StringToolInput("x")}},
		},
	}
	tl := []coretools.Tool{&mockToolForTest{}}
	exec := NewAgentExecutor(agent, tl).WithMaxIterations(2)

	result, err := exec.Invoke("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Agent stopped due to iteration limit" {
		t.Errorf("expected stop message, got %s", result)
	}
}

// 娴嬭瘯 tool not found 鍦烘櫙
func TestAgentExecutorToolNotFound(t *testing.T) {
	agent := &mockAgentForTest{
		responses: []AgentOutput{
			&AgentActionOutput{Action: AgentAction{Tool: "nonexistent", ToolInput: StringToolInput("x")}},
			&AgentFinishOutput{Finish: AgentFinish{Output: "recovered"}},
		},
	}
	exec := NewAgentExecutor(agent, nil)

	result, err := exec.Invoke("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "recovered" {
		t.Errorf("expected recovered, got %s", result)
	}
}
