package worker

import (
	"fmt"
	"sync"
	"time"
)

/*
Worker Pool Pattern:

Purpose:
✔ Limit concurrency
✔ Control resource usage
✔ Prevent goroutine explosion
*/

func DemoWorkerPool() {

	jobs := make(chan int, 5)
	var wg sync.WaitGroup

	for w := 1; w <= 3; w++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for job := range jobs {
				fmt.Printf("Worker %d processing job %d\n", id, job)
				time.Sleep(100 * time.Millisecond)
			}
		}(w)
	}

	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	wg.Wait()
}
