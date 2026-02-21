package patterns

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
Atomic Operations:

✔ Lock-free
✔ Very fast
✔ Good for counters / flags

Avoid for complex shared state.
*/

func DemoAtomic() {

	var counter int64
	var wg sync.WaitGroup

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()

	fmt.Println("Atomic Counter:", counter)
}
