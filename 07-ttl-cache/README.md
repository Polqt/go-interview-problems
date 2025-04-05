# Cache with Expiration
Implement a basic in-memory cache where each key can have an optional time-to-live (TTL).  
Keys should expire automatically after their TTL passes.

The cache should support:
- `Set(key, value, ttl)` — stores a key with a given TTL (0 = no expiration)
- `Get(key)` — returns value if key exists and hasn’t expired
- `Delete(key)` — deletes the key

## Tags
`Concurrency`
