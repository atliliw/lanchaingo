package handlers

import "fmt"

// StdOutHandler prints callbacks to stdout for debugging.
type StdOutHandler struct{}

func NewStdOutHandler() *StdOutHandler { return &StdOutHandler{} }

func (h *StdOutHandler) OnLLMStart(runName, input string) {
	fmt.Printf("[LLM START] %s\n  Input: %s\n", runName, truncate(input, 100))
}
func (h *StdOutHandler) OnLLMEnd(runName, output string) {
	fmt.Printf("[LLM END] %s\n  Output: %s\n", runName, truncate(output, 100))
}
func (h *StdOutHandler) OnLLMError(runName string, err error) {
	fmt.Printf("[LLM ERROR] %s: %v\n", runName, err)
}
func (h *StdOutHandler) OnChainStart(runName string, input map[string]any) {
	fmt.Printf("[CHAIN START] %s\n", runName)
}
func (h *StdOutHandler) OnChainEnd(runName string, output map[string]any) {
	fmt.Printf("[CHAIN END] %s\n", runName)
}
func (h *StdOutHandler) OnChainError(runName string, err error) {
	fmt.Printf("[CHAIN ERROR] %s: %v\n", runName, err)
}
func (h *StdOutHandler) OnToolStart(runName, toolName, input string) {
	fmt.Printf("[TOOL START] %s(%s)\n", toolName, truncate(input, 100))
}
func (h *StdOutHandler) OnToolEnd(runName, toolName, output string) {
	fmt.Printf("[TOOL END] %s: %s\n", toolName, truncate(output, 100))
}
func (h *StdOutHandler) OnToolError(runName, toolName string, err error) {
	fmt.Printf("[TOOL ERROR] %s: %v\n", toolName, err)
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n]) + "..."
}
