package react

import (
	"regexp"

	"github.com/atliliw/lanchaingo/agents"
)

// ReActOutputParser parses LLM output in ReAct format.
type ReActOutputParser struct {
	actionRegex        *regexp.Regexp
	finalAnswerMarker string
}

func NewReActOutputParser() *ReActOutputParser {
	return &ReActOutputParser{
		actionRegex:        regexp.MustCompile(`Action\s*:\s*(.*?)\s*\nAction\s*Input\s*:\s*(.*?)(?:\n|$)`),
		finalAnswerMarker: "Final Answer:",
	}
}

// Parse interprets LLM output as either an action or final answer.
func (p *ReActOutputParser) Parse(text string) (agents.AgentOutput, error) {
	text = trimSpace(text)

	if contains(text, p.finalAnswerMarker) {
		return p.parseFinalAnswer(text)
	}

	if action, err := p.parseAction(text); err != nil {
		return nil, err
	} else if action != nil {
		return action, nil
	}

	return nil, &agents.AgentError{
		Kind:    agents.ErrOutputParsing,
		Message: "Could not parse LLM output. Expected ReAct format with Action/Action Input or Final Answer.",
	}
}

func (p *ReActOutputParser) parseFinalAnswer(text string) (agents.AgentOutput, error) {
	parts := splitAfter(text, p.finalAnswerMarker)
	if len(parts) < 2 {
		return nil, &agents.AgentError{
			Kind:    agents.ErrOutputParsing,
			Message: "Final Answer marker found but no content after it",
		}
	}
	answer := trimSpace(parts[len(parts)-1])
	return &agents.AgentFinishOutput{
		Finish: agents.AgentFinish{Output: answer, Log: text},
	}, nil
}

func (p *ReActOutputParser) parseAction(text string) (agents.AgentOutput, error) {
	matches := p.actionRegex.FindStringSubmatch(text)
	if len(matches) < 3 {
		return nil, nil
	}

	tool := trimSpace(matches[1])
	toolInput := trimSpace(matches[2])
	ti := agents.StringToolInput(toolInput)

	return &agents.AgentActionOutput{
		Action: agents.AgentAction{
			Tool:      tool,
			ToolInput: ti,
			Log:       text,
		},
	}, nil
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}
	return s[start:end]
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func splitAfter(s, sep string) []string {
	var result []string
	for {
		idx := indexOf(s, sep)
		if idx < 0 {
			result = append(result, s)
			break
		}
		result = append(result, s[:idx+len(sep)])
		s = s[idx+len(sep):]
	}
	return result
}

func indexOf(s, sep string) int {
	for i := 0; i <= len(s)-len(sep); i++ {
		if s[i:i+len(sep)] == sep {
			return i
		}
	}
	return -1
}
