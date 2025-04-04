# Rate Limiter

Implement a rate limiter that allows up to `n` operations per second. It should provide two methods:

- `Take()` — blocks until the operation is allowed.
- `CanTake()` — returns immediately with `true` if the operation is allowed, otherwise `false`.

The limiter should pace operations using a token-based mechanism like a token bucket.

## Tags  
`Concurrency`
