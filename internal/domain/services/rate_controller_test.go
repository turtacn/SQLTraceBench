package services

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTokenBucketRateController(t *testing.T) {
	// Create a rate controller with a QPS of 100 and a max concurrency of 10
	rc := NewTokenBucketRateController(100, 10)

	// Start the rate controller
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rc.Start(ctx)

	// Try to acquire 10 tokens, which should be fast because of the max concurrency
	start := time.Now()
	for i := 0; i < 10; i++ {
		err := rc.Acquire(ctx)
		assert.NoError(t, err)
	}
	duration := time.Since(start)
	assert.Less(t, duration, 150*time.Millisecond, "acquiring the first 10 tokens should be fast")

	// Try to acquire another 10 tokens, which should take about 100ms because of the QPS limit
	start = time.Now()
	for i := 0; i < 10; i++ {
		err := rc.Acquire(ctx)
		assert.NoError(t, err)
	}
	duration = time.Since(start)
	assert.InDelta(t, float64(100*time.Millisecond), float64(duration), float64(50*time.Millisecond), "acquiring the next 10 tokens should take about 100ms")

	// Stop the rate controller
	rc.Stop()
}