package patterns

import (
	"fmt"
	"sync"
)

/*
Fan-Out:
✔ Distribute work to multiple goroutines

Fan-In:
✔ Merge results back into single channel
*/

func DemoFanOutFanIn() {

	in := make(chan int)
	out := make(chan int)

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()

		for n := range in {
			out <- n * 10
		}
	}

	// Fan-Out → multiple workers
	wg.Add(2)
	go worker()
	go worker()

	// Producer
	go func() {
		for i := 1; i <= 5; i++ {
			in <- i
		}
		close(in)

		wg.Wait()
		close(out)
	}()

	// Fan-In → single consumer
	for res := range out {
		fmt.Println(res)
	}
}
