# Go Pointers and Unsafe Study Notes

This note summarizes the basic usage of pointers in Go, pointer to pointer, `new`, the `unsafe` package, and related low-level operations. It is suitable for both beginners and advanced readers.

---

## Example Code

```go
package main

import (
    "fmt"
    "unsafe"
)

var p *int
var q *string

func main() {
    i := 10
    p = &i

    q = new(string)
    *q = "hello"

    q2 := &q

    println(p, q, q2)
    println(*p == i, *q == "hello", *q2 == q)

    *p = 20
    *q = "world"
    println(i, *q)

    p2 := &p
    **p2 = 50
    println(*p2, p, **p2, i)

    a := "Hello, world!"
    upA := uintptr(unsafe.Pointer(&a))
    upA += 1

    c := (*uint8)(unsafe.Pointer(upA))
    fmt.Println(*c)
}
```

---

## 1. Packages and Imports

- `package main`: Entry package for executable programs.
- `import ("fmt" "unsafe")`:
  - `fmt`: For formatted output (recommended for production code).
  - `unsafe`: Allows operations that bypass the type system and memory safety (use with caution).

---

## 2. Package-level Variables

```go
var p *int
var q *string
```
- `p` is of type `*int` (pointer to int), `q` is `*string`.
- Package-level variables are zero-initialized (the zero value for pointers is `nil`).

---

## 3. Local Variables and Address Taking

```go
i := 10
p = &i
```
- `i` is a local int variable, `&i` takes its address and assigns it to `p`, so `*p` can read/write `i`.

---

## 4. new and String Pointer

```go
q = new(string)
*q = "hello"
```
- `new(string)` allocates a zero value of type string and returns its pointer.
- `*q = "hello"` writes a string to the location pointed to by the pointer.

---

## 5. Pointer to Pointer

```go
q2 := &q
```
- `q2` is of type `**string`, i.e., a pointer to a string pointer.
- `*q2 == q`, `**q2 == *q`.

---

## 6. Modifying Original Value via Pointer

```go
*p = 20
*q = "world"
```
- `*p = 20` is equivalent to `i = 20`.
- `*q = "world"` modifies the string pointed to by q.

---

## 7. Pointer to Pointer and Double Dereference

```go
p2 := &p
**p2 = 50
```
- `p2` is of type `**int`, `*p2` is `*int`, `**p2` is `int`.
- `**p2 = 50` is equivalent to `i = 50`.

---

## 8. unsafe, uintptr, and Memory Offset (Dangerous Operation)

```go
a := "Hello, world!"
upA := uintptr(unsafe.Pointer(&a))
upA += 1
c := (*uint8)(unsafe.Pointer(upA))
fmt.Println(*c)
```
- Go strings are implemented as a struct (data pointer + length) under the hood.
- `unsafe.Pointer(&a)` gets the address of the string header, converts it to an integer for address arithmetic.
- After `upA += 1`, converting back to a pointer and dereferencing reads the first byte inside the string header (not the string content).
- **Risk**: This operation depends on platform and memory layout, is non-portable and unsafe, and may cause crashes or undefined behavior.

---

## Summary

- Pointers are used for indirect access and modification of variables, and support multiple levels of indirection.
- `new(T)` returns a pointer to type T, `&x` takes the address of a variable.
- The `unsafe` package enables low-level operations, but should be used with caution to avoid breaking memory safety.
- It is recommended to use the `fmt` family of functions for output, and avoid using `println`.
- Mastering pointer operations is fundamental for advanced Go development and helps you write more efficient