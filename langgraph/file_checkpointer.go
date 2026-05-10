package langgraph

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// FileCheckpointer stores checkpoints on disk.
type FileCheckpointer struct {
	mu   sync.Mutex
	dir  string
	keys []string
}

func NewFileCheckpointer(dir string) *FileCheckpointer {
	os.MkdirAll(dir, 0755)
	return &FileCheckpointer{dir: dir}
}

func (fc *FileCheckpointer) Save(state StateSchema) (string, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	id := "cp_" + itoa(len(fc.keys))
	data, _ := json.Marshal(state)
	os.WriteFile(filepath.Join(fc.dir, id+".json"), data, 0644)
	fc.keys = append(fc.keys, id)
	return id, nil
}

func (fc *FileCheckpointer) Load(id string) (StateSchema, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()

	data, err := os.ReadFile(filepath.Join(fc.dir, id+".json"))
	if err != nil {
		return nil, err
	}
	var state AgentState
	json.Unmarshal(data, &state)
	return &state, nil
}

func (fc *FileCheckpointer) List() ([]string, error) {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	r := make([]string, len(fc.keys))
	copy(r, fc.keys)
	return r, nil
}

func itoa(n int) string {
	if n == 0 { return "0" }
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
