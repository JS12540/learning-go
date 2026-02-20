package main

import "fmt"

func main() {
	// 1️⃣ Classic for loop
	for i := 0; i < 5; i++ {
		fmt.Println("Classic:", i)
	}

	// 2️⃣ While-style loop
	j := 0
	for j < 5 {
		fmt.Println("While-style:", j)
		j++
	}

	// 3️⃣ Infinite loop
	k := 0
	for {
		fmt.Println("Infinite loop:", k)
		k++
		if k == 3 {
			break
		}
	}

	// 4️⃣ Range loop (arrays/slices)
	nums := []int{10, 20, 30}

	for index, value := range nums {
		fmt.Println("Range:", index, value)
	}

	// 5️⃣ Range ignoring index
	for _, value := range nums {
		fmt.Println("Only value:", value)
	}

	// 6️⃣ Loop over string
	for i, ch := range "GoLang" {
		fmt.Println("Char:", i, string(ch))
	}
}
