# 🐹 Go: Structs vs OOP

> Go does **not** have classes. Instead, it achieves OOP concepts through **structs**, **methods**, **interfaces**, and **composition**.

---

## 📊 Quick Reference

| OOP Concept     | Go Equivalent              |
|-----------------|----------------------------|
| Class           | `struct`                   |
| Object          | Struct instance            |
| Method          | Function with a receiver   |
| Constructor     | Factory function (`NewX`)  |
| Inheritance     | Composition (embedding)    |
| Polymorphism    | Interfaces                 |
| Encapsulation   | Exported/unexported fields |

---

## 1. 🧱 Struct — The "Class"

A `struct` groups fields together, just like a class holds properties.

```go
package main

import "fmt"

// Define a "class" — just a struct in Go
type Animal struct {
    Name string
    Age  int
    // unexported (private) field — lowercase
    sound string
}

func main() {
    // Create an "object" — a struct instance
    dog := Animal{
        Name:  "Rex",
        Age:   3,
        sound: "Woof",
    }

    fmt.Println(dog.Name) // Rex
    fmt.Println(dog.Age)  // 3
}
```

---

## 2. 🏗️ Constructor — Factory Functions

Go has no `new` keyword for custom types. Use a factory function by convention (`NewX`).

```go
// "Constructor" pattern in Go
func NewAnimal(name string, age int, sound string) *Animal {
    return &Animal{
        Name:  name,
        Age:   age,
        sound: sound,
    }
}

func main() {
    cat := NewAnimal("Whiskers", 5, "Meow")
    fmt.Println(cat.Name) // Whiskers
}
```

---

## 3. ⚙️ Methods — Functions with Receivers

Methods in Go are functions tied to a type via a **receiver**.

```go
type Animal struct {
    Name  string
    sound string
}

// Method with a POINTER receiver (can modify the struct)
func (a *Animal) Speak() string {
    return a.Name + " says: " + a.sound
}

// Method with a VALUE receiver (read-only copy)
func (a Animal) Describe() string {
    return fmt.Sprintf("%s is an animal", a.Name)
}

func main() {
    dog := &Animal{Name: "Rex", sound: "Woof"}
    fmt.Println(dog.Speak())    // Rex says: Woof
    fmt.Println(dog.Describe()) // Rex is an animal
}
```

> 💡 **Pointer receiver** (`*Animal`) — use when you need to modify fields or avoid copying large structs.  
> 💡 **Value receiver** (`Animal`) — use for read-only operations on small structs.

---

## 4. 🔒 Encapsulation — Exported vs Unexported

Go uses **capitalization** to control visibility (no `public`/`private` keywords).

```go
package animals

type Dog struct {
    Name string  // ✅ Exported (public)  — accessible outside the package
    age  int     // ❌ Unexported (private) — only accessible within this package
}

// Getter for private field
func (d *Dog) GetAge() int {
    return d.age
}

// Setter for private field
func (d *Dog) SetAge(age int) {
    if age > 0 {
        d.age = age
    }
}
```

---

## 5. 🧩 Composition — Go's "Inheritance"

Go has **no inheritance**. Instead, you **embed** one struct inside another to reuse behavior.

```go
// Base "class"
type Animal struct {
    Name string
}

func (a Animal) Breathe() string {
    return a.Name + " is breathing"
}

// Dog "inherits" from Animal via embedding
type Dog struct {
    Animal        // Embedded struct — promotes Animal's fields and methods
    Breed  string
}

func (d Dog) Fetch() string {
    return d.Name + " is fetching the ball!" // Access Animal.Name directly
}

func main() {
    dog := Dog{
        Animal: Animal{Name: "Rex"},
        Breed:  "Labrador",
    }

    fmt.Println(dog.Name)      // Rex       — promoted from Animal
    fmt.Println(dog.Breathe()) // Rex is breathing — promoted method
    fmt.Println(dog.Fetch())   // Rex is fetching the ball!
    fmt.Println(dog.Breed)     // Labrador
}
```

> 💡 Embedding promotes fields and methods to the outer struct — it looks like inheritance but is really just **composition**.

---

## 6. 🎭 Interfaces — Polymorphism

An **interface** defines a set of methods. Any type that implements those methods satisfies the interface — **implicitly** (no `implements` keyword).

```go
// Define the interface
type Speaker interface {
    Speak() string
}

// Dog implements Speaker
type Dog struct{ Name string }
func (d Dog) Speak() string { return d.Name + " says: Woof!" }

// Cat implements Speaker
type Cat struct{ Name string }
func (c Cat) Speak() string { return c.Name + " says: Meow!" }

// Parrot implements Speaker
type Parrot struct{ Name string }
func (p Parrot) Speak() string { return p.Name + " says: Squawk!" }

// Polymorphic function — accepts any Speaker
func MakeNoise(s Speaker) {
    fmt.Println(s.Speak())
}

func main() {
    animals := []Speaker{
        Dog{Name: "Rex"},
        Cat{Name: "Whiskers"},
        Parrot{Name: "Polly"},
    }

    for _, animal := range animals {
        MakeNoise(animal) // Polymorphism in action!
    }
}

// Output:
// Rex says: Woof!
// Whiskers says: Meow!
// Polly says: Squawk!
```

> 💡 In Go, interface satisfaction is **implicit** — you never write `class Dog implements Speaker`. If a type has the method, it satisfies the interface.

---

## 7. 🏛️ Putting It All Together

A complete example combining structs, methods, embedding, and interfaces.

```go
package main

import "fmt"

// --- Interface ---
type Mover interface {
    Move() string
}

// --- Base struct ---
type LivingThing struct {
    Name string
    HP   int
}

func (l *LivingThing) TakeDamage(dmg int) {
    l.HP -= dmg
}

func (l LivingThing) IsAlive() bool {
    return l.HP > 0
}

// --- Warrior struct (composes LivingThing) ---
type Warrior struct {
    LivingThing        // Embed (composition over inheritance)
    Weapon      string
}

// Warrior implements Mover
func (w Warrior) Move() string {
    return w.Name + " charges forward with a " + w.Weapon + "!"
}

func (w Warrior) Attack() string {
    return w.Name + " swings their " + w.Weapon
}

// --- Mage struct ---
type Mage struct {
    LivingThing
    Spell string
}

// Mage implements Mover
func (m Mage) Move() string {
    return m.Name + " teleports forward!"
}

func (m Mage) CastSpell() string {
    return m.Name + " casts " + m.Spell + "!"
}

// --- Polymorphic function ---
func AdvanceAll(movers []Mover) {
    for _, m := range movers {
        fmt.Println(m.Move())
    }
}

// --- Main ---
func main() {
    warrior := Warrior{
        LivingThing: LivingThing{Name: "Aragorn", HP: 100},
        Weapon:      "sword",
    }

    mage := Mage{
        LivingThing: LivingThing{Name: "Gandalf", HP: 80},
        Spell:       "Fireball",
    }

    // Methods from embedded struct
    warrior.TakeDamage(20)
    fmt.Printf("%s HP: %d, Alive: %v\n", warrior.Name, warrior.HP, warrior.IsAlive())

    // Type-specific methods
    fmt.Println(warrior.Attack())
    fmt.Println(mage.CastSpell())

    // Polymorphism via interface
    fmt.Println("\n--- All units advance! ---")
    AdvanceAll([]Mover{warrior, mage})
}

// Output:
// Aragorn HP: 80, Alive: true
// Aragorn swings their sword
// Gandalf casts Fireball!
//
// --- All units advance! ---
// Aragorn charges forward with a sword!
// Gandalf teleports forward!
```

---

## 🔑 Key Takeaways

- **No classes** — use `struct` to group data
- **No constructors** — use `NewX()` factory functions by convention
- **No `public`/`private`** — capitalize for exported, lowercase for unexported
- **No inheritance** — use **embedding** (composition) to reuse behavior
- **No `implements`** — interface satisfaction is **implicit**; if it has the methods, it qualifies
- **Pointer vs value receivers** — use pointers when mutating state or working with large structs