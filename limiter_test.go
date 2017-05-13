package limiter

import (
	"sync"
	"testing"
)

func TestHello(t *testing.T) {
	t.Log("hello")
}

func TestExample(t *testing.T) {
	limit := NewConcurrencyLimiter(10)
	for i := 0; i < 1000; i++ {
		limit.Execute(func() {
			// do some work
		})
	}
	limit.Wait()
}

func TestLimit(t *testing.T) {
	LIMIT := 10
	N := 100

	c := NewConcurrencyLimiter(LIMIT)
	m := map[int]bool{}
	lock := &sync.Mutex{}

	max := int32(0)
	for i := 0; i < N; i++ {
		x := i
		c.Execute(func() {
			lock.Lock()
			m[x] = true
			currentMax := c.GetNumInProgress()
			if currentMax >= max {
				max = currentMax
			}
			lock.Unlock()
		})
	}

	// wait until the above completes
	c.Wait()

	t.Log("results:", len(m))
	t.Log("max:", max)

	if len(m) != int(N) {
		t.Error("invalid num of results", len(m))
	}

	if max > int32(LIMIT) {
		t.Error("invalid max", max)
	}
}
