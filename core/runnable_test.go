package core

import (
	"context"
	"testing"
)

func TestRunnableFunc(t *testing.T) {
	fn := RunnableFunc(func(ctx context.Context, input map[string]any) (any, error) { return "ok", nil })
	result, err := fn.Invoke(context.Background(), map[string]any{"key": "val"})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if result != "ok" { t.Errorf("expected ok, got %v", result) }
}
func TestRunnableFuncBatch(t *testing.T) {
	fn := RunnableFunc(func(ctx context.Context, input map[string]any) (any, error) { return input["x"], nil })
	inputs := []map[string]any{{"x": 1}, {"x": 2}, {"x": 3}}
	results, err := fn.Batch(context.Background(), inputs)
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	if len(results) != 3 { t.Fatalf("expected 3 results, got %d", len(results)) }
}
func TestRunnableFuncStream(t *testing.T) {
	fn := RunnableFunc(func(ctx context.Context, input map[string]any) (any, error) { return "streamed", nil })
	ch, err := fn.Stream(context.Background(), map[string]any{"x": 1})
	if err != nil { t.Fatalf("unexpected error: %v", err) }
	result := <-ch
	if result != "streamed" { t.Errorf("expected streamed, got %v", result) }
}
