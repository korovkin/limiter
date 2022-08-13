# go lang goroutine concurrency limiter

## builds

[![CircleCI](https://dl.circleci.com/status-badge/img/gh/JaderDias/limiter/tree/main.svg?style=svg)](https://dl.circleci.com/status-badge/redirect/gh/JaderDias/limiter/tree/main)

## Example

limit the number of concurrent go routines to 10:

```
  import "github.com/JaderDias/limiter"

  ...

  limiter.BoundedConcurrency(10, 1000, func(i int) {
      // do some work
  })
```

## Real World Example:

```
  import "github.com/JaderDias/limiter"

  ...

  urls := []string{
    "http://www.google.com",
    "http://www.apple.com",
  }
  results := make([]int, 2)
  limiter.BoundedConcurrency(10, 2, func(i int) {
    resp, err := http.Get(urls[i])
    Expect(err).To(BeNil())
    defer resp.Body.Close()
    results[i] = resp.StatusCode
  })

  log.Println("httpGoogle:", results[0])
  log.Println("httpApple:", results[1])
```

## Concurrent IO with Error tracking:

```
  import "github.com/JaderDias/limiter"
  
  ...

  errors := []error{
    errors.New("error a"),
    errors.New("error b"),
  }
  var firstError atomic.Value
  completed := int32(0)
  limiter.BoundedConcurrency(4, 2, func(i int) {
    atomic.AddInt32(&completed, 1)
    // Do some really slow IO ...
    // keep the error:
    firstError.CompareAndSwap(nil, errors[i])
  })

  firstErrorValue := firstError.Load().(error)
  Expect(firstErrorValue == errors[0] || firstErrorValue == errors[1]).To(BeTrue())
```
