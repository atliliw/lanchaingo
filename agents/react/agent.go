package react

import (
	"fmt"
	"strings"

	"github.com/atliliw/lanchaingo/agents"
	"github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/core/tools"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// ReActAgent implements the ReAct (Reasoning + Acting) pattern.
type ReActAgent struct {
	llm          language_models.BaseChatModel
	tools        []tools.Tool
	parser       *ReActOutputParser
	SystemPrompt string
}

func NewReActAgent(llm language_models.BaseChatModel, tl []tools.Tool, systemPrompt string) *ReActAgent {
	return &ReActAgent{
		llm:          llm,
		tools:        tl,
		parser:       NewReActOutputParser(),
		SystemPrompt: systemPrompt,
	}
}

func (a *ReActAgent) formatTools() string {
	var parts []string
	for _, t := range a.tools {
		parts = append(parts, fmt.Sprintf("%s: %s", t.Name(), t.Description()))
	}
	return strings.Join(parts, "\n")
}

func (a *ReActAgent) getToolNames() []string {
	names := make([]string, len(a.tools))
	for i, t := range a.tools {
		names[i] = t.Name()
	}
	return names
}

func (a *ReActAgent) buildPrompt(input string, steps []agents.AgentStep, history string) string {
	toolsDesc := a.formatTools()
	toolNames := a.getToolNames()
	scratchpad := FormatScratchpad(steps)

	prompt := BuildReActPrompt(toolsDesc, toolNames, input, scratchpad)

	if history != "" {
		prompt = fmt.Sprintf("Previous conversation history:\n%s\n\n%s", history, prompt)
	}
	if a.SystemPrompt != "" {
		prompt = fmt.Sprintf("%s\n\n%s", a.SystemPrompt, prompt)
	}
	return prompt
}

func (a *ReActAgent) Plan(steps []agents.AgentStep, inputs map[string]string) (agents.AgentOutput, error) {
	input, ok := inputs["input"]
	if !ok {
		return nil, &agents.AgentError{
			Kind:    agents.ErrOther,
			Message: "missing input parameter",
		}
	}

	history := inputs["history"]
	promptText := a.buildPrompt(input, steps, history)

	msg := messages.NewHumanMessage(promptText)
	result, err := a.llm.Chat(nil, []messages.Message{msg})
	if err != nil {
		return nil, &agents.AgentError{
			Kind:    agents.ErrOther,
			Message: fmt.Sprintf("LLM call failed: %v", err),
			Cause:   err,
		}
	}

	return a.parser.Parse(result.Content)
}

func (a *ReActAgent) InputKeys() []string {
	return []string{"input"}
}

func (a *ReActAgent) GetAllowedTools() []string {
	return a.getToolNames()
}
