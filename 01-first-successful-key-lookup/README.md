# First Successful Key Lookup

You are given a function `get(ctx, address, key)`, which retrieves a value for a given key from a remote address. However, you don't know which address has the value, and network calls can fail.

Your task is to implement the Get function that:
* Calls get for each address in parallel.
* Returns the first successful response.
* If all requests fail, returns an error.

## Tags
`Concurrency`

## Source
- [Mock-собеседование по Go от Team Lead из Яндекса](https://www.youtube.com/watch?v=x689QxR3AIc) by **it-interview**
