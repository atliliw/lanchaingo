package langgraph

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

// PostgresCheckpointer stores checkpoints in PostgreSQL.
type PostgresCheckpointer struct {
	mu   sync.Mutex
	db   *sql.DB
	keys []string
}

func NewPostgresCheckpointer(connStr string) (*PostgresCheckpointer, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres: ping: %w", err)
	}

	pc := &PostgresCheckpointer{db: db}
	if err := pc.init(); err != nil {
		return nil, err
	}
	return pc, nil
}

func (pc *PostgresCheckpointer) init() error {
	_, err := pc.db.Exec(`
		CREATE TABLE IF NOT EXISTS langgraph_checkpoints (
			id TEXT PRIMARY KEY,
			state JSONB,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	return err
}

func (pc *PostgresCheckpointer) Save(state StateSchema) (string, error) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	id := fmt.Sprintf("cp_%d", time.Now().UnixNano())
	data, err := json.Marshal(state)
	if err != nil {
		return "", err
	}

	_, err = pc.db.Exec(
		"INSERT INTO langgraph_checkpoints (id, state) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET state = $2",
		id, data,
	)
	if err != nil {
		return "", fmt.Errorf("postgres: save: %w", err)
	}
	pc.keys = append(pc.keys, id)
	return id, nil
}

func (pc *PostgresCheckpointer) Load(id string) (StateSchema, error) {
	var data []byte
	err := pc.db.QueryRow("SELECT state FROM langgraph_checkpoints WHERE id = $1", id).Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("postgres: load: %w", err)
	}
	var state AgentState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, err
	}
	return &state, nil
}

func (pc *PostgresCheckpointer) List() ([]string, error) {
	rows, err := pc.db.Query("SELECT id FROM langgraph_checkpoints ORDER BY created_at")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		rows.Scan(&id)
		ids = append(ids, id)
	}
	return ids, nil
}

func (pc *PostgresCheckpointer) Close() error {
	return pc.db.Close()
}
