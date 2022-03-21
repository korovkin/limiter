# go lang goroutine concurrency limiter

## builds

[![Build Status](https://travis-ci.org/korovkin/limiter.svg)](https://travis-ci.org/korovkin/limiter)

## Example

limit the number of concurrent go routines to 10:

```
  import "github.com/korovkin/limiter"

  ...

  limit := limiter.NewConcurrencyLimiter(10)
  defer limit.WaitAndClose()

  for i := 0; i < 1000; i++ {
    limit.Execute(func() {
      // do some work
    })
  }
```

## Real World Example:

```
  import "github.com/korovkin/limiter"

  ...

  limiter := limiter.NewConcurrencyLimiter(10)

  httpGoogle := int(0)
  limiter.Execute(func() {
    resp, err := http.Get("https://www.google.com/")
    Expect(err).To(BeNil())
    defer resp.Body.Close()
    httpGoogle = resp.StatusCode
  })

  httpApple := int(0)
  limiter.Execute(func() {
    resp, err := http.Get("https://www.apple.com/")
    Expect(err).To(BeNil())
    defer resp.Body.Close()
    httpApple = resp.StatusCode
  })

  limiter.WaitAndClose()

  log.Println("httpGoogle:", httpGoogle)
  log.Println("httpApple:", httpApple)
```
