package function_calling

import (
	"fmt"

	"github.com/atliliw/lanchaingo/agents"
	"github.com/atliliw/lanchaingo/core/language_models"
	"github.com/atliliw/lanchaingo/core/tools"
	"github.com/atliliw/lanchaingo/schema/messages"
)

// FunctionCallingAgent uses native LLM function calling (tool calling).
type FunctionCallingAgent struct {
	llm          language_models.BaseChatModel
	tools        []tools.Tool
	SystemPrompt string
}

func NewFunctionCallingAgent(llm language_models.BaseChatModel, tl []tools.Tool, systemPrompt string) *FunctionCallingAgent {
	defs := make([]language_models.ToolDefinition, len(tl))
	for i, t := range tl {
		defs[i] = language_models.ToolDefinition{
			Name:        t.Name(),
			Description: t.Description(),
			Parameters:  make(map[string]any),
		}
	}
	if len(defs) > 0 {
		llm.BindTools(defs)
	}
	return &FunctionCallingAgent{
		llm:          llm,
		tools:        tl,
		SystemPrompt: systemPrompt,
	}
}

func (a *FunctionCallingAgent) buildMessages(inputs map[string]string, steps []agents.AgentStep) []messages.Message {
	var msgs []messages.Message

	sys := a.SystemPrompt
	if sys == "" {
		sys = "You are a helpful assistant that can use tools to answer questions."
	}
	msgs = append(msgs, messages.NewSystemMessage(sys))
	msgs = append(msgs, messages.NewHumanMessage(inputs["input"]))

	for _, step := range steps {
		tc := language_models.ToolCall{
			ID:   step.Action.Log,
			Name: step.Action.Tool,
		}
		aiMsg := messages.Message{
			Content:     "",
			MessageType: messages.MessageTypeAI,
			ToolCalls:   []language_models.ToolCall{tc},
		}
		msgs = append(msgs, aiMsg)
		msgs = append(msgs, messages.NewToolMessage(step.Observation, step.Action.Log))
	}

	return msgs
}

func (a *FunctionCallingAgent) Plan(steps []agents.AgentStep, inputs map[string]string) (agents.AgentOutput, error) {
	msgs := a.buildMessages(inputs, steps)

	result, err := a.llm.Chat(nil, msgs)
	if err != nil {
		return nil, &agents.AgentError{
			Kind:    agents.ErrOther,
			Message: fmt.Sprintf("LLM call failed: %v", err),
			Cause:   err,
		}
	}

	if len(result.ToolCalls) > 0 {
		tc := result.ToolCalls[0]
		return &agents.AgentActionOutput{
			Action: agents.AgentAction{
				Tool:      tc.Name,
				ToolInput: agents.StringToolInput(tc.Arguments),
				Log:       tc.ID,
			},
		}, nil
	}

	return &agents.AgentFinishOutput{
		Finish: agents.AgentFinish{Output: result.Content},
	}, nil
}

func (a *FunctionCallingAgent) InputKeys() []string {
	return []string{"input"}
}

func (a *FunctionCallingAgent) GetAllowedTools() []string {
	names := make([]string, len(a.tools))
	for i, t := range a.tools {
		names[i] = t.Name()
	}
	return names
}
