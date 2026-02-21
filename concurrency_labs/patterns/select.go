package patterns

import (
	"fmt"
	"time"
)

/*
select statement:

Waits on multiple channel operations.

Behavior:
✔ Executes first available case
✔ If multiple ready → random selection
✔ Optional default prevents blocking
*/

func DemoSelect() {

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "from ch1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "from ch2"
	}()

	select {
	case msg := <-ch1:
		fmt.Println("Received:", msg)

	case msg := <-ch2:
		fmt.Println("Received:", msg)

		/*
			default:
				fmt.Println("No channel ready")

			Default makes select NON-BLOCKING
			Useful in polling / timeouts
		*/
	}
}
