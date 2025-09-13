# Go 条件语句 学习笔记

简要说明 Go 中常用的条件控制结构：if、switch 的用法、注意点与示例。适合作为 GitHub README 使用的参考文档。

---

## 目录

- [if 语句](#if-语句)
  - 基本用法
  - 初始化语句与作用域
  - 示例：checkNumber
- [switch 语句](#switch-语句)
  - 值匹配（基本用法）
  - 带初始化语句
  - 无表达式 switch（条件判断）
  - 类型判断（type switch）
  - 综合示例：analyzeInput
- [总结与建议](#总结与建议)

---

## if 语句

### 基本用法
```go
if condition {
    // do something
} else if anotherCondition {
    // do another
} else {
    // default case
}
```

### 初始化语句与作用域
Go 支持在 `if` 条件前写初始化语句，声明的变量只在该 `if`-`else` 链内有效：
```go
if x := 10; x > 5 {
    fmt.Println("x > 5")
} else {
    fmt.Println("x <= 5")
}
// 这里 x 已超出作用域
```

示例（变量作用域）：
```go
func main() {
    var a int = 10
    if b := 1; a > 10 {
        b = 2
        fmt.Println("a > 10")
    } else if c := 3; b > 1 {  // b 来自 if 初始化，c 只在 else if 内有效
        b = 3
        fmt.Println("b > 1")
    } else {
        fmt.Println("其他")
        fmt.Println(b) // b 仍然有效
        // fmt.Println(c) // 编译错误：c 超出作用域
    }
}
```

提示：如需在多个分支访问同一变量，可在外层提前声明（示例里用 `flag`, `m`）。

### 示例：checkNumber
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

## switch 语句

### 值匹配（基本用法）
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

### 带初始化语句
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

### 无表达式 switch（条件判断）
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

### 类型判断（type switch）
用于判断 `interface{}` 的具体类型：
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

### 综合示例：analyzeInput
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

调用示例（略）会输出针对不同类型和内容的判断结果。

---

## 总结

- if：适合复杂的布尔逻辑和需要初始化临时变量的场景，但注意初始化变量的作用域限制。
- switch：更简洁，适用于多分支判断、类型分支或无需复杂条件的场景。
- 一般建议：
  - 简单多分支 → 用 switch。
  - 复杂逻辑判断 → 用 if。

---