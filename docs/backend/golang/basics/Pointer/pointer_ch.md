# Go 指针与 unsafe 学习笔记

本笔记总结了 Go 语言中指针的基本用法、指针的指针、new、unsafe 包及相关底层操作，适合初学者和进阶读者参考。

---

## 示例代码

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

## 1. 包与导入

- `package main`：可执行程序的入口包。
- `import ("fmt" "unsafe")`：
  - `fmt`：格式化输出（推荐用于生产代码）。
  - `unsafe`：允许进行绕过类型系统和内存安全的操作（需谨慎使用）。

---

## 2. 包级变量

```go
var p *int
var q *string
```
- `p` 是 `*int` 类型（指向 int 的指针），`q` 是 `*string`。
- 包级变量会被零值初始化（指针的零值是 `nil`）。

---

## 3. 局部变量与取地址

```go
i := 10
p = &i
```
- `i` 是 int 局部变量，`&i` 取其地址赋给 `p`，此后 `*p` 可读写 `i`。

---

## 4. new 与字符串指针

```go
q = new(string)
*q = "hello"
```
- `new(string)` 分配一个 string 类型的零值并返回指针。
- `*q = "hello"` 给指针所指向的位置写入字符串。

---

## 5. 指针的指针

```go
q2 := &q
```
- `q2` 类型为 `**string`，即指向 string 指针的指针。
- `*q2 == q`，`**q2 == *q`。

---

## 6. 通过指针修改原值

```go
*p = 20
*q = "world"
```
- `*p = 20` 等价于 `i = 20`。
- `*q = "world"` 修改 q 指向的字符串。

---

## 7. 指向指针的指针并二次解引用

```go
p2 := &p
**p2 = 50
```
- `p2` 类型为 `**int`，`*p2` 是 `*int`，`**p2` 是 `int`。
- `**p2 = 50` 等价于 `i = 50`。

---

## 8. unsafe、uintptr 与内存偏移（危险操作）

```go
a := "Hello, world!"
upA := uintptr(unsafe.Pointer(&a))
upA += 1
c := (*uint8)(unsafe.Pointer(upA))
fmt.Println(*c)
```
- Go 字符串底层是一个结构体（数据指针+长度）。
- `unsafe.Pointer(&a)` 取 string header 的地址，转为整数做地址运算。
- `upA += 1` 后再转回指针并解引用，读取的是 string header 内部的第 1 个字节（不是字符串内容）。
- **风险**：此操作依赖平台、内存布局，不可移植且不安全，可能导致程序崩溃或行为异常。

---

## 总结

- 指针用于间接访问和修改变量，支持多级指针。
- `new(T)` 返回指向类型 T 的指针，`&x` 取变量地址。
- `unsafe` 包可实现底层操作，但需谨慎使用，避免破坏内存安全。
- 推荐使用 `fmt` 系列函数进行输出，避免使用 `println`。
- 指针操作是 Go 进阶开发的重要基础，理解其原理有助于写出更高效和安全的代码