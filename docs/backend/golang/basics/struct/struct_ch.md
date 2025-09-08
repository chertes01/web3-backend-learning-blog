# Go 学习笔记：结构体、方法、值与指针

---

## 1. 结构体字段标签（json / gorm）

```go
type User struct {
    Name string `json:"name" gorm:"column:user_name"`
}
```

- `json:"name"`：控制 Go 结构体与 JSON 的映射关系。
    - JSON 输入：`{"name": "Tom"}` → 自动填充到 `User.Name`
    - JSON 输出：结构体转 JSON 时字段叫 "name"
- `gorm:"column:user_name"`：控制 ORM 框架 GORM 的行为。
    - 指定数据库表中的列名是 `user_name`
    - 插入/查询时会用 `user_name` 作为映射字段

✅ 典型流程：前端 JSON → Go 结构体 → GORM → 数据库

---

## 2. 匿名结构体 vs 具名结构体

**具名结构体**：有名字，可复用

```go
type Person struct {
    Name string
    Age  int
}
```

**匿名结构体**：定义时直接使用，不具备名字，常用于一次性场景

```go
var tmp = struct {
    Name string
    Age  int
}{"Tom", 18}
```

> ⚠️ 匿名结构体指的是「结构体本身没有名字」，不是字段类型匿名。

---

## 3. 变量作用域与遮蔽（Shadowing）

```go
var globalVar = struct{}{} // 全局变量

func main() {
    var globalVar = struct {
        Msg string
    }{Msg: "local"}
    fmt.Println(globalVar) // 打印局部变量
}
```

- Go 允许全局和局部变量同名
- 局部变量会遮蔽全局变量（shadowing）
- 在局部作用域中优先使用局部变量

---

## 4. 空结构体通道（chan struct{}）

```go
done := make(chan struct{})

go func() {
    fmt.Println("goroutine running")
    done <- struct{}{} // 发送信号
}()

<-done // 等待信号
```

- `struct{}` 空结构体占用 0 字节，常用于信号通知
- `make(chan struct{}, 0)`：无缓冲通道，发送和接收需要同步

---

## 5. 结构体嵌套与字段遮蔽

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

- 字段查找优先级：结构体自身 → 内嵌结构体 → 更深层嵌套
- 同名字段会被外层字段遮蔽
- ⚠️ `first := One{a: "Apple"}` 并不会影响 `third.One.a`，因为它们是不同实例。

---

## 6. 方法接收者：值 vs 指针

```go
type User struct {
    name string
}

// 值接收者
func (u User) SetNameCopy(newName string) {
    u.name = newName // 修改副本，不影响原对象
}

// 指针接收者
func (u *User) SetNamePtr(newName string) {
    u.name = newName // 修改原对象
}
```

- 值接收者 `(u User)`：方法参数是副本，修改不影响原对象
- 指针接收者 `(u *User)`：方法参数是指针，修改会影响原对象
- 调用时 Go 会自动转换：
    - 值也能调用指针方法（自动取地址）
    - 指针也能调用值方法（自动解引用）

---

## 7. 函数参数：值传递 vs 指针传递

```go
func UpdateUserByValue(u User, newName string) {
    u.name = newName // 修改副本
}

func UpdateUserByPointer(u *User, newName string) {
    u.name = newName // 修改原对象
}
```

- 传值：拷贝一份，不会影响外部变量
- 传指针：直接修改原对象

---

## 8. 示例运行流程

```go
user := User{name: "Initial"}

UpdateUserByValue(user, "ValueFunc") // 不变
UpdateUserByPointer(&user, "PointerFunc") // 改变

user.SetNameCopy("SetNameCopy") // 不变
user.SetNamePtr("SetNamePtr")   // 改变
```

**输出：**

```
初始值: Initial
UpdateUserByValue 后: Initial
UpdateUserByPointer 后: PointerFunc
SetNameCopy 方法后: PointerFunc
SetNamePtr 方法后: SetNamePtr
```

---

## ✅ 总结

- 结构体嵌套：外层字段优先，可能遮蔽内层字段
- 作用域：局部变量可以遮蔽全局变量
- 函数参数：传值不改原对象，传指针能改原对象
- 方法接收者：值接收者操作副本，指针接收者操作原对象
- Go 自动转换：值/指针都能调用两类方法