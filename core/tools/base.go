package tools

import "context"

// Tool is the interface for tools that agents can execute.
//
// Maps to Rust langchainrust::core::tools::base::Tool.
type Tool interface {
	// Name returns the tool's unique identifier.
	Name() string

	// Description returns a human-readable description of what the tool does.
	Description() string

	// Run executes the tool with the given input string.
	Run(ctx context.Context, input string) (string, error)
}

// ToolError represents an error that occurred during tool execution.
type ToolError struct {
	ToolName string
	Message  string
	Cause    error
}

func (e *ToolError) Error() string {
	if e.Cause != nil {
		return "tool " + e.ToolName + ": " + e.Message + ": " + e.Cause.Error()
	}
	return "tool " + e.ToolName + ": " + e.Message
}

func (e *ToolError) Unwrap() error {
	return e.Cause
}

// NewToolError creates a new ToolError.
func NewToolError(toolName, message string, cause error) *ToolError {
	return &ToolError{
		ToolName: toolName,
		Message:  message,
		Cause:    cause,
	}
}
