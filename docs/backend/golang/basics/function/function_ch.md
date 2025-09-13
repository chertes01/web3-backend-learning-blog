# Go Functions, Methods, and Closures — Study Notes

本笔记总结了 Go 语言中的函数、方法、闭包及匿名函数立即执行（IIFE）的用法，适合作为 GitHub README 参考。

---

## 1. 函数（Function）

在 Go 中，函数是一等公民，可以像变量一样被声明、赋值、传递和返回。

### 1.1 基本声明

```go
func add(a int, b int) int {
    return a + b
}
```

### 1.2 匿名函数

匿名函数（anonymous function）即没有名字的函数，可以直接赋值给变量使用：

```go
square := func(x int) int {
    return x * x
}
fmt.Println(square(5)) // 输出 25
```

---

## 2. 方法（Method）

方法是绑定到类型上的函数。与普通函数相比，它多了一个接收者（receiver）。

### 2.1 方法声明

```go
type Counter struct {
    value int
}

// Add 方法，接收者是 *Counter
func (c *Counter) Add(n int) int {
    c.value += n
    return c.value
}
```

调用方式：

```go
c := Counter{0}
fmt.Println(c.Add(3)) // 输出 3
fmt.Println(c.Add(5)) // 输出 8
```

### 2.2 方法值（Method Value）

方法可以赋值给函数变量：

```go
fn := c.Add
fmt.Println(fn(2)) // 相当于 c.Add(2)，输出 10
```

> 注意：fn 绑定的是当前实例 c，即使之后 c 被修改为新的对象，fn 仍然指向原来的实例。

---

## 3. 闭包（Closure）

闭包 = 函数 + 它引用的外部变量环境。闭包能“记住”它定义时的上下文。

### 3.1 闭包的基本用法

```go
func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

posSum := adder()
fmt.Println(posSum(1)) // 1
fmt.Println(posSum(2)) // 3
fmt.Println(posSum(3)) // 6
```

### 3.2 闭包工厂

```go
func makeMultiplier(factor int) func(y int) int {
    return func(y int) int {
        return factor * y
    }
}

times10 := makeMultiplier(10)
fmt.Println(times10(5)) // 输出 50

// 多个独立闭包环境
times2 := makeMultiplier(2)
times3 := makeMultiplier(3)
fmt.Println(times2(5)) // 10
fmt.Println(times3(5)) // 15
```

---

## 4. 匿名函数立即执行（IIFE）

在 Go 中，可以声明匿名函数并立即执行：

```go
returnFunc := func() func(int, string) (int, string) {
    fmt.Println("this is an anonymous function")
    return func(i int, s string) (int, string) {
        return i, s
    }
}() // 立即调用

ret1, ret2 := returnFunc(1, "test")
fmt.Println(ret1, ret2) // 输出: 1 test
```

---

## 5. 综合示例

将函数、方法、闭包放到一起：

```go
package main

import "fmt"

type Counter struct {
    value int
}

func (c *Counter) Add(n int) int {
    c.value += n
    return c.value
}

func makeMultiplier(factor int) func(y int) int {
    return func(y int) int {
        return factor * y
    }
}

func main() {
    // 方法
    counter := Counter{0}
    fn := counter.Add
    fmt.Println(fn(3)) // 3
    fmt.Println(fn(5)) // 8

    // 函数
    square := func(x int) int { return x * x }
    fmt.Println(square(5)) // 25

    // 闭包
    times10 := makeMultiplier(10)
    fmt.Println(times10(5)) // 50

    // 匿名函数立即执行
    returnFunc := func() func(int, string) (int, string) {
        fmt.Println("this is an anonymous function")
        return func(i int, s string) (int, string) {
            return i, s
        }
    }()
    ret1, ret2 := returnFunc(1, "GoLang")
    fmt.Println(ret1, ret2) // 1 GoLang
}
```

---

## 📌 总结

- **函数**：最基本的可复用代码块，可以作为值传递。
- **方法**：绑定到类型的函数，带有接收者。
- **闭包**：函数与其外部变量的组合，能保持状态。
- **匿名函数立即执行**：声明函数后立刻执行，常用于初始化逻