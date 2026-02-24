package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/korovkin/limiter"
)

type TaskState int

const (
	kPending TaskState = iota
	kRunning
	kFinished
)

// Configuration constants for the TUI example
const (
	kNumTasks         = 40
	kConcurrencyLimit = 5
	kSeparatorLength  = 82
	kUIRefreshRate    = 100 * time.Millisecond
	kGridWidth        = 4
	kBaseTaskDuration = 1000 * time.Millisecond
	kTaskJitter       = 500 * time.Millisecond
	kFinalWait        = 200 * time.Millisecond
)

type Task struct {
	ID    int
	State TaskState
}

func main() {
	limit := limiter.NewConcurrencyLimiter(kConcurrencyLimit)
	tasks := make([]*Task, kNumTasks)
	for i := 0; i < kNumTasks; i++ {
		tasks[i] = &Task{ID: i + 1, State: kPending}
	}

	sep_line := strings.Repeat("-", kSeparatorLength)

	var mu sync.Mutex
	updateUI := func() {
		mu.Lock()
		defer mu.Unlock()

		// Move cursor to top (or clear screen)
		fmt.Print("\033[H\033[2J")
		fmt.Printf("Concurrency Limiter Visualization (Limit: %d, Active: %d)\n", kConcurrencyLimit, limit.GetNumInProgress())
		fmt.Println(sep_line)

		for _, t := range tasks {
			status := ""
			switch t.State {
			case kPending:
				status = "\033[33m[WAITING ]\033[0m"
			case kRunning:
				status = "\033[32m[RUNNING ]\033[0m"
			case kFinished:
				status = "\033[34m[FINISHED]\033[0m"
			}
			fmt.Printf("Task %2d: %s  ", t.ID, status)
			if (t.ID % kGridWidth) == 0 {
				fmt.Println()
			}
		}
		fmt.Println("\n" + sep_line)
	}

	// UI update loop
	stopUI := make(chan struct{})
	go func() {
		ticker := time.NewTicker(kUIRefreshRate)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				updateUI()
			case <-stopUI:
				updateUI() // final update
				return
			}
		}
	}()

	var wg sync.WaitGroup
	wg.Add(kNumTasks)

	for i := 0; i < kNumTasks; i++ {
		t := tasks[i]
		go func() {
			defer wg.Done()
			// Each goroutine starts "Pending" and then calls Execute which blocks
			// until a concurrency slot (ticket) is available.
			_, _ = limit.Execute(func() {
				mu.Lock()
				t.State = kRunning
				mu.Unlock()

				// simulate variable work time
				time.Sleep(kBaseTaskDuration + time.Duration(t.ID%5)*kTaskJitter)

				mu.Lock()
				t.State = kFinished
				mu.Unlock()
			})
		}()
	}

	wg.Wait()
	limit.WaitAndClose()
	close(stopUI)
	time.Sleep(kFinalWait) // wait for final UI update
	fmt.Println("\n\n=> All tasks completed!\n\n")
}
