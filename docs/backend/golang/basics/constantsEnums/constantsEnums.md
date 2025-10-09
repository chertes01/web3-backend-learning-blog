# Go Constants and Enums — Study Notes

This note summarizes the various ways to declare constants (`const`), grouped declarations, type inference, using `iota` for enums, custom types with method bindings, pointer receivers, and common pitfalls. Suitable as a GitHub README reference.


## Constant Declaration Methods and Rules

### 1. Explicitly Typed Constants

```go
const a int = 1
```
- The type of constant `a` is explicitly declared as `int`.
- Constants are determined at compile time and cannot be modified.

---

### 2. Untyped Constants

```go
const b = "test"
```
- Constants in Go can be untyped, known as **untyped constants**.
- The compiler infers the type based on the context.

Example:
```go
var s string = b  // Inferred as string
var r rune = b[0] // Inferred as rune
```

---

### 3. Multiple Constants in One Line

```go
const c, d = 2, "hello"
```
- Multiple constants can be declared in one line.
- Types are inferred from the right-hand side.

---

### 4. Explicit Types with Multiple Assignments

```go
const e, f bool = true, false
```
- Multiple constants can have explicit types, but all must share the same type.

---

### 5. Grouped Declarations

```go
const (
    h    byte = 3
    i         = "value"
    j, k      = "v", 4
    l, m      = 5, false
)
```
- Grouped declarations improve readability and organization.
- Within a group:
  - `iota` automatically increments.
  - Omitting the type or expression reuses the previous line's type and value.

---

### 6. Independent Groups with `iota`

```go
const (
    n = 6
)
```
- Each `const (...)` group is independent.
- `iota` resets to 0 in each new group.

---

## Named Types and Enums

### 1. Defining a New Type

```go
type Gender string
```
- This creates a new named type, not an alias.
- **Difference**:
  - `type Gender = string` creates an alias.
  - `type Gender string` defines a new type, allowing method bindings.

---

### 2. Enum Constants

```go
const (
    Male   Gender = "Male"
    Female Gender = "Female"
)
```
- This creates an enum-like structure using `const` and a custom type.
- The constants are fixed values restricted to the `Gender` type.

---

### 3. Binding Methods to a Type

```go
func (g *Gender) IsMale() bool {
    return *g == Male
}

func (g *Gender) IsFemale() bool {
    return *g == Female
}
```
- Methods can be bound to the `Gender` type.
- **Pointer Receiver**:
  - If the receiver is a variable (addressable), the compiler automatically takes its address.
  - Constants or literals cannot directly call pointer methods (e.g., `Male.IsMale()` ❌, but `var g = Male; g.IsMale()` ✅).
  - **Caution**: If `g` is a `nil` pointer, it will cause a runtime panic.

---

## Using `iota` for Enums

### 1. Defining a Month Enum

```go
type Month int

const (
    January Month = 1 + iota
    February
    March
    April
    May
    June
    July
    August
    September
    October
    November
    December
)
```
- **How `iota` Works**:
  - `iota` is a constant counter that starts at 0 in each `const` group and increments by 1 for each new line.
  - If an expression is omitted, it reuses the previous line's expression (but `iota` still increments).

Example:
```go
January = 1 + 0 = 1
February = 1 + 1 = 2
...
December = 1 + 11 = 12
```

---

### 2. Common `iota` Patterns

```go
const (
    _  = iota             // Skip 0
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
    TB                    // 1 << 40
)
```

---

## Pointer Receivers and Auto-Addressing

### Pointer Receiver Example

```go
func (g *Gender) IsMale() bool {
    return *g == Male
}
```

### Auto-Addressing Rules

- If the receiver is an addressable variable, the compiler automatically takes its address:
  ```go
  var g = Male
  g.IsMale() // Valid
  ```
- Constants or literals cannot directly call pointer methods:
  ```go
  Male.IsMale() // Compilation error
  ```
- If the pointer is `nil`, calling a pointer method will cause a runtime panic.

---

## Common Pitfalls

| Issue                  | Explanation                                                                 |
|------------------------|-----------------------------------------------------------------------------|
| `iota` Reset           | Each `const (...)` block resets `iota` to 0.                               |
| Pointer Receiver Method| Literals cannot call pointer methods; only variables can.                  |
| Untyped Constants      | Flexible but inferred type depends on the context during assignment.       |
| Constant Expressions   | Must be computable at compile time (cannot depend on runtime variables).   |
| Grouped Declaration    | Omitting values reuses the previous line's expression, but `iota` increments. |

---

## Recommended Additions: String Method and `iota` Patterns

### 1. Adding a `String()` Method for Enums

```go
func (m Month) String() string {
    names := [...]string{
        "January", "February", "March", "April", "May", "June",
        "July", "August", "September", "October", "November", "December",
    }
    if m < January || m > December {
        return "Unknown"
    }
    return names[m-1]
}
```
- Example:
  ```go
  fmt.Println(September) // Outputs "September"
  ```

---

### 2. Typical `iota` Patterns

```go
const (
    _  = iota             // Skip 0
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
)
```

---

### 3. Benefits of Enum Types

- **Improved Readability**: Avoids using plain numbers or strings.
- **Type Safety**: Compiler checks prevent mixing types.
- **Extensibility**: Methods (e.g., validation, printing, comparison) can be added.

---

## Summary of Key Points

| Category              | Details                                                                 |
|-----------------------|-------------------------------------------------------------------------|
| Constant Declaration  | `const name [type] = value`, supports grouping.                        |
| Type Inference        | Untyped constants are inferred based on context.                      |
| Enums                 | Use `iota` for auto-incrementing constants.                           |
| Named Types           | `type NewType BaseType` defines a new type.                           |
| Pointer Receiver      | Allows auto-addressing, but beware of `nil` pointers.                 |
| Constant Rules        | Only compile-time computable values are allowed.                      |
| Enum String Output    | Implement `String()` for better readability.                          |

---