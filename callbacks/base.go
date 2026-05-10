package callbacks

// CallbackHandler defines the interface for execution tracing and monitoring.
type CallbackHandler interface {
	OnLLMStart(runName string, input string)
	OnLLMEnd(runName string, output string)
	OnLLMError(runName string, err error)
	OnChainStart(runName string, input map[string]any)
	OnChainEnd(runName string, output map[string]any)
	OnChainError(runName string, err error)
	OnToolStart(runName string, toolName string, input string)
	OnToolEnd(runName string, toolName string, output string)
	OnToolError(runName string, toolName string, err error)
}

// CallbackManager manages multiple callback handlers.
type CallbackManager struct {
	handlers []CallbackHandler
}

func NewCallbackManager() *CallbackManager {
	return &CallbackManager{handlers: make([]CallbackHandler, 0)}
}

func (m *CallbackManager) AddHandler(h CallbackHandler) {
	m.handlers = append(m.handlers, h)
}

func (m *CallbackManager) Handlers() []CallbackHandler {
	return m.handlers
}

func (m *CallbackManager) IsEmpty() bool {
	return len(m.handlers) == 0
}
