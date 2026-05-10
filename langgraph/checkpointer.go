package langgraph

import "fmt"

// Checkpointer allows saving and loading graph state for persistence/resume.
type Checkpointer interface {
	Save(state StateSchema) (string, error)
	Load(id string) (StateSchema, error)
	List() ([]string, error)
}

// MemoryCheckpointer stores checkpoints in memory.
type MemoryCheckpointer struct {
	checkpoints map[string]StateSchema
	order       []string
}

func NewMemoryCheckpointer() *MemoryCheckpointer {
	return &MemoryCheckpointer{
		checkpoints: make(map[string]StateSchema),
		order:       make([]string, 0),
	}
}

func (mc *MemoryCheckpointer) Save(state StateSchema) (string, error) {
	id := fmt.Sprintf("cp_%d", len(mc.order))
	mc.checkpoints[id] = state.CloneState()
	mc.order = append(mc.order, id)
	return id, nil
}

func (mc *MemoryCheckpointer) Load(id string) (StateSchema, error) {
	state, ok := mc.checkpoints[id]
	if !ok {
		return nil, fmt.Errorf("checkpoint '%s' not found", id)
	}
	return state.CloneState(), nil
}

func (mc *MemoryCheckpointer) List() ([]string, error) {
	result := make([]string, len(mc.order))
	copy(result, mc.order)
	return result, nil
}
