# Go Loop and Control Statements Notes

本笔记总结了 Go 语言中循环、控制语句、map/slice 排序、随机数等常用技巧与注意事项，适合初学者查阅与复习。

---

## 1. for 循环的基本用法

Go 只有 for 一种循环结构，但用法灵活：

- **经典计数循环**
    ```go
    for i := 0; i < 10; i++ {
        fmt.Println("第", i, "次循环")
    }
    ```

- **条件循环**
    ```go
    i := 0
    for i < 10 {
        fmt.Println(i)
        i++
    }
    ```

- **无限循环**
    ```go
    for {
        fmt.Println("死循环")
        break
    }
    ```

---

## 2. 遍历数组、切片、map

- **数组/切片**
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

- **map（无序）**
    ```go
    m := map[string]int{"a":1, "b":2, "c":3}
    for k, v := range m {
        fmt.Println(k, v)
    }
    ```

- **有序遍历 map**
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

## 3. sort 包与排序技巧

- **字符串排序（字典序）**
    ```go
    sort.Strings(keys)
    // 注意："Student1", "Student10", "Student2" 排序后是 Student1, Student10, Student2
    ```

- **自定义排序（按数字）**
    ```go
    sort.Slice(keys, func(i, j int) bool {
        ni, _ := strconv.Atoi(strings.TrimPrefix(keys[i], "Student"))
        nj, _ := strconv.Atoi(strings.TrimPrefix(keys[j], "Student"))
        return ni < nj
    })
    ```

- **常用排序**
    ```go
    sort.Ints([]int)
    sort.Float64s([]float64)
    sort.Strings([]string)
    ```

---

## 4. break / continue / 标签

- **break 跳出当前循环或 switch**
    ```go
    for i := 0; i < 5; i++ {
        if i == 3 {
            break
        }
        fmt.Println(i)
    }
    ```

- **continue 跳过本次迭代**
    ```go
    for i := 0; i < 5; i++ {
        if i == 3 {
            continue
        }
        fmt.Println(i)
    }
    ```

- **标签 break/continue 跳出多层循环**
    ```go
    outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 1 {
                break outer // 跳出外层循环
            }
        }
    }
    ```

---

## 5. switch 的灵活用法

- **无表达式 switch**
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

- **带表达式 switch**
    ```go
    switch score {
    case 100:
        fmt.Println("Perfect")
    case 90, 95:
        fmt.Println("High Score")
    }
    ```

- **fallthrough 强制执行下一个 case**
    ```go
    switch x := 2; x {
    case 2:
        fmt.Println("Two")
        fallthrough
    case 3:
        fmt.Println("Three") // 也会执行
    }
    ```

---

## 6. 随机数 rand

- **生成随机整数**
    ```go
    rand.Seed(time.Now().UnixNano()) // 建议每次运行前设置
    n := rand.Intn(10) // [0,10)
    ```

- **生成指定区间 [min, max] 的随机数**
    ```go
    n := rand.Intn(max-min+1) + min
    ```

---

## 7. map 和 slice 初始化

- **map 初始化**
    ```go
    m1 := make(map[string]int)
    m2 := map[string]int{"A":1, "B":2}
    var m3 map[string]int // nil map，不能直接赋值
    ```

- **slice 初始化**
    ```go
    s := make([]string, 0, len(m))
    s = append(s, "A")
    ```

---

## 8. 命名规范与导出规则

- **函数/变量首字母大写：对外可见（导出）**
- **小写：仅包内可见**

    ```go
    func PublicFunc() {}   // 公有
    func privateFunc() {}  // 私有
    ```

---

## 9. 常见错误与注意事项

- map 遍历无序，如需有序输出需排序 key
- rand.Intn 只有单参数，区间需手动计算
- break 只退出一层循环，需标签或 return 跳出多层
- 小写函数名不能被导出，跨包调用需大写
- rand.Seed 必须设置，否则每次结果相同

---

## 10. 总结表

| 语句           | 作用                         | 范围         |
| -------------- | ---------------------------- | ------------ |
| break          | 跳出当前循环/switch           | 当前层       |
| break label    | 跳出到指定标签                | 可跨多层循环 |
| continue       | 跳过当前迭代进入下一次循环    | 当前层       |

---

> **Tip:** Go 的 for-range 更适合遍历 map/slice，普通 for 更适合计数和数组下标操作。map 的遍历顺序每次都可能不同，排序后输出才有序。