package main

import "fmt"

func Counter() {

	// counter is a function that RETURNS another function
	createCounter := func() func() int {

		// This variable belongs to the outer function
		// but will be "captured" by the inner function
		count := 0

		// Return an anonymous function (closure)
		return func() int {

			// This function remembers the 'count' variable
			// even after counter() has finished execution
			count++

			return count
		}
	}

	// Call counter() → returns the inner function
	inc := createCounter()

	// Each call updates the SAME captured variable
	fmt.Println("Closure call 1:", inc()) // 1
	fmt.Println("Closure call 2:", inc()) // 2
	fmt.Println("Closure call 3:", inc()) // 3
}
