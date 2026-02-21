package patterns

import (
	"fmt"
	"sync"
)

/*
Race Condition:

Multiple goroutines modifying shared variable → undefined behavior.

Mutex (Mutual Exclusion Lock):
✔ Only ONE goroutine enters critical section
✔ Others wait

Critical Section = Code accessing shared memory
*/

func DemoMutex() {

	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			mu.Lock()   // acquire lock
			counter++   // critical section
			mu.Unlock() // release lock
		}()
	}

	wg.Wait()

	fmt.Println("Final Counter:", counter)
}
