# Go Arrays and Pointers — Study Notes

本笔记总结了 Go 语言中数组与指针的基础用法、常见陷阱与注意事项，适合作为 GitHub README 参考。

---

## 1. 一维数组基础

### 1.1 声明数组

```go
var a [5]int             // 长度为5的int数组，元素零值为0
b := [5]int{1,2,3,4,5}   // 类型推导声明
c := [...]int{1,2,3}     // ...让编译器自动计算长度
```

- 数组是值类型，存储固定长度的同类型元素。
- 零值初始化：int → 0，string → ""，bool → false，指针 → nil

### 1.2 初始化指定下标元素

```go
positionInit := [5]string{1: "position1", 3: "position3"}
fmt.Println(positionInit) // ["", "position1", "", "position3", ""]
```
- 未指定的元素使用零值
- 数组长度不能被初始化值超过

### 1.3 数组与切片的区别

```go
var marr [2]map[string]string
fmt.Println(marr) // [map[] map[]]，但元素为 nil
```
- `[2]map[string]string` 是数组，长度固定
- `[]map[string]string` 是切片，长度可变
- 数组是值类型，切片是引用类型

---

## 2. 多维数组

### 2.1 二维数组

```go
a := [3][2]int{
    {0,1},
    {2,3},
    {4,5},
}
```
- 访问方式：`a[row][col]`
- 零值初始化：
    ```go
    var arr2D [3][3]int
    arr2D[0][1] = 10
    arr2D[2][2] = 90
    ```

### 2.2 三维数组

```go
b := [3][2][2]int{
    {{0,1},{2,3}},
    {{4,5},{6,7}},
    {{8,9},{10,11}},
}
```

- 遍历多维数组需嵌套 for 循环：

    ```go
    for i, v := range b {
        for j, inner := range v {
            fmt.Println(inner)
        }
    }
    ```

### 2.3 多维数组赋值注意

- 在函数内部可以修改数组元素
- 全局数组必须用字面量初始化：

    ```go
    var arr2D = [3][3]int{
        {0, 10, 0},
        {0, 0, 90},
        {0, 0, 0},
    }
    ```
- 不能在全局用运行时赋值（如 `arr2D[0][1]=10`）

---

## 3. 数组传参

### 3.1 值传递

```go
func modifyArray(param [5]int) {
    param[0] = 100
}
a := [5]int{1,2,3,4,5}
modifyArray(a)
fmt.Println(a) // a[0] 仍然是 1
```
- 数组是值类型，函数接收的是副本，内部修改不会影响外部

### 3.2 指针传递

```go
func modifyArrayPtr(param *[5]int) {
    param[0] = 100
}
modifyArrayPtr(&a)
fmt.Println(a) // a[0] = 100
```
- 传递数组指针，函数内部修改会影响原数组

---

## 4. 结构体与指针数组

### 4.1 结构体定义

```go
type Person struct {
    Name string
    Age  int
}
```

### 4.2 结构体指针数组

```go
var pArr = [3]*Person{
    &Person{"Alice", 30},
    &Person{"Bob", 25},
    &Person{"Charlie", 35},
}
```
- 存放指针，函数传参后，指针仍然指向同一块内存

#### 遍历

```go
func printPersons(param [3]*Person) {
    for i, p := range param {
        fmt.Printf("index=%d, pointer=%p, Name=%s, Age=%d\n", i, p, p.Name, p.Age)
    }
}
```

---

## 5. var vs := 声明

| 特性     | var           | :=                 |
|----------|---------------|--------------------|
| 使用范围 | 全局/局部     | 只能函数内部       |
| 初始化   | 可指定类型或推导 | 类型由编译器推导 |
| 零值     | 自动赋零值     | 必须声明并赋值     |

**例子：**
```go
var x = 10      // 全局合法
x := 10         // ❌ 全局非法，只能在函数内部
```

---

## 6. 全局数组初始化注意事项

- 必须使用编译期可确定的字面量
- 不能在全局使用运行时赋值语句
- 支持索引指定元素初始化：

    ```go
    var arr = [3][3]int{
        [0][1]: 10,
        [2][2]: 90,
    }
    ```

---

## 7. 知识点总结

- Go 数组是值类型，切片是引用类型。
- 多维数组初始化可以使用字面量或索引指定元素。
- 函数传参：
    - 直接传数组 → 值传递 → 不会修改原数组
    - 传数组指针 → 引用传递 → 会修改原数组
- 结构体数组：
    - 存放结构体指针 → 可以在函数中修改指针指向的值
    - 存放结构体本身 → 函数内部修改不会影响原数组
- `:=` 只能在函数内部使用，全局声