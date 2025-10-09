# 🧠 Go Slice Study Notes

---

## 1. Basic Concepts of Slices

- A slice is a **dynamic view** based on an array, supporting expansion, slicing, and element modification.
- Unlike arrays, slices do not require a fixed length.

```go
s := []int{1, 2, 3}   // Define a slice
a := [3]int{1, 2, 3}  // Define an array
```

- Array: fixed length `[3]int`
- Slice: dynamic length `[]int`

A slice consists of three parts under the hood:
- Pointer: points to a position in the underlying array
- Length (`len`): number of elements in the slice
- Capacity (`cap`): size from the slice's start to the end of the underlying array

---

## 2. Ways to Create Slices

1. **Literal creation**
    ```go
    slice1 := []int{1, 2, 3, 4}
    ```
2. **Using make()**
    ```go
    slice2 := make([]int, 4, 5)
    // len = 4, cap = 5
    ```
3. **Slicing from an array**
    ```go
    array := [...]string{"zero", "one", "two", "three", "four", "five", "six"}
    slice3 := array[3:7]
    // slice3 = [three four five six]
    // len(slice3) = 4, cap(slice3) = 4
    ```

---

## 3. Length and Capacity

- `len`: number of accessible elements in the slice
- `cap`: available space from the start position to the end of the underlying array

```go
fmt.Printf("slice1:%v ,length:%d ,capacity:%d\n", slice1, len(slice1), cap(slice1))
```

Slices share the same underlying array. `len` is just the view size, and `cap` determines if more elements can be appended.

---

## 4. Accessing and Modifying Slices

```go
s1 := []int{5, 4, 3, 2, 1}
fmt.Println("s1:", s1[:3]) // [5 4 3]
s1 = append([]int{10, 9, 8}, s1[3:]...)
fmt.Println("s1 after modification:", s1) // [10 9 8 2 1]
```

- `s1[:3]` gets the first 3 elements
- `append([]int{10, 9, 8}, s1[3:]...)` inserts new elements at the front

---

## 5. Iterating Over Slices

```go
for v := range s {
    fmt.Println(v, s[v])
}
```
- `range` returns the index, `s[v]` returns the value

---

## 6. nil Slice vs Empty Slice

```go
var nilSlice []int // nil slice
s3 := []int{}      // empty slice
```

| Type      | len | cap | Is nil | Memory Allocation   |
|-----------|-----|-----|--------|--------------------|
| nil slice |  0  |  0  | ✅ Yes | ❌ No array        |
| empty     |  0  |  0  | ❌ No  | ✅ Has array       |

```go
fmt.Println(nilSlice == nil) // true
fmt.Println(s3 == nil)       // false
```

---

## 7. Using append and Expansion Mechanism

```go
s := []string{}
s = append(s, "hello")
s = append(s, "world", "wellcome")
s = append(s[:3], []string{"to", "go", "language"}...)
```

- append automatically expands the slice if capacity is insufficient
- Expansion rules:
    - Capacity < 1024: roughly doubles
    - Capacity ≥ 1024: increases by 25% each time
- After expansion, a new underlying array is allocated and the original content is copied

---

## 8. Insert and Delete Operations

**Insert element:**
```go
s5 := array[0:7]
s5 = append(s5[0:2], append([]string{"3"}, s5[2:7]...)...)
fmt.Println(s5) // Insert "3" at index 2
```

**Delete element:**
```go
s6 := array[0:7]
s6 = append(s6[0:3], s6[4:7]...)
fmt.Println(s6) // Delete element at index 3
```

---

## 9. copy Function

```go
src := []int{1, 2, 3, 4}
dst := make([]int, 2)
copy(dst, src)
```
- Only copies the first `len(dst)` elements
- Modifying dst does not affect src
- If both share the underlying array (e.g., `dst = src[1:3]`), modifications affect both

---

## 10. Expansion and Sharing Notes

- append sometimes returns a new underlying array, sometimes not
    - If cap is sufficient: reuses the original array
    - If cap is insufficient: allocates a new array (detaches from the original slice)

```go
s1 := []int{1, 2, 3}
s2 := append(s1, 4)
```
- If cap(s1) is sufficient, s1 and s2 share the underlying array; otherwise, s2 points to a new array

---

## 11. Experiment: Expansion Pattern Example

```go
s := []int{}
for i := 0; i < 20; i++ {
    s = append(s, i)
    fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
}
```

Sample output:

```
len=1 cap=1
len=2 cap=2
len=3 cap=4
len=5 cap=8
len=9 cap=16
len=17 cap=32
```
- Less than 1024: doubles in size

---

## 🧾 Summary: Slice Key Points

| Operation      | Syntax                                      | Description                |
|----------------|---------------------------------------------|----------------------------|
| Define slice   | s := []int{}                                | Empty slice                |
| Use make       | make([]T, len, cap)                         | Specify length and capacity|
| From array     | a[low:high]                                 | Slice range [low, high)    |
| Insert         | append(s[:i], append([]T{x}, s[i:]...)...)  | Insert element             |
| Delete         | append(s[:i], s[i+1:]...)                   | Delete element             |
| Copy           | copy(dst, src)                              | Copy minimum length elements|
| nil vs empty   | nilSlice == nil                             | Empty slice is not nil     |
| Expansion      | auto double or +25%                         | New array allocated if needed|

---