# ratelimiter

A modular rate limiting library for Go. Provides multiple algorithms behind a single interface, with a pluggable storage backend for in-memory or distributed use.

## Installation

```bash
go get github.com/aditi-2400/ratelimiter
```

## Usage

```go
import "github.com/aditi-2400/ratelimiter"

limiter, err := ratelimiter.New(ratelimiter.TokenBucket, ratelimiter.Config{
    Rate:  10,        // tokens per second
    Burst: 10,        // maximum burst capacity
    Store: myStore,   // see Store section below
})
if err != nil {
    log.Fatal(err)
}

allowed, err := limiter.Allow(ctx, "user-123")
if err != nil {
    // infrastructure failure (e.g. store unreachable)
}
if !allowed {
    // request denied — rate limit exceeded
}
```

## Algorithms

| Algorithm | Status |
|---|---|
| Token Bucket | Available |
| Leaky Bucket | Coming soon |
| Sliding Window Log | Coming soon |
| Sliding Window Counter | Coming soon |

## Config

| Field | Type | Description |
|---|---|---|
| `Rate` | `float64` | Requests (or tokens) permitted per second |
| `Burst` | `int` | Maximum burst capacity |
| `Period` | `time.Duration` | Window duration (used by sliding window algorithms) |
| `Store` | `Store` | Storage backend for per-key state |

## Store

The `Store` interface decouples the algorithm from its storage backend:

```go
type Store interface {
    Get(ctx context.Context, key string) (any, error)
    Set(ctx context.Context, key string, value any, ttl time.Duration) error
    Del(ctx context.Context, key string) error
}
```

Implement this interface to use any backend (in-memory, Redis, etc.). Keys are per-user or per-IP identifiers passed to `Allow`.

## Design

- **Strategy pattern** — each algorithm implements the `Limiter` interface
- **Factory pattern** — `New(algorithm, config)` constructs the right implementation
- **Pluggable storage** — algorithms are decoupled from state storage via `Store`

## License

MIT
