package patterns

import (
	"fmt"
	"time"
)

/*
DemoBasicGoroutines demonstrates:

1️⃣ How to launch a goroutine
2️⃣ That main() does NOT wait automatically
3️⃣ Why Sleep() is sometimes needed in demos
*/

func DemoBasicGoroutines() {

	// Launching a goroutine using "go" keyword
	// This runs concurrently with main goroutine
	go func() {
		fmt.Println("Hello from goroutine!")

		/*
			This executes independently.
			If main() exits before this runs → it gets killed.
		*/
	}()

	fmt.Println("Hello from main!")

	/*
		Why Sleep?

		main goroutine exits immediately after this function.
		If we don’t sleep → goroutine may never execute.

		In real apps → use sync.WaitGroup instead.
	*/
	time.Sleep(300 * time.Millisecond)
}
