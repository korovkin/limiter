package main

import (
	"fmt"
	"time"

	"github.com/korovkin/limiter"
)

func main() {
	// create a limiter with 3 concurrent workers
	limit := limiter.NewConcurrencyLimiter(5)

	// we want to process 10 items
	for i := 1; i <= 10; i++ {
		item := i
		limit.Execute(func() {
			fmt.Printf("Processing item %d (Active: %d)\n", item, limit.GetNumInProgress())
			// simulate some work
			time.Sleep(100 * time.Millisecond)
			fmt.Printf("Finished item %d\n", item)
		})
	}

	// wait for all jobs to finish
	limit.WaitAndClose()
	fmt.Println("All items processed")
}
