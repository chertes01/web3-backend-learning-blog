# Go Arrays and Pointers — Study Notes

This note summarizes the basics, common pitfalls, and important points of arrays and pointers in Go. Suitable as a GitHub README reference.

---

## 1. One-Dimensional Array Basics

### 1.1 Declaring Arrays

```go
var a [5]int             // int array of length 5, elements default to 0
b := [5]int{1,2,3,4,5}   // Type inference declaration
c := [...]int{1,2,3}     // ... lets the compiler infer the length
```

- Arrays are value types and store a fixed number of elements of the same type.
- Zero value initialization: int → 0, string → "", bool → false, pointer → nil

### 1.2 Initialize Specific Index Elements

```go
positionInit := [5]string{1: "position1", 3: "position3"}
fmt.Println(positionInit) // ["", "position1", "", "position3", ""]
```
- Unspecified elements use zero values
- The number of initialized values cannot exceed the array length

### 1.3 Difference Between Arrays and Slices

```go
var marr [2]map[string]string
fmt.Println(marr) // [map[] map[]], but elements are nil
```
- `[2]map[string]string` is an array, fixed length
- `[]map[string]string` is a slice, variable length
- Arrays are value types, slices are reference types

---

## 2. Multidimensional Arrays

### 2.1 Two-Dimensional Arrays

```go
a := [3][2]int{
    {0,1},
    {2,3},
    {4,5},
}
```
- Access: `a[row][col]`
- Zero value initialization:
    ```go
    var arr2D [3][3]int
    arr2D[0][1] = 10
    arr2D[2][2] = 90
    ```

### 2.2 Three-Dimensional Arrays

```go
b := [3][2][2]int{
    {{0,1},{2,3}},
    {{4,5},{6,7}},
    {{8,9},{10,11}},
}
```

- Traversing multidimensional arrays requires nested for loops:

    ```go
    for i, v := range b {
        for j, inner := range v {
            fmt.Println(inner)
        }
    }
    ```

### 2.3 Multidimensional Array Assignment Notes

- Array elements can be modified inside functions
- Global arrays must be initialized with literals:

    ```go
    var arr2D = [3][3]int{
        {0, 10, 0},
        {0, 0, 90},
        {0, 0, 0},
    }
    ```
- You cannot assign values at runtime globally (e.g., `arr2D[0][1]=10`)

---

## 3. Array Parameter Passing

### 3.1 Value Passing

```go
func modifyArray(param [5]int) {
    param[0] = 100
}
a := [5]int{1,2,3,4,5}
modifyArray(a)
fmt.Println(a) // a[0] is still 1
```
- Arrays are value types; functions receive a copy, so modifications inside do not affect the original

### 3.2 Pointer Passing

```go
func modifyArrayPtr(param *[5]int) {
    param[0] = 100
}
modifyArrayPtr(&a)
fmt.Println(a) // a[0] = 100
```
- Passing an array pointer allows the function to modify the original array

---

## 4. Structs and Pointer Arrays

### 4.1 Struct Definition

```go
type Person struct {
    Name string
    Age  int
}
```

### 4.2 Array of Struct Pointers

```go
var pArr = [3]*Person{
    &Person{"Alice", 30},
    &Person{"Bob", 25},
    &Person{"Charlie", 35},
}
```
- Stores pointers; after passing as a function parameter, the pointers still point to the same memory

#### Traversal

```go
func printPersons(param [3]*Person) {
    for i, p := range param {
        fmt.Printf("index=%d, pointer=%p, Name=%s, Age=%d\n", i, p, p.Name, p.Age)
    }
}
```

---

## 5. `var` vs `:=` Declarations

| Feature   | var                | :=                   |
|-----------|--------------------|----------------------|
| Scope     | Global/Local       | Only inside function |
| Init      | Can specify or infer type | Type inferred by compiler |
| Zero value| Auto zero value    | Must assign a value  |

**Example:**
```go
var x = 10      // Valid globally
x := 10         // ❌ Invalid globally, only inside functions
```

---

## 6. Global Array Initialization Notes

- Must use literals determinable at compile time
- Cannot use runtime assignment globally
- Supports index-specified element initialization:

    ```go
    var arr = [3][3]int{
        [0][1]: 10,
        [2][2]: 90,
    }
    ```

---

## 7. Key Points Summary

- Go arrays are value types; slices are reference types.
- Multidimensional arrays can be initialized with literals or by specifying indices.
- Function parameter passing:
    - Passing arrays directly → value passing → does not modify the original array
    - Passing array pointers → reference passing → modifies the original array
- Struct arrays:
    - Storing struct pointers → can modify the pointed value in functions
    - Storing structs themselves → modifications inside functions do not affect the original array
- `:=` can only be used inside functions; global declarations must use `var