package limiter_test

import (
	"log"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/korovkin/limiter"

	. "github.com/onsi/gomega"
)

func TestExample(t *testing.T) {
	RegisterTestingT(t)

	t.Run("TestExample", func(*testing.T) {
		x := int32(1000)
		limit := limiter.NewConcurrencyLimiter(10)
		for i := 0; i < 1000; i++ {
			limit.Execute(func() {
				// do some work:
				atomic.AddInt32(&x, -1)
			})
		}
		limit.Wait()
		defer limit.Close()
		Expect(limit.GetNumInProgress()).To(BeEquivalentTo(0))
		Expect(x).To(BeEquivalentTo(0))
	})
}

func TestLimit(t *testing.T) {
	RegisterTestingT(t)

	t.Run("TestLimit", func(*testing.T) {
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

		Expect(max).To(BeEquivalentTo(10))
		Expect(len(m)).To(BeEquivalentTo(N))
		Expect(c.GetNumInProgress()).To(BeEquivalentTo(0))

		c.Close()
		_, err := c.Execute(func() {
			log.Println("more")
		})
		Expect(err).ToNot(BeNil())
	})
}

func TestExecuteWithTicket(t *testing.T) {
	RegisterTestingT(t)

	t.Run("TestExecuteWithTicket", func(t *testing.T) {
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

		Expect(sum).To(BeEquivalentTo(N))
		Expect(c.GetNumInProgress()).To(BeEquivalentTo(0))

		c.Close()
		_, err := c.Execute(func() {
			log.Println("more")
		})
		Expect(err).ToNot(BeNil())
	})
}
