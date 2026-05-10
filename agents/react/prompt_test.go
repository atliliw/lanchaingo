package react

import (
	"testing"

	"github.com/atliliw/lanchaingo/agents"
)

// 娴嬭瘯 BuildReActPrompt 鍩烘湰鍔熻兘
func TestBuildReActPrompt(t *testing.T) {
	prompt := BuildReActPrompt(
		"calculator: 璁＄畻鏁板琛ㄨ揪寮?,
		[]string{"calculator"},
		"璁＄畻 2 + 2",
		"",
	)

	if !containsStr(prompt, "calculator: 璁＄畻鏁板琛ㄨ揪寮?) {
		t.Error("prompt should contain tool description")
	}
	if !containsStr(prompt, "璁＄畻 2 + 2") {
		t.Error("prompt should contain input")
	}
	if !containsStr(prompt, "[calculator]") {
		t.Error("prompt should contain tool names in brackets")
	}
}

// 娴嬭瘯 FormatScratchpad
func TestFormatScratchpad(t *testing.T) {
	steps := []agents.AgentStep{
		{
			Action: agents.AgentAction{
				Tool:      "calculator",
				ToolInput: agents.StringToolInput("2 + 2"),
				Log:       "I need to calculate",
			},
			Observation: "4",
		},
	}

	scratchpad := FormatScratchpad(steps)

	if !containsStr(scratchpad, "calculator") {
		t.Error("scratchpad should contain tool name")
	}
	if !containsStr(scratchpad, "4") {
		t.Error("scratchpad should contain observation")
	}
}
