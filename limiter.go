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

// enforce a maximum Concurrency of limit
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

// if num of go routines allocated by this instance is < limit
// launch a new go routine to execut job
// else wait until a go routine becomes available
func (c *ConcurrencyLimiter) Execute(job func()) int {
	ticket := <-c.tickets
	atomic.AddInt32(&c.numInProgress, 1)
	go func() {
		defer func() {
			c.tickets <- ticket
			atomic.AddInt32(&c.numInProgress, -1)

		}()

		// run the job
		job()
	}()
	return ticket
}

// if num of go routines allocated by this instance is < limit
// launch a new go routine to execut job
// else wait until a go routine becomes available
func (c *ConcurrencyLimiter) ExecuteWithTicket(job func(ticket int)) int {
	ticket := <-c.tickets
	atomic.AddInt32(&c.numInProgress, 1)
	go func() {
		defer func() {
			c.tickets <- ticket
			atomic.AddInt32(&c.numInProgress, -1)
		}()

		// run the job
		job(ticket)
	}()
	return ticket
}

// wait until all the previously Executed jobs completed running
//
// IMPORTANT: calling the Wait function while keep calling Execute leads to
//            un-desired race conditions
func (c *ConcurrencyLimiter) Wait() {
	for i := 0; i < c.limit; i++ {
		_ = <-c.tickets
	}
}

// get a racy counter of how many go routines are active right now
func (c *ConcurrencyLimiter) GetNumInProgress() int32 {
	return c.numInProgress
}
