package handlers

import (
	"fmt"
	"os"
	"sync"
	"time"
)

// FileCallbackHandler writes callbacks to a file.
type FileCallbackHandler struct {
	mu   sync.Mutex
	file *os.File
}

func NewFileCallbackHandler(filePath string) (*FileCallbackHandler, error) {
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open callback file: %w", err)
	}
	return &FileCallbackHandler{file: f}, nil
}

func (h *FileCallbackHandler) write(event string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	fmt.Fprintf(h.file, "[%s] %s\n", time.Now().Format(time.RFC3339), event)
}

func (h *FileCallbackHandler) OnLLMStart(runName, input string)  { h.write(fmt.Sprintf("LLM START %s", runName)) }
func (h *FileCallbackHandler) OnLLMEnd(runName, output string)   { h.write(fmt.Sprintf("LLM END %s", runName)) }
func (h *FileCallbackHandler) OnLLMError(runName string, err error) { h.write(fmt.Sprintf("LLM ERROR %s: %v", runName, err)) }
func (h *FileCallbackHandler) OnChainStart(runName string, input map[string]any) { h.write(fmt.Sprintf("CHAIN START %s", runName)) }
func (h *FileCallbackHandler) OnChainEnd(runName string, output map[string]any) { h.write(fmt.Sprintf("CHAIN END %s", runName)) }
func (h *FileCallbackHandler) OnChainError(runName string, err error) { h.write(fmt.Sprintf("CHAIN ERROR %s: %v", runName, err)) }
func (h *FileCallbackHandler) OnToolStart(runName, toolName, input string) { h.write(fmt.Sprintf("TOOL START %s", toolName)) }
func (h *FileCallbackHandler) OnToolEnd(runName, toolName, output string) { h.write(fmt.Sprintf("TOOL END %s", toolName)) }
func (h *FileCallbackHandler) OnToolError(runName, toolName string, err error) { h.write(fmt.Sprintf("TOOL ERROR %s: %v", toolName, err)) }

func (h *FileCallbackHandler) Close() error {
	return h.file.Close()
}
