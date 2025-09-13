# Go Conditional Statements — Study Notes

A brief overview of common conditional control structures in Go: usage, caveats, and examples for `if` and `switch`. Suitable as a GitHub README reference.

---

## Table of Contents

- [if statement](#if-statement)
  - Basic usage
  - Init statement and scope
  - Example: checkNumber
- [switch statement](#switch-statement)
  - Value matching (basic usage)
  - With init statement
  - Expressionless switch (condition checks)
  - Type switch
  - Combined example: analyzeInput
- [Summary & Recommendations](#summary--recommendations)

---

## if statement

### Basic usage
```go
if condition {
    // do something
} else if anotherCondition {
    // do another
} else {
    // default case
}
```

### Init statement and scope
Go allows an initialization statement before the `if` condition. Variables declared here are only valid within the `if`-`else` chain:
```go
if x := 10; x > 5 {
    fmt.Println("x > 5")
} else {
    fmt.Println("x <= 5")
}
// x is out of scope here
```

Example (variable scope):
```go
func main() {
    var a int = 10
    if b := 1; a > 10 {
        b = 2
        fmt.Println("a > 10")
    } else if c := 3; b > 1 {  // b comes from if-init, c is only valid in this else-if
        b = 3
        fmt.Println("b > 1")
    } else {
        fmt.Println("other")
        fmt.Println(b) // b is still valid here
        // fmt.Println(c) // compile error: c is out of scope
    }
}
```

Tip: If you need to access the same variable in multiple branches, declare it outside (see `flag`, `m` in the examples).

### Example: checkNumber
```go
func checkNumber(n int) {
    flag := false
    m := -1

    if n > 100 {
        flag = true
        fmt.Println("n is a big number")
    } else if m = n % 2; m == 0 {
        fmt.Println("n is even")
    } else {
        fmt.Println("n is odd")
        if flag {
            fmt.Println("but flag is true")
        }
        fmt.Println("m is", m)
    }

    fmt.Println("flag is", flag)
    fmt.Println("----")
}
```

---

## switch statement

### Value matching (basic usage)
```go
a := "test string"

switch a {
case "test":
    fmt.Println("a = test")
case "t", "test string":
    fmt.Println("a = test string or t")
default:
    fmt.Println("default case")
}
```

### With init statement
```go
switch b := 5; b {
case 1:
    fmt.Println("b = 1")
case 3, 4:
    fmt.Println("b = 3 or 4")
case 5:
    fmt.Println("b = 5")
default:
    fmt.Println("b =", b)
}
```

### Expressionless switch (condition checks)
```go
b := 5
switch {
case b == 3:
    fmt.Println("b = 3")
case b == 5:
    fmt.Println("b = 5")
default:
    fmt.Println("default case")
}
```

### Type switch
Used to determine the concrete type of an `interface{}`:
```go
var d interface{}
d = 1

switch t := d.(type) {
case int:
    fmt.Println("d is int:", t)
case string:
    fmt.Println("d is string:", t)
case float64:
    fmt.Println("d is float64:", t)
default:
    fmt.Println("unknown type:", t)
}
```

### Combined example: analyzeInput
```go
func analyzeInput(x interface{}) {
    switch v := x.(type) {
    case string:
        switch v {
        case "hello":
            fmt.Println("got hello")
        case "world", "hi":
            fmt.Println("got greeting")
        default:
            fmt.Println("Please provide a greeting")
        }

        length := len(v)
        switch {
        case length > 5:
            fmt.Println("long string")
        case length%2 == 1:
            fmt.Println("odd length")
        case length%2 == 0:
            fmt.Println("even length")
        }
        fmt.Println("string", v)

    case int:
        switch {
        case v == 0:
            fmt.Println("zero")
        case v == 1:
            fmt.Println("one")
        case v > 10:
            fmt.Println("n is big number")
        default:
            fmt.Println("n is a number")
        }

        if v%2 == 0 {
            fmt.Println("even number")
        } else {
            fmt.Println("odd number")
        }
        fmt.Println("int", v)

    case float64:
        fmt.Println("float64", v)

    default:
        fmt.Println("unknown type:", v)
    }
    fmt.Println("----")
}
```

Example calls (omitted) will output results based on different types and values.

---

## Summary & Recommendations

- `if`: Suitable for complex boolean logic and when you need to initialize temporary variables, but watch the scope of init variables.
- `switch`: More concise, suitable for multi-branch, type branching, or simple conditions.
- General advice:
  - Simple multi-branch → use `switch`.
  - Complex logic → use `if`.

---