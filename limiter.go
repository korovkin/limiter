package limiter

import (
	"sync/atomic"
)

const (
	DEFAULT_LIMIT = 100
)

type ConcurrencyLimiter struct {
	limit         int      `json:"limit"`
	tickets       chan int `json:"tickets"`
	numInProgress int32    `json:"in_progress"`
}

func NewConcurrencyLimiter(limit int) *ConcurrencyLimiter {
	if limit <= 0 {
		limit = DEFAULT_LIMIT
	}

	// allocate a limiter instance
	c := &ConcurrencyLimiter{
		limit:   limit,
		tickets: make(chan int, limit),
	}

	// allocate the tickets:
	for i := 0; i < c.limit; i++ {
		c.tickets <- i
	}

	return c
}

func (c *ConcurrencyLimiter) Execute(job func()) {
	ticket := <-c.tickets
	atomic.AddInt32(&c.numInProgress, 1)
	go func() {
		defer func() {
			c.tickets <- ticket
			atomic.AddInt32(&c.numInProgress, -1)
		}()

		job()
	}()
}

func (c *ConcurrencyLimiter) GetNumInProgress() int32 {
	return c.numInProgress
}

func (c *ConcurrencyLimiter) Wait() {
	for i := 0; i < c.limit; i++ {
		_ = <-c.tickets
	}
}
