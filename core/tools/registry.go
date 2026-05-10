package tools

import (
	"fmt"
	"sync"
)

// ToolRegistry manages tool registration and lookup.
// Agents use the registry to discover and invoke tools by name.
//
// Maps to Rust langchainrust::core::tools::registry::ToolRegistry.
type ToolRegistry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewToolRegistry creates an empty ToolRegistry.
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool to the registry.
// Returns an error if a tool with the same name already exists.
func (r *ToolRegistry) Register(tool Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	name := tool.Name()
	if _, exists := r.tools[name]; exists {
		return fmt.Errorf("tool %q already registered", name)
	}
	r.tools[name] = tool
	return nil
}

// RegisterAll registers multiple tools at once.
// On first conflict, stops and returns the error.
func (r *ToolRegistry) RegisterAll(tools ...Tool) error {
	for _, tool := range tools {
		if err := r.Register(tool); err != nil {
			return err
		}
	}
	return nil
}

// Get retrieves a tool by name.
// Returns nil and false if the tool is not found.
func (r *ToolRegistry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, ok := r.tools[name]
	return tool, ok
}

// MustGet retrieves a tool by name, panicking if not found.
func (r *ToolRegistry) MustGet(name string) Tool {
	tool, ok := r.Get(name)
	if !ok {
		panic(fmt.Sprintf("tool %q not found in registry", name))
	}
	return tool
}

func (r *ToolRegistry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tools, name)
}

// List returns all registered tool names.
func (r *ToolRegistry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.tools))
	for name := range r.tools {
		names = append(names, name)
	}
	return names
}

// ToToolDefinitions converts all registered tools to ToolDefinition slice
// for use in LLM function calling API requests.
func (r *ToolRegistry) ToToolDefinitions() []ToolDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	defs := make([]ToolDefinition, 0, len(r.tools))
	for _, tool := range r.tools {
		defs = append(defs, ToolDefinition{
			Name:        tool.Name(),
			Description: tool.Description(),
			Parameters:  make(map[string]any),
		})
	}
	return defs
}

// Len returns the number of registered tools.
func (r *ToolRegistry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.tools)
}
