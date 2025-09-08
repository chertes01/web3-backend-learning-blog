...existing code...
# Go Study Notes: Structs, Methods, Value vs Pointer

---

## 1. Struct Field Tags (json / gorm)

```go
type User struct {
    Name string `json:"name" gorm:"column:user_name"`
}
```

- `json:"name"`: controls mapping between Go struct and JSON.
    - JSON input: `{"name": "Tom"}` → fills `User.Name`
    - JSON output: the field is named "name" when struct is marshaled
- `gorm:"column:user_name"`: controls ORM (GORM) behavior.
    - Specifies the column name in the database as `user_name`
    - Insert/query operations use `user_name` for mapping

✅ Typical flow: frontend JSON → Go struct → GORM → database

---

## 2. Anonymous Struct vs Named Struct

**Named struct**: has a name and can be reused

```go
type Person struct {
    Name string
    Age  int
}
```

**Anonymous struct**: defined inline, has no name, often used for one-off scenarios

```go
var tmp = struct {
    Name string
    Age  int
}{"Tom", 18}
```

> ⚠️ "Anonymous struct" refers to the struct type having no name, not an anonymous field type.

---

## 3. Variable Scope and Shadowing

```go
var globalVar = struct{}{} // global variable

func main() {
    var globalVar = struct {
        Msg string
    }{Msg: "local"}
    fmt.Println(globalVar) // prints the local variable
}
```

- Go allows global and local variables with the same name
- Local variables shadow global variables
- The local variable is used inside its scope

---

## 4. Empty-struct Channel (chan struct{})

```go
done := make(chan struct{})

go func() {
    fmt.Println("goroutine running")
    done <- struct{}{} // send signal
}()

<-done // wait for signal
```

- `struct{}` (empty struct) occupies 0 bytes and is often used for signaling
- `make(chan struct{}, 0)`: unbuffered channel; send and receive synchronize

---

## 5. Struct Embedding and Field Shadowing

```go
type One struct{ a string }
type Two struct {
    One
    b string
}
type Three struct {
    One
    Two
    a string
    b string
    c string
}

third := Three{a: "A", b: "B", c: "C"}
fmt.Println(third.a)        // "A"
fmt.Println(third.One.a)    // ""
fmt.Println(third.Two.a)    // ""
fmt.Println(third.Two.b)    // ""
```

- Field lookup priority: the struct itself → embedded structs → deeper embeddings
- Fields with the same name are shadowed by outer fields
- ⚠️ `first := One{a: "Apple"}` does not affect `third.One.a` — they are different instances

---

## 6. Method Receivers: Value vs Pointer

```go
type User struct {
    name string
}

// Value receiver
func (u User) SetNameCopy(newName string) {
    u.name = newName // modifies a copy, does not affect the original
}

// Pointer receiver
func (u *User) SetNamePtr(newName string) {
    u.name = newName // modifies the original
}
```

- Value receiver `(u User)`: the method operates on a copy; modifications do not affect the original
- Pointer receiver `(u *User)`: the method operates on the original via pointer; modifications persist
- Go automatically converts when calling methods:
    - A value can call a pointer method (address is taken automatically)
    - A pointer can call a value method (automatic dereference)

---

## 7. Function Parameters: Pass by Value vs Pass by Pointer

```go
func UpdateUserByValue(u User, newName string) {
    u.name = newName // modifies a copy
}

func UpdateUserByPointer(u *User, newName string) {
    u.name = newName // modifies the original
}
```

- Passing by value: a copy is made; the original is not changed
- Passing by pointer: the original object is modified

---

## 8. Example Execution Flow

```go
user := User{name: "Initial"}

UpdateUserByValue(user, "ValueFunc") // no change
UpdateUserByPointer(&user, "PointerFunc") // changes

user.SetNameCopy("SetNameCopy") // no change
user.SetNamePtr("SetNamePtr")   // changes
```

**Output:**

```
Initial value: Initial
After UpdateUserByValue: Initial
After UpdateUserByPointer: PointerFunc
After SetNameCopy: PointerFunc
After SetNamePtr: SetNamePtr
```

---

## ✅ Summary

- Struct embedding: outer fields take precedence and can shadow inner fields
- Scope: local variables can shadow globals
- Function parameters: pass-by-value does not change the original, pass-by-pointer can
- Method receivers: value receivers operate on copies, pointer receivers operate on the original
- Go's automatic conversions: both values and pointers can call either kind of method