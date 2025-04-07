# Merge Channels
Implement the `Merge` function that takes a variable number of input channels and returns a single output channel. The output channel should receive all values from all input channels. Once all input channels are closed and all values have been sent to the output channel, the output channel should be closed.

Requirements:

1. The function should handle any number of input channels, including zero.
2. All values from all input channels must be sent to the output channel.
3. The output channel should be closed only after all input channels are closed and all values have been forwarded.
4. Your implementation must not leak goroutines.
5. Values should be forwarded as soon as they are available from any input channel.

## Tags
`Concurrency`
