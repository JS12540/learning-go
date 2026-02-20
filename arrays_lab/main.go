package main

import "fmt"

func main() {

	// 1️⃣ Static Array (fixed size)
	var numbers [3]int = [3]int{10, 20, 30}
	fmt.Println("Static array:", numbers)

	// Short declaration
	names := [2]string{"Jay", "Go"}
	fmt.Println("Names:", names)

	// Access element
	fmt.Println("First number:", numbers[0])

	// ------------------------------------------------

	// 2️⃣ Slice (dynamic array)
	scores := []int{95, 88, 76}
	fmt.Println("Slice:", scores)

	// Append elements
	scores = append(scores, 100)
	fmt.Println("After append:", scores)

	// ------------------------------------------------

	// 3️⃣ make() → create slice
	ages := make([]int, 3) // length = 3
	ages[0] = 25
	ages[1] = 30
	ages[2] = 35

	fmt.Println("make slice:", ages)
	fmt.Println("Length:", len(ages))
	fmt.Println("Capacity:", cap(ages))

	// ------------------------------------------------

	// 4️⃣ Looping styles

	// Classic loop
	for i := 0; i < len(scores); i++ {
		fmt.Println("Classic loop:", scores[i])
	}

	// Range loop (preferred)
	for index, value := range scores {
		fmt.Println("Range loop:", index, value)
	}

	// Ignore index
	for _, value := range scores {
		fmt.Println("Only value:", value)
	}

	// ------------------------------------------------

	// 5️⃣ Slice from array
	arr := [5]int{1, 2, 3, 4, 5}
	subSlice := arr[1:4]

	fmt.Println("Array:", arr)
	fmt.Println("Sub-slice:", subSlice)

	// ------------------------------------------------

	// 6️⃣ Copy slice
	src := []int{1, 2, 3}
	dst := make([]int, len(src))

	copy(dst, src)

	fmt.Println("Source:", src)
	fmt.Println("Destination:", dst)

	// ------------------------------------------------

	// 7️⃣ Length vs Capacity demo
	data := make([]int, 0, 5)

	data = append(data, 10)
	data = append(data, 20)

	fmt.Println("Data:", data)
	fmt.Println("Length:", len(data))
	fmt.Println("Capacity:", cap(data))
}
