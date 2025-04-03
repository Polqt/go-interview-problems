# Non-Blocking Cache
Implement a non-blocking cache that stores previously fetched responses from URLs. The cache should:

* Return cached data if the requested URL has been fetched before.
* Fetch and store new data asynchronously if the URL is not in the cache.
* Ensure multiple concurrent requests for the same URL do not trigger multiple fetches.

## Tags
`Concurrency`

## Source
- [Let's implement a concurrent non-blocking cache in Go](https://youtu.be/KlDWmTcyXdA?si=2Vz9-Y1tp_a-Qow1) by **Konrad Reiche**
