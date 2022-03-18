package limiter_test

import (
	"sync"
	"testing"

	"github.com/korovkin/limiter"
)

func TestExample(t *testing.T) {
	limit := limiter.NewConcurrencyLimiter(10)
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

	c := limiter.NewConcurrencyLimiter(LIMIT)
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

	if len(m) != N {
		t.Error("invalid num of results", len(m))
	}

	if max > int32(LIMIT) || max == 0 {
		t.Error("invalid max", max)
	}
}

func TestExecuteWithTicket(t *testing.T) {
	LIMIT := 10
	N := 100
	c := limiter.NewConcurrencyLimiter(LIMIT)
	m := map[int]int{}
	lock := &sync.Mutex{}

	for i := 0; i < N; i++ {
		c.ExecuteWithTicket(func(ticket int) {
			lock.Lock()
			m[ticket] += 1
			if ticket > LIMIT-1 {
				t.Errorf("expected max ticket: %d, got %d", LIMIT, ticket)
			}
			lock.Unlock()
		})
	}
	c.Wait()

	sum := 0
	for _, count := range m {
		sum += count
	}
	if sum != N {
		t.Errorf("invalid num of results: %d, expected %d", sum, N)
	}
}
