package core

import (
	"context"
	"sync"
	"time"
)

// RateLimiter controls request rate with a token bucket.
type RateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	max      float64
	interval time.Duration
	last     time.Time
}

// NewRateLimiter creates a RateLimiter allowing n operations per interval.
func NewRateLimiter(n int, interval time.Duration) *RateLimiter {
	return &RateLimiter{
		tokens:   float64(n),
		max:      float64(n),
		interval: interval,
		last:     time.Now(),
	}
}

// Wait blocks until a token is available or ctx is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	for {
		r.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(r.last)
		r.tokens += float64(elapsed) / float64(r.interval) * r.max
		if r.tokens > r.max {
			r.tokens = r.max
		}
		if r.tokens >= 1 {
			r.tokens--
			r.last = now
			r.mu.Unlock()
			return nil
		}
		// Need to wait: calculate time until next token
		wait := r.interval - time.Duration(r.tokens/r.max*float64(r.interval))
		r.mu.Unlock()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(wait):
		}
	}
}
