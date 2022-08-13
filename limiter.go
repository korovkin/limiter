package limiter

func BoundedConcurrencyWithDoneProcessor[C any](
	limit,
	jobs int,
	job func(int) C,
	done func(C),
) {
	inputChan := make(chan int, limit)
	defer close(inputChan)
	outputChan := make(chan C)
	defer close(outputChan)
	for workerNumber := 0; workerNumber < limit; workerNumber++ {
		go func() {
			for i := range inputChan {
				result := job(i)
				outputChan <- result
			}
		}()
	}

	go func() {
		for i := 0; i < jobs; i++ {
			inputChan <- i
		}
	}()

	for i := 0; i < jobs; i++ {
		data := <-outputChan
		done(data)
	}
}

func BoundedConcurrency(
	limit,
	jobs int,
	job func(int),
) {
	BoundedConcurrencyWithDoneProcessor(
		limit,
		jobs,
		func(i int) int {
			job(i)
			return i
		},
		func(int) {},
	)
}
