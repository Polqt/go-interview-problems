# Rate Tracker
You are given a `Handler` with a method `Handle()` that performs some work. Your task is to implement the `LogRate` method, which reports how many times `Handle()` is called during each time interval of duration `d`.

The `LogRate` method should periodically report the number of `Handle()` calls that occurred in each interval by calling `monitor.SendRate(int)`.

## Tags
`Concurrency`
