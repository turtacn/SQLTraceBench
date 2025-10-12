package services

import (
	"context"
	"time"
)

// RateController defines the interface for controlling the rate and concurrency of execution.
type RateController interface {
	Start(ctx context.Context)
	Acquire(ctx context.Context) error
	Stop()
	MaxConcurrency() int
}

// TokenBucketRateController implements a rate controller using the token bucket algorithm.
type TokenBucketRateController struct {
	qps            int
	maxConcurrency int
	tokens         chan struct{}
	ticker         *time.Ticker
	stopCh         chan struct{}
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
		tokens:         make(chan struct{}, maxConcurrency),
		stopCh:         make(chan struct{}),
	}
}

// Start begins the token generation process.
func (c *TokenBucketRateController) Start(ctx context.Context) {
	c.ticker = time.NewTicker(time.Second / time.Duration(c.qps))
	go func() {
		for i := 0; i < c.maxConcurrency; i++ {
			c.tokens <- struct{}{}
		}
		for {
			select {
			case <-c.stopCh:
				c.ticker.Stop()
				return
			case <-c.ticker.C:
				select {
				case c.tokens <- struct{}{}:
				default:
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// Acquire attempts to get a token from the bucket.
func (c *TokenBucketRateController) Acquire(ctx context.Context) error {
	select {
	case <-c.tokens:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Stop terminates the token generation goroutine.
func (c *TokenBucketRateController) Stop() {
	close(c.stopCh)
}

// MaxConcurrency returns the maximum concurrency of the rate controller.
func (c *TokenBucketRateController) MaxConcurrency() int {
	return c.maxConcurrency
}