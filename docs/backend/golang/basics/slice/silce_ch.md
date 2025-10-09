# 🧠 Go 切片（Slice）学习笔记

---

## 一、切片的基本概念

- 切片（Slice）是基于数组的**动态视图**，可扩容、截取和修改元素。
- 与数组不同，切片不需要指定固定长度。

```go
s := []int{1, 2, 3}   // 定义切片
a := [3]int{1, 2, 3}  // 定义数组
```

- 数组：固定长度 `[3]int`
- 切片：动态长度 `[]int`

切片底层由三部分组成：
- 指针：指向底层数组中的某个位置
- 长度（len）：切片中元素个数
- 容量（cap）：从切片起始位置到底层数组末尾的大小

---

## 二、切片的创建方式

1. **字面量创建**
    ```go
    slice1 := []int{1, 2, 3, 4}
    ```
2. **使用 make()**
    ```go
    slice2 := make([]int, 4, 5)
    // len = 4，cap = 5
    ```
3. **从数组切割**
    ```go
    array := [...]string{"zero", "one", "two", "three", "four", "five", "six"}
    slice3 := array[3:7]
    // slice3 = [three four five six]
    // len(slice3) = 4, cap(slice3) = 4
    ```

---

## 三、长度与容量

- `len`：当前切片可访问元素数量
- `cap`：从起始位置到底层数组末尾的可用空间

```go
fmt.Printf("slice1:%v ,length:%d ,capacity:%d\n", slice1, len(slice1), cap(slice1))
```

切片在底层共享同一个数组，`len` 只是视图大小，`cap` 决定是否能继续追加。

---

## 四、访问与修改切片

```go
s1 := []int{5, 4, 3, 2, 1}
fmt.Println("s1:", s1[:3]) // [5 4 3]
s1 = append([]int{10, 9, 8}, s1[3:]...)
fmt.Println("s1 after modification:", s1) // [10 9 8 2 1]
```

- `s1[:3]` 取前3个元素
- `append([]int{10, 9, 8}, s1[3:]...)` 在前面插入新元素

---

## 五、切片遍历

```go
for v := range s {
    fmt.Println(v, s[v])
}
```
- `range` 返回索引，`s[v]` 返回值

---

## 六、nil 切片 vs 空切片

```go
var nilSlice []int // nil 切片
s3 := []int{}      // 空切片
```

| 类型     | len | cap | 是否为 nil | 内存分配      |
|----------|-----|-----|------------|---------------|
| nil 切片 |  0  |  0  | ✅ 是      | ❌ 无底层数组 |
| 空切片   |  0  |  0  | ❌ 否      | ✅ 有底层数组 |

```go
fmt.Println(nilSlice == nil) // true
fmt.Println(s3 == nil)       // false
```

---

## 七、append 的使用与扩容机制

```go
s := []string{}
s = append(s, "hello")
s = append(s, "world", "wellcome")
s = append(s[:3], []string{"to", "go", "language"}...)
```

- append 会在容量不足时自动扩容
- 扩容规则：
    - 容量 < 1024 时，大致翻倍
    - 容量 ≥ 1024 时，每次增加 25%
- 扩容后会分配新底层数组并 copy 原内容

---

## 八、插入与删除操作

**插入元素：**
```go
s5 := array[0:7]
s5 = append(s5[0:2], append([]string{"3"}, s5[2:7]...)...)
fmt.Println(s5) // 在索引2处插入"3"
```

**删除元素：**
```go
s6 := array[0:7]
s6 = append(s6[0:3], s6[4:7]...)
fmt.Println(s6) // 删除索引3的元素
```

---

## 九、copy 函数

```go
src := []int{1, 2, 3, 4}
dst := make([]int, 2)
copy(dst, src)
```
- 只复制前 `len(dst)` 个元素
- 修改 dst 不影响 src
- 若两者共享底层数组（如 `dst = src[1:3]`），修改会互相影响

---

## 十、扩容与共享注意点

- append 有时会返回新的底层数组，有时不会
    - 如果 cap 足够：复用原数组
    - 如果 cap 不够：分配新数组（与原切片脱离关系）

```go
s1 := []int{1, 2, 3}
s2 := append(s1, 4)
```
- 若 cap(s1) 足够，s1 和 s2 共享底层数组；否则 s2 指向新数组

---

## 十一、实验：扩容规律示例

```go
s := []int{}
for i := 0; i < 20; i++ {
    s = append(s, i)
    fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
}
```

输出类似：

```
len=1 cap=1
len=2 cap=2
len=3 cap=4
len=5 cap=8
len=9 cap=16
len=17 cap=32
```
- 小于 1024 时为翻倍增长

---

## 🧾 总结：切片学习重点

| 操作   | 语法                                      | 说明                 |
|--------|-------------------------------------------|----------------------|
| 定义切片 | s := []int{}                             | 空切片               |
| 使用 make | make([]T, len, cap)                    | 指定长度与容量       |
| 从数组创建 | a[low:high]                            | 截取范围 [low, high) |
| 插入     | append(s[:i], append([]T{x}, s[i:]...)...) | 插入元素             |
| 删除     | append(s[:i], s[i+1:]...)               | 删除元素             |
| 拷贝     | copy(dst, src)                          | 拷贝最小长度元素     |
| nil vs 空切片 | nilSlice == nil                     | 空切片不为 nil       |
| 扩容     | 自动翻倍或+25%                           | 超出容量自动分配新数组 |

---