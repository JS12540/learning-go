package pipeline

import "fmt"

/*
Pipeline Pattern:

✔ Stages connected by channels
✔ Each stage transforms data
✔ Enables streaming
*/

func DemoPipeline() {

	gen := func(nums ...int) <-chan int {
		out := make(chan int)

		go func() {
			for _, n := range nums {
				out <- n
			}
			close(out)
		}()

		return out
	}

	sq := func(in <-chan int) <-chan int {
		out := make(chan int)

		go func() {
			for n := range in {
				out <- n * n
			}
			close(out)
		}()

		return out
	}

	for n := range sq(gen(1, 2, 3)) {
		fmt.Println(n)
	}
}
