package ratelimiter

import (
	"context"
	"errors"
	"time"
)

// tokenBucket implements Limiter using the token bucket algorithm.
// Tokens accumulate at a fixed rate up to a maximum burst capacity.
// Each allowed request consumes one token.
type tokenBucket struct {
	rate  float64 // tokens added per second
	burst int     // maximum token capacity
	store Store
}

// tokenBucketState is the per-key state persisted in the Store.
type tokenBucketState struct {
	tokens     float64
	lastRefill time.Time
}

// newTokenBucket creates a TokenBucket limiter from cfg.
func newTokenBucket(cfg Config) (*tokenBucket, error) {

	if cfg.Rate <= 0 {
		return nil, errors.New("ratelimiter: Rate must be greater than zero")
	}

	if cfg.Burst <= 0 {
		return nil, errors.New("ratelimiter: Burst must be greater than zero")
	}

	return &tokenBucket{
		rate:  cfg.Rate,
		burst: cfg.Burst,
		store: cfg.Store,
	}, nil
}

// Allow returns true if the request for key is within the rate limit.
// It refills tokens based on time elapsed since the last request, then
// consumes one token if available. Safe for concurrent use via the Store.
func (tb *tokenBucket) Allow(ctx context.Context, key string) (bool, error) {
	val, err := tb.store.Get(ctx, key)

	if err != nil {
		return false, err
	}

	state, ok := val.(*tokenBucketState)

	if !ok || state == nil {

		state = &tokenBucketState{

			tokens:     float64(tb.burst),
			lastRefill: time.Now(),
		}
	}

	elapsed := time.Since(state.lastRefill)
	tokensEarned := elapsed.Seconds() * tb.rate
	state.tokens += tokensEarned

	if state.tokens > float64(tb.burst) {
		state.tokens = float64(tb.burst)
	}

	state.lastRefill = time.Now()
	ttl := time.Duration(float64(tb.burst) / tb.rate * float64(time.Second))

	if state.tokens >= 1 {

		state.tokens -= 1
		err = tb.store.Set(ctx, key, state, ttl)
		if err != nil {

			return false, err
		}
		return true, nil

	}

	err = tb.store.Set(ctx, key, state, ttl)

	if err != nil {

		return false, err

	}
	return false, nil

}
