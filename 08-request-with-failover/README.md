# Request With Failover
Attempts to retrieve data from multiple sources with timeout-based failover logic.

Behavior:
- Launches a request to each address with a 500ms timeout.
- If a request errors, the next address is tried immediately.
- If a request hangs, the next one is attempted in parallel after 500ms, but the previous request is still allowed to complete.
- Returns the first successful result from any address.
- If all attempts fail,  returns `ErrRequestsFailed`

## Tags
`Concurrency` `Context` `Failover` 
