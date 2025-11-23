package services

import (
	"context"
	"sync"
	"time"
)

// RateController defines the interface for controlling the rate and concurrency of execution.
type RateController interface {
	Start(ctx context.Context)
	Acquire(ctx context.Context) error
	Stop()
	MaxConcurrency() int
}

// TokenBucketRateController implements a rate controller using the token bucket algorithm
// combined with time.Sleep to ensure smooth distribution.
type TokenBucketRateController struct {
	qps            int
	maxConcurrency int

	// mu protects the state of the rate limiter
	mu              sync.Mutex
	nextAllowedTime time.Time
}

// NewTokenBucketRateController creates a new token bucket rate controller.
func NewTokenBucketRateController(targetQPS, maxConcurrency int) *TokenBucketRateController {
	if targetQPS <= 0 {
		targetQPS = 100
	}
	if maxConcurrency <= 0 {
		maxConcurrency = 10
	}

	return &TokenBucketRateController{
		qps:            targetQPS,
		maxConcurrency: maxConcurrency,
	}
}

// Start initializes the rate controller.
// For this implementation, no background goroutine is needed as we calculate sleeps on demand.
func (c *TokenBucketRateController) Start(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextAllowedTime = time.Now()
}

// Acquire blocks until a request is allowed to proceed based on the configured QPS.
// It uses time.Sleep to smooth out the request rate.
func (c *TokenBucketRateController) Acquire(ctx context.Context) error {
	c.mu.Lock()
	now := time.Now()
	interval := time.Second / time.Duration(c.qps)

	// If we are behind schedule (idle), reset nextAllowedTime to now.
	if now.After(c.nextAllowedTime) {
		c.nextAllowedTime = now
	}

	// Calculate the time this request is allowed to start.
	targetTime := c.nextAllowedTime
	// Advance the next allowed time for the next request.
	c.nextAllowedTime = targetTime.Add(interval)
	c.mu.Unlock()

	// Calculate how long we need to sleep.
	sleepDuration := targetTime.Sub(now)
	if sleepDuration > 0 {
		select {
		case <-time.After(sleepDuration):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

// Stop cleans up resources.
func (c *TokenBucketRateController) Stop() {
	// No resources to clean up in this implementation.
}

// MaxConcurrency returns the maximum concurrency of the rate controller.
func (c *TokenBucketRateController) MaxConcurrency() int {
	return c.maxConcurrency
}
