package tools

import (
	"context"
	"time"
)

// DateTimeTool returns current date and time.
type DateTimeTool struct{}

func NewDateTimeTool() *DateTimeTool { return &DateTimeTool{} }

func (d *DateTimeTool) Name() string { return "datetime" }

func (d *DateTimeTool) Description() string {
	return "returns current date and time"
}

func (d *DateTimeTool) Run(_ context.Context, _ string) (string, error) {
	return time.Now().Format("2006-01-02 15:04:05"), nil
}
