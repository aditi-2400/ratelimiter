// Package ratelimiter provides rate limiting algorithms with a unified interface.
package ratelimiter

import (
	"context"
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
	Rate   int
	Period time.Duration
	Burst  int
	Store  Store
}
