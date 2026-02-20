package main

import (
	"fmt"
	"math"
)

func greet() {
	fmt.Println("Hello from fucntion")
}

func add(a int, b int) int {
	return a + b
}

// 3️⃣ Multiple returns
func divide(a float64, b float64) (float64, float64) {
	return a / b, math.Mod(a, b)
}

// 4️⃣ Named return
func rectangle(w, h float64) (area float64) {
	area = w * h
	return
}

// 5️⃣ Variadic
func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func main() {
	greet()
	result := add(10, 34)
	fmt.Println("Result of addition:", result)
	q, r := divide(10, 3)
	fmt.Println("Divide:", q, r)

	fmt.Println("Area:", rectangle(5, 4))

	fmt.Println("Sum:", sum(1, 2, 3, 4))

	// 6️⃣ Anonymous function
	multiply := func(a int, b int) int {
		return a * b
	}

	fmt.Println("Multiply:", multiply(3, 4))

	// 7️⃣ Closure
	counter := func() func() int {
		count := 0
		return func() int {
			count++
			return count
		}
	}

	inc := counter()
	fmt.Println("Closure:", inc(), inc(), inc())
}
