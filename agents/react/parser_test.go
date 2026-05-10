package react

import (
	"testing"

	"github.com/atliliw/lanchaingo/agents"
)

// 娴嬭瘯瑙ｆ瀽 Action 鏍煎紡杈撳嚭
func TestParseAction(t *testing.T) {
	parser := NewReActOutputParser()

	text := `Thought: I need to calculate this
Action: calculator
Action Input: {"expression": "2 + 3"}`

	result, err := parser.Parse(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	action, ok := result.(*agents.AgentActionOutput)
	if !ok {
		t.Fatal("expected AgentActionOutput")
	}
	if action.Action.Tool != "calculator" {
		t.Errorf("expected calculator, got %s", action.Action.Tool)
	}
}

// 娴嬭瘯瑙ｆ瀽 Final Answer 鏍煎紡
func TestParseFinalAnswer(t *testing.T) {
	parser := NewReActOutputParser()

	text := `Thought: I know the answer
Final Answer: The answer is 42`

	result, err := parser.Parse(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	finish, ok := result.(*agents.AgentFinishOutput)
	if !ok {
		t.Fatal("expected AgentFinishOutput")
	}
	if finish.Finish.Output != "The answer is 42" {
		t.Errorf("expected 'The answer is 42', got '%s'", finish.Finish.Output)
	}
}

// 娴嬭瘯瑙ｆ瀽瀛楃涓插伐鍏疯緭鍏?func TestParseStringInput(t *testing.T) {
	parser := NewReActOutputParser()

	text := `Thought: Need weather
Action: weather
Action Input: Beijing`

	result, err := parser.Parse(text)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	action := result.(*agents.AgentActionOutput)
	if action.Action.ToolInput.String() != "Beijing" {
		t.Errorf("expected Beijing, got %s", action.Action.ToolInput.String())
	}
}

// 娴嬭瘯鏃犳硶瑙ｆ瀽鐨勮緭鍑?func TestParseInvalidOutput(t *testing.T) {
	parser := NewReActOutputParser()

	_, err := parser.Parse("This is invalid output")
	if err == nil {
		t.Error("expected error for invalid output")
	}
}
