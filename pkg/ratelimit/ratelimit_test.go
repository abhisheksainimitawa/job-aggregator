package ratelimit

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter_Wait(t *testing.T) {
	rl := NewRateLimiter(10) // 10 requests per second
	defer rl.Stop()

	ctx := context.Background()

	// Test acquiring tokens
	start := time.Now()
	for i := 0; i < 10; i++ {
		err := rl.Wait(ctx)
		if err != nil {
			t.Fatalf("Wait() error = %v", err)
		}
	}

	// Should complete almost instantly (within 100ms)
	duration := time.Since(start)
	if duration > 100*time.Millisecond {
		t.Errorf("Expected fast acquisition, took %v", duration)
	}

	// Next request should wait for refill
	start = time.Now()
	err := rl.Wait(ctx)
	if err != nil {
		t.Fatalf("Wait() error = %v", err)
	}

	duration = time.Since(start)
	if duration < 900*time.Millisecond {
		t.Errorf("Expected to wait ~1s for refill, waited %v", duration)
	}
}

func TestRateLimiter_ContextCancellation(t *testing.T) {
	rl := NewRateLimiter(1)
	defer rl.Stop()

	// Exhaust tokens
	ctx := context.Background()
	rl.Wait(ctx)

	// Create context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Should return error immediately
	err := rl.Wait(ctx)
	if err == nil {
		t.Error("Expected error from cancelled context")
	}
}
