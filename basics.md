# Go (Golang) Fundamentals – Beginner Guide

This guide explains:

- Hello World syntax breakdown
- Go vs Python analogy
- Variables & types
- Control flow
- Functions
- Structs
- Memory concepts
- Pointers
- Arrays & Slices
- Concurrency basics

---

# 1️⃣ Understanding Your First Go Program

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

---

## 🔹 package main

Every Go file belongs to a package.

If you want the program to run directly, it must be:

```go
package main
```

### Python Analogy

Python does not require this explicitly.

Similar concept:

```python
if __name__ == "__main__":
```

---

## 🔹 import "fmt"

Imports Go's built-in formatting package.

Used for printing to console.

### Python vs Go

Python:
```python
print("Hello")
```

Go:
```go
fmt.Println("Hello")
```

---

## 🔹 func main()

The entry point of the program.

Program execution always starts from `main()`.

Python equivalent:

```python
if __name__ == "__main__":
```

---

## 🔹 Curly Braces `{}`

Go uses `{}` to define code blocks.

Python uses indentation.

Python:
```python
if x > 5:
    print(x)
```

Go:
```go
if x > 5 {
    fmt.Println(x)
}
```

---

# 2️⃣ Variables in Go

## Method 1 – Explicit Type

```go
var age int = 25
```

## Method 2 – Type Inference (Most Common)

```go
age := 25
```

### Python Equivalent

Python:
```python
age = 25
```

Go:
```go
age := 25
```

---

## Static vs Dynamic Typing

Python:
- Dynamic typing
- Type decided at runtime

Go:
- Static typing
- Type decided at compile time

---

# 3️⃣ Basic Data Types

| Type | Example |
|------|---------|
| int | 10 |
| float64 | 3.14 |
| string | "hello" |
| bool | true / false |

Example:

```go
name := "Jay"
age := 30
height := 5.9
isDev := true
```

---

# 4️⃣ Control Flow

## If Statement

```go
if age > 18 {
    fmt.Println("Adult")
}
```

Python:
```python
if age > 18:
    print("Adult")
```

---

## For Loop (Only Loop in Go)

```go
for i := 0; i < 5; i++ {
    fmt.Println(i)
}
```

Python:
```python
for i in range(5):
    print(i)
```

Go does not have a separate while loop.

---

# 5️⃣ Functions

```go
func add(a int, b int) int {
    return a + b
}
```

Python:
```python
def add(a, b):
    return a + b
```

Go requires:
- Parameter types
- Return type

---

# 6️⃣ Structs (Like Classes but Simpler)

Go uses structs instead of traditional classes.

```go
type Person struct {
    Name string
    Age  int
}
```

Python equivalent:

```python
class Person:
    def __init__(self, name, age):
        self.name = name
        self.age = age
```

---

# 7️⃣ Memory Concepts in Go

Go manages memory using:

- Stack
- Heap
- Garbage Collector
- Pointers

---

## Stack

Fast memory for local variables.

```go
x := 10
```

Stored on stack (usually).

---

## Heap

Used when memory needs to survive beyond function scope.

Go automatically manages this.

---

# 8️⃣ Pointers

Pointer stores the address of a variable.

```go
x := 10
p := &x
```

- `&x` → address of x
- `*p` → value stored at that address

Python hides pointer concepts.
Go exposes them clearly.

---

# 9️⃣ Arrays and Slices

## Array (Fixed Size)

```go
var arr [3]int
```

## Slice (Dynamic – Most Used)

```go
nums := []int{1,2,3}
```

Python:
```python
nums = [1,2,3]
```

Slice is closer to Python list.

---

# 🔟 Maps (Like Python Dictionary)

```go
m := map[string]int{
    "Jay": 30,
}
```

Python:
```python
m = {"Jay": 30}
```

---

# 1️⃣1️⃣ Error Handling

Go does not use exceptions like Python.

Functions return errors explicitly.

```go
result, err := someFunction()
if err != nil {
    fmt.Println("Error:", err)
}
```

Python:
```python
try:
    something()
except Exception as e:
    print(e)
```

---

# 1️⃣2️⃣ Concurrency – Go's Superpower

Go uses goroutines.

```go
go someFunction()
```

Goroutine = lightweight thread.

Python equivalent (complex):

```python
import threading
```

Go makes concurrency simple and built-in.

---

# 1️⃣3️⃣ Core Differences: Python vs Go

| Concept | Python | Go |
|----------|---------|-----|
| Typing | Dynamic | Static |
| Indentation | Required | Not required |
| Speed | Interpreted | Compiled |
| Concurrency | Complex | Built-in |
| Classes | Yes | Structs |
| Pointers | Hidden | Explicit |
| Error Handling | Exceptions | Explicit return |

---

# 1️⃣4️⃣ Mental Model of Go

Think of Go as:

Python simplicity  
+  
C-level performance  
+  
Built-in concurrency  
+  
Strict structure  

---

# 1️⃣5️⃣ Recommended Learning Order

1. Variables & Types
2. Functions
3. Structs & Methods
4. Pointers
5. Slices & Maps
6. Error Handling
7. Interfaces
8. Goroutines & Channels
9. Building APIs

---

# 🎉 You Now Understand Go Fundamentals

You are ready to start building:

- CLI tools
- REST APIs
- Microservices
- Concurrent programs

Happy Coding 🚀
