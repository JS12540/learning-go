package patterns

import "fmt"

/*
Channels = Communication pipes between goroutines

Unbuffered Channel:
- Sender blocks until receiver ready
- Receiver blocks until sender ready
- Synchronization happens automatically

Buffered Channel:
- Sender blocks only when buffer full
*/

func DemoChannels() {

	// Unbuffered channel
	ch := make(chan string)

	go func() {
		/*
			This send operation BLOCKS
			until someone receives from channel.
		*/
		ch <- "Message from goroutine"
	}()

	// Receive blocks until value available
	msg := <-ch
	fmt.Println(msg)

	// Buffered channel with capacity 2
	buffered := make(chan int, 2)

	buffered <- 1 // doesn't block
	buffered <- 2 // doesn't block

	// buffered <- 3 → WOULD BLOCK (buffer full)

	fmt.Println(<-buffered, <-buffered)
}
