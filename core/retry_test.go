package core

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestDoWithRetrySuccess(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 2

	attempts := 0
	result, err := DoWithRetry(ctx, cfg, func(ctx context.Context) (string, error) {
		attempts++
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "ok" {
		t.Errorf("expected 'ok', got '%s'", result)
	}
	if attempts != 1 {
		t.Errorf("expected 1 attempt, got %d", attempts)
	}
}

func TestDoWithRetryRetriesOnFailure(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.BaseDelay = time.Millisecond

	attempts := 0
	_, err := DoWithRetry(ctx, cfg, func(ctx context.Context) (string, error) {
		attempts++
		return "", errors.New("server error")
	})
	if err == nil {
		t.Fatal("expected error")
	}
	if attempts != 4 { // 1 initial + 3 retries
		t.Errorf("expected 4 attempts, got %d", attempts)
	}
}

func TestDoWithRetrySucceedsAfterRetry(t *testing.T) {
	ctx := context.Background()
	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.BaseDelay = time.Millisecond

	attempts := 0
	result, err := DoWithRetry(ctx, cfg, func(ctx context.Context) (string, error) {
		attempts++
		if attempts < 3 {
			return "", errors.New("transient error")
		}
		return "recovered", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "recovered" {
		t.Errorf("expected 'recovered', got '%s'", result)
	}
	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestDoWithRetryContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := DefaultRetryConfig()
	cfg.MaxRetries = 3
	cfg.BaseDelay = time.Millisecond

	_, err := DoWithRetry(ctx, cfg, func(ctx context.Context) (string, error) {
		return "", errors.New("error")
	})
	if err == nil {
		t.Fatal("expected context cancelled error")
	}
}

func TestRetryableHTTP(t *testing.T) {
	if !RetryableHTTP(429) {
		t.Error("429 should be retryable")
	}
	if !RetryableHTTP(500) {
		t.Error("500 should be retryable")
	}
	if RetryableHTTP(400) {
		t.Error("400 should not be retryable")
	}
	if RetryableHTTP(200) {
		t.Error("200 should not be retryable")
	}
}
