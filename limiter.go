package limiter

func BoundedConcurrency(limit, jobs int, job func(i int)) {
	inputChan := make(chan int, limit)
	defer close(inputChan)
	outputChan := make(chan int)
	defer close(outputChan)
	for workerNumber := 0; workerNumber < limit; workerNumber++ {
		go func() {
			for i := range inputChan {
				job(i)
				outputChan <- i
			}
		}()
	}

	go func() {
		for i := 0; i < jobs; i++ {
			inputChan <- i
		}
	}()

	for i := 0; i < jobs; i++ {
		<-outputChan
	}
}
