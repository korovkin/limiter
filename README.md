# limiter
go lang concurrency limiter.

## example

limit number of concurrent go routines to 10

```
	limit := NewConcurrencyLimiter(10)
	for i := 0; i < 1000; i++ {
		limit.Execute(func() {
			// do some work
		})
	}
	limit.Wait()
```
