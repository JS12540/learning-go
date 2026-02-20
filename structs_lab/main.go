package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func (p Person) greet() {
	fmt.Println("Hello, my name is", p.Name)
}

func greet_new(p Person) {
	fmt.Println("Hello, my name is", p.Name)
}

func (p *Person) birthday() {
	p.Age++
}

func birthday_change(p *Person) {
	p.Age++
}

func main() {
	fmt.Println("Hello")

	p1 := Person{Name: "Jay", Age: 25}

	fmt.Println(p1)

	p1.greet()

	greet_new(p1)

	p1.birthday()
	fmt.Println("After birthday:", p1.Age)

	birthday_change(&p1)
	fmt.Println("After birthday:", p1.Age)

}

// Polymorphism via Interface

type Animal interface {
	Speak() string
}

type Dog struct{}

func (d Dog) Speak() string { return "Woof" }

type Cat struct{}

func (c Cat) Speak() string { return "Meow" }

func makeSound(a Animal) {
	fmt.Println(a.Speak())
}
