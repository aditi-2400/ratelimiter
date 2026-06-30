// Package ratelimiter provides rate limiting algorithms with a unified interface.
package ratelimiter

import (
	"context"
	"fmt"
	"time"
)

type Limiter interface {

	// Allow returns true if the request identified by key(ip, user) is permitted.
	// Returns an error only for infrastructure failures (e.g. store unreachable),
	// not for rate-limited requests — those return false, nil.
	Allow(ctx context.Context, key string) (bool, error)
}

type Algorithm int

const (
	TokenBucket Algorithm = iota
	LeakyBucket
	SlidingWindowLog
	SlidingWindowCounter
)

type Config struct {
	Rate   float64
	Period time.Duration
	Burst  int
	Store  Store
}

// New returns a Limiter using the given algorithm and config.
func New(algo Algorithm, cfg Config) (Limiter, error) {

	switch algo {

	case TokenBucket:
		return newTokenBucket(cfg)

	case LeakyBucket, SlidingWindowLog, SlidingWindowCounter:
		return nil, fmt.Errorf("ratelimiter: %v not yet implemented", algo)

	default:
		return nil, fmt.Errorf("ratelimiter: unknown algorithm %v", algo)

	}
}
