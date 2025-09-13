# Go Functions, Methods, and Closures — Study Notes

This note summarizes the usage of functions, methods, closures, and immediately-invoked function expressions (IIFE) in Go. Suitable as a GitHub README reference.

---

## 1. Functions

In Go, functions are first-class citizens. They can be declared, assigned, passed, and returned just like variables.

### 1.1 Basic Declaration

```go
func add(a int, b int) int {
    return a + b
}
```

### 1.2 Anonymous Functions

Anonymous functions are functions without a name and can be assigned directly to variables:

```go
square := func(x int) int {
    return x * x
}
fmt.Println(square(5)) // Output: 25
```

---

## 2. Methods

A method is a function bound to a type. Compared to ordinary functions, it has an extra receiver.

### 2.1 Method Declaration

```go
type Counter struct {
    value int
}

// Add method, receiver is *Counter
func (c *Counter) Add(n int) int {
    c.value += n
    return c.value
}
```

Usage:

```go
c := Counter{0}
fmt.Println(c.Add(3)) // Output: 3
fmt.Println(c.Add(5)) // Output: 8
```

### 2.2 Method Value

Methods can be assigned to function variables:

```go
fn := c.Add
fmt.Println(fn(2)) // Equivalent to c.Add(2), Output: 10
```

> Note: `fn` is bound to the current instance `c`. Even if `c` is later changed to a new object, `fn` still points to the original instance.

---

## 3. Closures

Closure = function + the external variable environment it references. Closures can "remember" the context in which they were defined.

### 3.1 Basic Usage of Closures

```go
func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

posSum := adder()
fmt.Println(posSum(1)) // 1
fmt.Println(posSum(2)) // 3
fmt.Println(posSum(3)) // 6
```

### 3.2 Closure Factory

```go
func makeMultiplier(factor int) func(y int) int {
    return func(y int) int {
        return factor * y
    }
}

times10 := makeMultiplier(10)
fmt.Println(times10(5)) // Output: 50

// Multiple independent closure environments
times2 := makeMultiplier(2)
times3 := makeMultiplier(3)
fmt.Println(times2(5)) // 10
fmt.Println(times3(5)) // 15
```

---

## 4. Immediately-Invoked Function Expression (IIFE)

In Go, you can declare an anonymous function and execute it immediately:

```go
returnFunc := func() func(int, string) (int, string) {
    fmt.Println("this is an anonymous function")
    return func(i int, s string) (int, string) {
        return i, s
    }
}() // Immediately invoked

ret1, ret2 := returnFunc(1, "test")
fmt.Println(ret1, ret2) // Output: 1 test
```

---

## 5. Comprehensive Example

Combining functions, methods, and closures:

```go
package main

import "fmt"

type Counter struct {
    value int
}

func (c *Counter) Add(n int) int {
    c.value += n
    return c.value
}

func makeMultiplier(factor int) func(y int) int {
    return func(y int) int {
        return factor * y
    }
}

func main() {
    // Method
    counter := Counter{0}
    fn := counter.Add
    fmt.Println(fn(3)) // 3
    fmt.Println(fn(5)) // 8

    // Function
    square := func(x int) int { return x * x }
    fmt.Println(square(5)) // 25

    // Closure
    times10 := makeMultiplier(10)
    fmt.Println(times10(5)) // 50

    // Immediately-invoked anonymous function
    returnFunc := func() func(int, string) (int, string) {
        fmt.Println("this is an anonymous function")
        return func(i int, s string) (int, string) {
            return i, s
        }
    }()
    ret1, ret2 := returnFunc(1, "GoLang")
    fmt.Println(ret1, ret2) // 1 GoLang
}
```

---

## 📌 Summary

- **Function**: The most basic reusable code block, can be passed as a value.
- **Method**: A function bound to a type, with a receiver.
- **Closure**: A combination of a function and its external variables, can maintain state.
- **Immediately-invoked anonymous function**: Declare and execute a function immediately, often