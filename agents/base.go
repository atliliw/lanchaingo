package agents

import (
	"fmt"
	"sync"

	"github.com/atliliw/lanchaingo/core/tools"
	"github.com/atliliw/lanchaingo/memory"
)

// BaseAgent defines the interface for agent planning.
type BaseAgent interface {
	Plan(intermediateSteps []AgentStep, inputs map[string]string) (AgentOutput, error)
	InputKeys() []string
	GetAllowedTools() []string
}

// AgentExecutor runs the agent loop: Plan 鈫?Act 鈫?Observe.
type AgentExecutor struct {
	Agent         BaseAgent
	Tools         []tools.Tool
	MaxIterations int
	Verbose       bool
	Memory        *memory.ConversationBufferMemory
	mu            sync.Mutex
}

func NewAgentExecutor(agent BaseAgent, tl []tools.Tool) *AgentExecutor {
	return &AgentExecutor{
		Agent:         agent,
		Tools:         tl,
		MaxIterations: 10,
	}
}

func (e *AgentExecutor) WithMaxIterations(n int) *AgentExecutor {
	e.MaxIterations = n
	return e
}

func (e *AgentExecutor) WithVerbose(v bool) *AgentExecutor {
	e.Verbose = v
	return e
}

func (e *AgentExecutor) WithMemory(m *memory.ConversationBufferMemory) *AgentExecutor {
	e.Memory = m
	return e
}

func (e *AgentExecutor) Invoke(input string) (string, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	inputs := map[string]string{"input": input}

	if e.Memory != nil {
		vars, err := e.Memory.LoadMemoryVariables(inputs)
		if err == nil {
			if h, ok := vars["history"]; ok {
				if hs, ok := h.(string); ok && hs != "" {
					inputs["history"] = hs
				}
			}
		}
	}

	var steps []AgentStep
	var result string

	for iter := 0; iter < e.MaxIterations; iter++ {
		if e.Verbose {
			fmt.Printf("\n=== Iteration %d ===\n", iter+1)
		}

		output, err := e.Agent.Plan(steps, inputs)
		if err != nil {
			return "", err
		}

		if finish, ok := output.(*AgentFinishOutput); ok {
			result = finish.Finish.Output
			if e.Verbose {
				fmt.Printf("Final answer: %s\n", result)
			}
			break
		}

		if action, ok := output.(*AgentActionOutput); ok {
			if e.Verbose {
				fmt.Printf("Action: %s(%s)\n", action.Action.Tool, action.Action.ToolInput)
			}

			observation := e.executeTool(action.Action)
			if e.Verbose {
				fmt.Printf("Observation: %s\n", observation)
			}

			steps = append(steps, AgentStep{
				Action:      action.Action,
				Observation: observation,
			})
		}
	}

	if result == "" {
		result = "Agent stopped due to iteration limit"
	}

	if e.Memory != nil {
		_ = e.Memory.SaveContext(inputs, map[string]string{"output": result})
	}

	return result, nil
}

func (e *AgentExecutor) executeTool(action AgentAction) string {
	for _, t := range e.Tools {
		if t.Name() == action.Tool {
			res, err := t.Run(nil, action.ToolInput.String())
			if err != nil {
				return fmt.Sprintf("Error: %v", err)
			}
			return res
		}
	}
	return fmt.Sprintf("Tool %s not found", action.Tool)
}
