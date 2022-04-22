package limiter_test

import (
	"errors"
	"log"
	"net/http"
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
		limit.WaitAndClose()
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
		c.WaitAndClose()

		Expect(max).To(BeEquivalentTo(10))
		Expect(len(m)).To(BeEquivalentTo(N))
		Expect(c.GetNumInProgress()).To(BeEquivalentTo(0))

		_, err := c.Execute(func() {
			log.Println("more")
		})
		Expect(err).ToNot(BeNil())
	})
}

func TestExecuteWithTicket(t *testing.T) {
	RegisterTestingT(t)

	t.Run("TestExecuteWithTicket", func(*testing.T) {
		LIMIT := 10
		N := 100
		c := limiter.NewConcurrencyLimiter(LIMIT)
		m := map[int]int{}
		lock := &sync.Mutex{}

		for i := 0; i < N; i++ {
			c.ExecuteWithTicket(func(ticket int) {
				lock.Lock()
				m[ticket] += 1
				Expect(ticket).To(BeNumerically("<", LIMIT))
				lock.Unlock()
			})
		}
		c.WaitAndClose()

		sum := 0
		for _, count := range m {
			sum += count
		}

		Expect(sum).To(BeEquivalentTo(N))
		Expect(c.GetNumInProgress()).To(BeEquivalentTo(0))

		_, err := c.Execute(func() {
			log.Println("more ...")
		})
		Expect(err).ToNot(BeNil())
	})
}

func TestConcurrentIO(t *testing.T) {
	RegisterTestingT(t)

	t.Run("TestConcurrentIO", func(*testing.T) {
		c := limiter.NewConcurrencyLimiter(10)

		httpGoogle := int(0)
		c.Execute(func() {
			resp, err := http.Get("https://www.google.com/")
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			httpGoogle = resp.StatusCode
		})
		httpApple := int(0)
		c.Execute(func() {
			resp, err := http.Get("https://www.apple.com/")
			Expect(err).To(BeNil())
			defer resp.Body.Close()
			httpApple = resp.StatusCode
		})
		c.WaitAndClose()

		Expect(httpGoogle).To(BeEquivalentTo(200))
		Expect(httpApple).To(BeEquivalentTo(200))
	})
}

func TestConcurrently(t *testing.T) {
	RegisterTestingT(t)

	t.Run("TestConcurrently", func(*testing.T) {
		a := errors.New("error a")
		b := errors.New("error b")
		completed := int32(0)

		concurrently := limiter.NewConcurrencyLimiterForIO(limiter.DefaultConcurrencyLimitIO)
		concurrently.Execute(func() {
			atomic.AddInt32(&completed, 1)
			// Do some really slow IO ...
			// keep the error:
			concurrently.FirstErrorStore(a)
		})
		concurrently.Execute(func() {
			atomic.AddInt32(&completed, 1)
			// Do some really slow IO ...
			// keep the error:
			concurrently.FirstErrorStore(b)
		})
		concurrently.WaitAndClose()

		Expect(completed).To(BeEquivalentTo(2))
		firstErr := concurrently.FirstErrorGet()
		Expect(firstErr).ToNot(BeNil())
		Expect(firstErr == a || firstErr == b).To(BeTrue())
	})
}

func TestEmpty(t *testing.T) {
	RegisterTestingT(t)

	t.Run("TestEmpty", func(*testing.T) {
		c := limiter.NewConcurrencyLimiter(10)
		c.WaitAndClose()
	})
}
