package agents

import "fmt"

// AgentErrorKind categorizes agent errors.
type AgentErrorKind int

const (
	ErrOutputParsing  AgentErrorKind = iota
	ErrToolNotFound
	ErrToolExecution
	ErrMaxIterations
	ErrOther
)

// AgentError represents errors during agent execution.
type AgentError struct {
	Kind    AgentErrorKind
	Message string
	Cause   error
}

func (e *AgentError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("agent: %s: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("agent: %s", e.Message)
}

func (e *AgentError) Unwrap() error { return e.Cause }

func NewAgentError(kind AgentErrorKind, msg string, cause error) *AgentError {
	return &AgentError{Kind: kind, Message: msg, Cause: cause}
}
