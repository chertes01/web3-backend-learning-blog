# Go Loop and Control Statements Notes

This note summarizes common techniques and tips for loops, control statements, map/slice sorting, random numbers, etc. in Go, suitable for beginners to review and reference.

---

## 1. Basic Usage of for Loops

Go only has the `for` loop, but it's very flexible:

- **Classic counting loop**
    ```go
    for i := 0; i < 10; i++ {
        fmt.Println("Loop number", i)
    }
    ```

- **Conditional loop**
    ```go
    i := 0
    for i < 10 {
        fmt.Println(i)
        i++
    }
    ```

- **Infinite loop**
    ```go
    for {
        fmt.Println("Infinite loop")
        break
    }
    ```

---

## 2. Iterating Arrays, Slices, and Maps

- **Array/Slice**
    ```go
    a := [3]string{"A", "B", "C"}
    for i, v := range a {
        fmt.Printf("a[%d] = %s\n", i, v)
    }

    s := []int{1, 2, 3}
    for i, v := range s {
        fmt.Printf("s[%d] = %d\n", i, v)
    }
    ```

- **Map (unordered)**
    ```go
    m := map[string]int{"a":1, "b":2, "c":3}
    for k, v := range m {
        fmt.Println(k, v)
    }
    ```

- **Ordered map iteration**
    ```go
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Println(k, m[k])
    }
    ```

---

## 3. The sort Package and Sorting Techniques

- **String sorting (lexicographical)**
    ```go
    sort.Strings(keys)
    // Note: "Student1", "Student10", "Student2" will be sorted as Student1, Student10, Student2
    ```

- **Custom sorting (by number)**
    ```go
    sort.Slice(keys, func(i, j int) bool {
        ni, _ := strconv.Atoi(strings.TrimPrefix(keys[i], "Student"))
        nj, _ := strconv.Atoi(strings.TrimPrefix(keys[j], "Student"))
        return ni < nj
    })
    ```

- **Common sorts**
    ```go
    sort.Ints([]int)
    sort.Float64s([]float64)
    sort.Strings([]string)
    ```

---

## 4. break / continue / Labels

- **break exits the current loop or switch**
    ```go
    for i := 0; i < 5; i++ {
        if i == 3 {
            break
        }
        fmt.Println(i)
    }
    ```

- **continue skips the current iteration**
    ```go
    for i := 0; i < 5; i++ {
        if i == 3 {
            continue
        }
        fmt.Println(i)
    }
    ```

- **Labelled break/continue to exit multiple loops**
    ```go
    outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                break outer // exit the outer loop
            }
        }
    }
    ```

---

## 5. Flexible Usage of switch

- **Switch without an expression**
    ```go
    switch {
    case score < 60:
        fmt.Println("Failed")
    case score < 80:
        fmt.Println("Passed")
    default:
        fmt.Println("Excellent")
    }
    ```

- **Switch with an expression**
    ```go
    switch score {
    case 100:
        fmt.Println("Perfect")
    case 90, 95:
        fmt.Println("High Score")
    }
    ```

- **fallthrough to force execution of the next case**
    ```go
    switch x := 2; x {
    case 2:
        fmt.Println("Two")
        fallthrough
    case 3:
        fmt.Println("Three") // will also execute
    }
    ```

---

## 6. Random Numbers with rand

- **Generate random integer**
    ```go
    rand.Seed(time.Now().UnixNano()) // Set before each run
    n := rand.Intn(10) // [0,10)
    ```

- **Generate random number in [min, max]**
    ```go
    n := rand.Intn(max-min+1) + min
    ```

---

## 7. Initializing Maps and Slices

- **Map initialization**
    ```go
    m1 := make(map[string]int)
    m2 := map[string]int{"A":1, "B":2}
    var m3 map[string]int // nil map, cannot assign directly
    ```

- **Slice initialization**
    ```go
    s := make([]string, 0, len(m))
    s = append(s, "A")
    ```

---

## 8. Naming Conventions and Export Rules

- **Function/variable names starting with uppercase: exported (public)**
- **Lowercase: only visible within the package**

    ```go
    func PublicFunc() {}   // public
    func privateFunc() {}  // private
    ```

---

## 9. Common Mistakes and Tips

- Map iteration is unordered; sort keys for ordered output
- rand.Intn only takes one argument; calculate range manually
- break only exits one loop; use labels or return to exit multiple
- Lowercase function names are not exported; use uppercase for cross-package calls
- Always set rand.Seed, otherwise results are the same every run

---

## 10. Summary Table

| Statement      | Effect                              | Scope         |
| -------------- | ----------------------------------- | ------------ |
| break          | Exit current loop/switch            | Current level |
| break label    | Exit to specified label             | Multi-level   |
| continue       | Skip to next iteration              | Current level |

---

> **Tip:** Go's for-range is better for iterating maps/slices, while classic for is better for counting and array index operations. Map iteration order is random; sort keys for ordered