package ratelimiter

import (
	"context"
	"time"
)

// Store is a key-value backend for persisting per-key rate limit state.
// Implementations must be safe for concurrent use.
type Store interface {
	// Get returns the state stored for key, or nil if no entry exists.
	Get(ctx context.Context, key string) (any, error)

	// Set saves value for key and removes it automatically after ttl elapses.
	Set(ctx context.Context, key string, value any, ttl time.Duration) error

	// Del removes the state for key immediately.
	Del(ctx context.Context, key string) error
}
