package ratelimit

import (
	"context"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	rate       int           // requests per second
	tokens     int           // available tokens
	maxTokens  int           // maximum tokens
	ticker     *time.Ticker
	mu         sync.Mutex
	stopCh     chan struct{}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	rl := &RateLimiter{
		rate:      requestsPerSecond,
		tokens:    requestsPerSecond,
		maxTokens: requestsPerSecond,
		ticker:    time.NewTicker(time.Second),
		stopCh:    make(chan struct{}),
	}

	go rl.refillTokens()
	return rl
}

// Wait blocks until a token is available
func (rl *RateLimiter) Wait(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if rl.tryAcquire() {
				return nil
			}
			time.Sleep(10 * time.Millisecond)
		}
	}
}

// tryAcquire attempts to acquire a token
func (rl *RateLimiter) tryAcquire() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if rl.tokens > 0 {
		rl.tokens--
		return true
	}
	return false
}

// refillTokens refills the bucket with tokens
func (rl *RateLimiter) refillTokens() {
	for {
		select {
		case <-rl.ticker.C:
			rl.mu.Lock()
			rl.tokens = rl.maxTokens
			rl.mu.Unlock()
		case <-rl.stopCh:
			rl.ticker.Stop()
			return
		}
	}
}

// Stop stops the rate limiter
func (rl *RateLimiter) Stop() {
	close(rl.stopCh)
}
