# Go Map Study Notes

本笔记总结了 Go 语言中 map 的基础用法、并发安全、嵌套结构、常见错误、高级技巧及底层原理，适合作为 GitHub README 参考。

---

## 1. Map 基础概念

### 1.1 定义与初始化

```go
var m map[string]int // 声明但未初始化，m 为 nil，不能直接写入
m := make(map[string]int)         // 使用 make 初始化
m2 := make(map[string]int, 10)    // 指定初始容量
m3 := map[string]int{}            // 字面量初始化
m4 := map[string]int{"x": 1, "y": 2}
```
- 未初始化的 map 不能写入元素，否则 panic
- make 初始化空 map，容量仅优化性能

---

## 2. 增、删、改、查

### 2.1 插入 / 修改

```go
m["a"] = 100 // 新增或覆盖 key
```

### 2.2 查询

```go
v := m["a"]          // key 存在 → 返回 value
v = m["not_exist"]   // key 不存在 → 返回零值
v, ok := m["a"]      // 推荐用 ok 判断 key 是否存在
```

### 2.3 删除

```go
delete(m, "a") // 删除 key，不存在也不会报错
```

### 2.4 遍历

```go
for k, v := range m {
    fmt.Println(k, v)
}
```

嵌套 map 遍历：

```go
students := map[string]map[string]float64{
    "Alice": {"math": 95.5, "english": 88.0},
    "Bob":   {"math": 76.0, "english": 90.0},
}
for name, subjects := range students {
    fmt.Println("Student:", name)
    for sub, score := range subjects {
        fmt.Printf("  %s: %.2f\n", sub, score)
    }
}
```

---

## 3. 嵌套 map 的使用

- 写入前需确保内层 map 已初始化，否则 panic
- 读取不存在 key 返回零值，不报错
- 判断 key 是否存在需用 ok

```go
if student[name] == nil {
    student[name] = make(map[string]float64)
}
student[name][sub] = score
```

---

## 4. 零值与 nil

- 值类型（int, float64, bool, struct）：未写入 key 返回零值
- 引用类型（slice, map, pointer, channel, interface）：未写入 key 返回 nil
- 判断 key 是否存在用 ok，判断 value 是否为 nil 仅对引用类型有效

---

## 5. 并发安全

- Go 的 map 不是并发安全的，多 goroutine 读写需加锁，否则 panic

### 5.1 加锁方式

```go
var mu sync.Mutex
mu.Lock()
m[key] = value
v := m[key]
mu.Unlock()
```

读多写少场景推荐 RWMutex：

```go
var rw sync.RWMutex
rw.RLock()
v := m[key]
rw.RUnlock()
rw.Lock()
m[key] = value
rw.Unlock()
```

---

## 6. 常见错误

- nil map 写入会 panic
- 并发读写 map 会 panic
- 未初始化嵌套 map 写入会 panic

---

## 7. 高级应用：封装管理结构

线程安全的 map 管理器示例：

```go
type StudentManager struct {
    mu      sync.Mutex
    student map[string]map[string]float64
}

func NewManager() *StudentManager {
    return &StudentManager{
        student: make(map[string]map[string]float64),
    }
}

func (m *StudentManager) AddStudent(name string) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if _, ok := m.student[name]; !ok {
        m.student[name] = make(map[string]float64)
    }
}

func (m *StudentManager) AddScore(name, sub string, score float64) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if m.student[name] == nil {
        m.student[name] = make(map[string]float64)
    }
    m.student[name][sub] = score
}

func (m *StudentManager) GetScore(name, sub string) (float64, bool) {
    m.mu.Lock()
    defer m.mu.Unlock()
    if subjects, ok := m.student[name]; ok {
        score, ok2 := subjects[sub]
        return score, ok2
    }
    return 0, false
}

func (m *StudentManager) PrintAll() {
    m.mu.Lock()
    defer m.mu.Unlock()
    for name, subjects := range m.student {
        fmt.Println("Name:", name)
        for sub, score := range subjects {
            fmt.Printf("  %s: %.2f\n", sub, score)
        }
    }
}
```

---

## 8. 实践总结与经验

- map 声明后必须初始化（make 或字面量）
- 嵌套 map 写入前需初始化内层 map
- 判断 key 是否存在用 ok
- 并发场景必须加锁或用 sync.Map

---

## 9. 进阶技巧与底层原理

### 9.1 key 类型限制

- 只能用可比较类型（int, string, bool, 指针、数组等）作 key
- 不可用 slice、map、函数作 key
- 浮点数作 key 有精度风险，建议转为字符串或整数

### 9.2 map 的迭代顺序

- range 遍历 map 顺序随机
- 固定顺序需排序 key 列表

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

### 9.3 容量与性能优化

- make(map[string]int, n) 可指定初始容量，减少扩容开销

### 9.4 sync.Map 用法

Go 1.9+ 提供线程安全 map：

```go
var sm sync.Map
sm.Store("Alice", 100)
v, ok := sm.Load("Alice")
sm.Delete("Alice")
sm.Range(func(key, value any) bool {
    fmt.Println(key, value)
    return true
})
```

### 9.5 map 的内存管理与 GC

- 删除 key 后，value 引用会被 GC 回收，但桶结构不会立即释放
- 大 map 删除大量元素后可新建 map 复制有效数据，减少内存占用

---

## 10. 总结

- map 是高效的键值对数据结构，需注意初始化和并发安全
- 嵌套 map 写入前需初始化内层
- 并发读写需加锁或用 sync.Map
- 删除 key 只释放 value 引用，桶结构仍占内存
- 大 map 批量删除后建议新建 map 优化内存
- key 类型必须可比较，遍历顺序随机