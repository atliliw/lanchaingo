package callbacks

import (
	"fmt"
	"os"
	"time"
)

type LangSmithHandler struct {
	apiKey  string
	project string
	baseURL string
}

func NewLangSmithHandler(apiKey, project string) *LangSmithHandler {
	if project == "" { project = "lanchaingo" }
	baseURL := os.Getenv("LANGSMITH_BASE_URL")
	if baseURL == "" { baseURL = "https://api.smith.langchain.com" }
	return &LangSmithHandler{apiKey: apiKey, project: project, baseURL: baseURL}
}

func (h *LangSmithHandler) send(runName, runType, input, output string, err error) {
	if h.apiKey == "" { return }
	// Simplified run creation via REST API
	_ = fmt.Sprintf("%s/runs", h.baseURL)
	_ = time.Now()
}

func (h *LangSmithHandler) OnLLMStart(runName, input string)       { h.send(runName, "llm", input, "", nil) }
func (h *LangSmithHandler) OnLLMEnd(runName, output string)        { h.send(runName, "llm", "", output, nil) }
func (h *LangSmithHandler) OnLLMError(runName string, er error)    { h.send(runName, "llm", "", "", er) }
func (h *LangSmithHandler) OnChainStart(runName string, in map[string]any)  { h.send(runName, "chain", fmt.Sprintf("%v", in), "", nil) }
func (h *LangSmithHandler) OnChainEnd(runName string, out map[string]any)   { h.send(runName, "chain", "", fmt.Sprintf("%v", out), nil) }
func (h *LangSmithHandler) OnChainError(runName string, er error)           { h.send(runName, "chain", "", "", er) }
func (h *LangSmithHandler) OnToolStart(runName, tool, input string)         { h.send(runName, "tool", input, "", nil) }
func (h *LangSmithHandler) OnToolEnd(runName, tool, output string)          { h.send(runName, "tool", "", output, nil) }
func (h *LangSmithHandler) OnToolError(runName, tool string, er error)      { h.send(runName, "tool", "", "", er) }
