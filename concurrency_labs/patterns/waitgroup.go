package patterns

import (
	"fmt"
	"sync"
)

/*
sync.WaitGroup:

Purpose → Wait for multiple goroutines to finish.

Mechanism:
✔ Add(n) before launching goroutines
✔ Done() inside goroutine
✔ Wait() blocks until counter = 0
*/

func DemoWaitGroups() {

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {

		wg.Add(1) // increment counter

		go func(id int) {
			defer wg.Done() // decrement counter

			fmt.Println("Worker", id, "done")
		}(i)
	}

	// Block until all Done() called
	wg.Wait()

	fmt.Println("All workers finished")
}
