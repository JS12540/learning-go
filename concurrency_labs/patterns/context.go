package patterns

import (
	"context"
	"fmt"
	"time"
)

/*
context.Context:

Used for:
✔ Cancellation
✔ Deadlines
✔ Timeouts
✔ Propagating request-scoped values
*/

func DemoContext() {

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	ch := make(chan string)

	go func() {
		time.Sleep(500 * time.Millisecond)
		ch <- "Finished work"
	}()

	select {
	case <-ctx.Done():
		fmt.Println("Cancelled:", ctx.Err())

	case msg := <-ch:
		fmt.Println(msg)
	}
}
