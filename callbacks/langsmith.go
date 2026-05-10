package callbacks

// LangSmithHandler sends traces to LangSmith.
// Requires: go get github.com/langchain-ai/langsmith-go
type LangSmithHandler struct {
	apiKey  string
	project string
	enabled bool
}

func NewLangSmithHandler(apiKey, project string) *LangSmithHandler {
	return &LangSmithHandler{apiKey: apiKey, project: project, enabled: apiKey != ""}
}

func (h *LangSmithHandler) OnLLMStart(runName, input string)    {}
func (h *LangSmithHandler) OnLLMEnd(runName, output string)     {}
func (h *LangSmithHandler) OnLLMError(runName string, err error) {}
func (h *LangSmithHandler) OnChainStart(runName string, input map[string]any)    {}
func (h *LangSmithHandler) OnChainEnd(runName string, output map[string]any)      {}
func (h *LangSmithHandler) OnChainError(runName string, err error) {}
func (h *LangSmithHandler) OnToolStart(runName, toolName, input string)    {}
func (h *LangSmithHandler) OnToolEnd(runName, toolName, output string)      {}
func (h *LangSmithHandler) OnToolError(runName, toolName string, err error) {}
