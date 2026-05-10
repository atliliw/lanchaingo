package tools

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// PythonREPLTool executes Python code.
type PythonREPLTool struct{}

func NewPythonREPLTool() *PythonREPLTool { return &PythonREPLTool{} }

func (p *PythonREPLTool) Name() string { return "python_repl" }

func (p *PythonREPLTool) Description() string {
	return "execute Python code and return the result"
}

func (p *PythonREPLTool) Run(_ context.Context, input string) (string, error) {
	cmd := exec.Command("python3", "-c", input)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return "", fmt.Errorf("python error: %s", string(exitErr.Stderr))
		}
		return "", fmt.Errorf("failed to execute python: %w", err)
	}

	result := strings.TrimSpace(string(output))
	if result == "" {
		return "(no output)", nil
	}
	return result, nil
}
