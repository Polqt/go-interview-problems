## Web Crawler

In this exercise you'll use Go's concurrency features to parallelize a web crawler.

Modify the `Crawl` function to fetch URLs in parallel without fetching the same URL twice.

~Hint:~ you can keep a cache of the URLs that have been fetched on a map, but maps alone are not safe for concurrent use!

## Tags
`Concurrency`

## Source
- [A Tour of Go](https://go.dev/tour/concurrency/10)
