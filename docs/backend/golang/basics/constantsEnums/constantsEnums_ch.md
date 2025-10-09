# Go 常量与枚举 — 学习笔记

本笔记总结了 Go 中常量（`const`）的多种声明方式、分组声明、类型推断、使用 `iota` 实现枚举、自定义类型与方法绑定、指针接收者及常见易错点，适合作为 GitHub README 参考。

---

## 常量声明方式与规则

### 1. 显式类型常量

```go
const a int = 1
```
- 常量 `a` 的类型被显式声明为 `int`。
- 常量在编译期确定，不能修改。

---

### 2. 无类型常量

```go
const b = "test"
```
- Go 的常量可以没有类型，称为**无类型常量**（untyped constant）。
- 编译器会根据上下文推断类型。

示例：
```go
var s string = b  // 推断为 string
var r rune = b[0] // 推断为 rune
```

---

### 3. 多常量同一行声明

```go
const c, d = 2, "hello"
```
- 可以在一行声明多个常量。
- 类型由右侧推断。

---

### 4. 显式类型 + 多赋值

```go
const e, f bool = true, false
```
- 多个常量可以显式指定类型，但类型必须一致。

---

### 5. 分组声明

```go
const (
    h    byte = 3
    i         = "value"
    j, k      = "v", 4
    l, m      = 5, false
)
```
- 分组声明提升可读性和组织性。
- 分组内：
  - `iota` 自动递增。
  - 省略类型或表达式会复用上一行的类型和值。

---

### 6. 独立分组与 `iota`

```go
const (
    n = 6
)
```
- 每个 `const (...)` 分组是独立的。
- `iota` 在每个新分组内都会重置为 0。

---

## 命名类型与枚举

### 1. 定义新类型

```go
type Gender string
```
- 创建一个新的命名类型，而不是别名。
- 区别：
  - `type Gender = string` 是别名。
  - `type Gender string` 是新类型，可绑定方法。

---

### 2. 枚举型常量

```go
const (
    Male   Gender = "Male"
    Female Gender = "Female"
)
```
- 通过 `const` 和自定义类型实现类似枚举的结构。
- 常量值受限于 `Gender` 类型。

---

### 3. 为类型绑定方法

```go
func (g *Gender) IsMale() bool {
    return *g == Male
}

func (g *Gender) IsFemale() bool {
    return *g == Female
}
```
- 可以为 `Gender` 类型绑定方法。
- **指针接收者**：
  - 如果接收者是变量（可寻址），编译器会自动取地址。
  - 常量或字面量不能直接调用指针方法（如 `Male.IsMale()` ❌，但 `var g = Male; g.IsMale()` ✅）。
  - **注意**：如果 `g` 是 nil 指针，运行时会 panic。

---

## 使用 `iota` 实现枚举

### 1. 定义月份枚举

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
- **iota 工作原理**：
  - `iota` 是常量计数器，在每个 `const` 分组从 0 开始，每行递增 1。
  - 省略表达式时会复用上一行表达式（但 `iota` 仍递增）。

示例：
```go
January = 1 + 0 = 1
February = 1 + 1 = 2
...
December = 1 + 11 = 12
```

---

### 2. 常见 `iota` 模式

```go
const (
    _  = iota             // 跳过0
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
    TB                    // 1 << 40
)
```

---

## 指针接收者与自动取址

### 指针接收者示例

```go
func (g *Gender) IsMale() bool {
    return *g == Male
}
```

### 自动取址规则

- 如果接收者是可寻址变量，编译器会自动取地址：
  ```go
  var g = Male
  g.IsMale() // 合法
  ```
- 常量或字面量不能直接调用指针方法：
  ```go
  Male.IsMale() // 编译错误
  ```
- 如果指针为 nil，调用指针方法会导致运行时 panic。

---

## 常见易错点

| 问题                  | 说明                                                                 |
|-----------------------|----------------------------------------------------------------------|
| `iota` 重置           | 每个 `const (...)` 块都会重置 `iota` 为 0。                         |
| 指针接收者方法        | 字面量不能调用指针方法，只能变量调用。                              |
| 无类型常量            | 灵活但赋值时类型由上下文推断。                                      |
| 常量表达式            | 必须在编译期可计算（不能依赖运行时变量）。                          |
| 分组声明              | 省略值会复用上一行表达式，但 `iota` 仍递增。                        |

---

## 建议补充：String 方法与 iota 模式

### 1. 为枚举添加 `String()` 方法

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
- 示例：
  ```go
  fmt.Println(September) // 输出 "September"
  ```

---

### 2. 典型 iota 模式

```go
const (
    _  = iota             // 跳过0
    KB = 1 << (10 * iota) // 1 << 10 = 1024
    MB                    // 1 << 20
    GB                    // 1 << 30
)
```

---

### 3. 枚举类型的好处

- **可读性提升**：避免使用普通数字或字符串。
- **类型安全**：编译器检查避免混淆。
- **可扩展性**：可添加方法（如判断、打印、比较）。

---

## 知识点总结

| 分类                | 内容                                                                 |
|---------------------|----------------------------------------------------------------------|
| 常量声明            | `const name [type] = value`，支持分组。                              |
| 类型推断            | 无类型常量根据上下文推断类型。                                       |
| 枚举                | 使用 `iota` 自动递增常量。                                           |
| 命名类型            | `type NewType BaseType` 定义新类型。                                 |
| 指针接收者          | 支持自动取址，但注意 nil 指针。                                      |
| 常量规则            | 只能使用编译期可计算的值。                                           |
| 枚举字符串输出      | 实现 `String()` 方法提升可读性。                                     |

---