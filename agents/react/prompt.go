package react

import (
	"fmt"
	"strings"

	"github.com/atliliw/lanchaingo/agents"
)

const reactPrefix = `Answer the following questions as best you can. You have access to the following tools:

{tools}

Use the following format:

Question: the input question you must answer
Thought: you should always think about what to do
Action: the action to take, should be one of [{tool_names}]
Action Input: the input to the action
Observation: the result of the action
... (this Thought/Action/Action Input/Observation can repeat N times)
Thought: I now know the final answer
Final Answer: the final answer to the original input question

Begin!

Question: {input}
Thought:{agent_scratchpad}`

// BuildReActPrompt constructs the full ReAct prompt.
func BuildReActPrompt(toolsDescription string, toolNames []string, input string, scratchpad string) string {
	p := strings.ReplaceAll(reactPrefix, "{tools}", toolsDescription)
	p = strings.ReplaceAll(p, "{tool_names}", strings.Join(toolNames, ", "))
	p = strings.ReplaceAll(p, "{input}", input)
	p = strings.ReplaceAll(p, "{agent_scratchpad}", scratchpad)
	return p
}

// FormatScratchpad formats intermediate steps into the agent scratchpad.
func FormatScratchpad(steps []agents.AgentStep) string {
	var parts []string
	for _, step := range steps {
		logLine := ""
		if step.Action.Log != "" {
			logLine = strings.Split(step.Action.Log, "\n")[0]
		}
		parts = append(parts, fmt.Sprintf(
			" %s\nAction: %s\nAction Input: %s\nObservation: %s\n",
			logLine,
			step.Action.Tool,
			step.Action.ToolInput.String(),
			step.Observation,
		))
	}
	return strings.Join(parts, "")
}
