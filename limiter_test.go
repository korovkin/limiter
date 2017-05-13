package limiter

import (
	"sync"
	"testing"
)

func TestHello(t *testing.T) {
	t.Log("hello")
}

func TestLimit(t *testing.T) {
	c := NewConcurrencyLimiter(10)
	m := map[int]bool{}
	lock := &sync.Mutex{}

	N := 100

	for i := 0; i < N; i++ {
		x := i
		c.Execute(func() {
			lock.Lock()
			m[x] = true
			lock.Unlock()
		})
	}

	c.Wait()

	if len(m) != N {
		t.Error("invalid num of results", len(m))
	}
}
