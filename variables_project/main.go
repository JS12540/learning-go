package main

import (
	"fmt"
	"os"
	"reflect"
)

func main() {
	// 1️⃣ Explicit declaration
	var name string = "Jay"
	var age int = 25
	var height float64 = 5.9
	var isDev bool = true

	// No new line
	fmt.Print("Hello ")
	fmt.Print("World")

	// New line
	fmt.Println("Explicit:", name, age, height, isDev)

	// Formatted Print
	fmt.Printf("name type: %T\n", name)
	fmt.Printf("age type: %T\n", age)
	fmt.Printf("height type: %T\n", height)
	fmt.Printf("isDev type: %T\n", isDev)
	fmt.Printf("Value: %v, Type: %T\n", name, name)
	fmt.Println(reflect.TypeOf(name))

	var city = "Mumbai"
	var score = 99.5

	fmt.Println("Inferred:", city, score)

	// 3️⃣ Short declaration (most common)
	country := "India"
	version := 1

	fmt.Println("Short:", country, version)

	var a, b, c int = 1, 2, 3
	x, y := 10, "Hello"

	fmt.Println("Multiple:", a, b, c)
	fmt.Println("Mixed:", x, y)

	// 5️⃣ Constants
	const pi = 3.14159
	const appName string = "GoLab"

	fmt.Println("Constants:", pi, appName)

	// fmt.Sprintf() — Format WITHOUT Printing
	msg := fmt.Sprintf("Name: %s, Age: %d", name, age)
	fmt.Println(msg)

	// Print to Writer (Advanced)
	fmt.Fprintln(os.Stdout, "Hello to stdout")
	fmt.Fprintln(os.Stderr, "Error message")
}
