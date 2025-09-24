# Go Map Study Notes

This note summarizes the basics, concurrency safety, nested structures, common mistakes, advanced techniques, and underlying principles of maps in Go. Suitable as a GitHub README reference.

---

## 1. Map Basics

### 1.1 Definition and Initialization

```go
var m map[string]int // Declared but not initialized, m is nil and cannot be written to
m := make(map[string]int)         // Initialize with make
m2 := make(map[string]int, 10)    // Specify initial capacity
m3 := map[string]int{}            // Literal initialization
m4 := map[string]int{"x": 1, "y": 2}
```
- Uninitialized maps cannot be written to, otherwise panic
- make initializes an empty map; capacity only optimizes performance

---

## 2. Add, Delete, Update, Query

### 2.1 Insert / Update

```go
m["a"] = 100 // Add or overwrite key
```

### 2.2 Query

```go
v := m["a"]          // Key exists → returns value
v = m["not_exist"]   // Key does not exist → returns zero value
v, ok := m["a"]      // Recommended: use ok to check if key exists
```

### 2.3 Delete

```go
delete(m, "a") // Delete key; no error if key doesn't exist
```

### 2.4 Iterate

```go
for k, v := range m {
    fmt.Println(k, v)
}
```

Nested map iteration:

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

## 3. Using Nested Maps

- Before writing, ensure the inner map is initialized, otherwise panic
- Reading a non-existent key returns zero value, no error
- Use ok to check if key exists

```go
if student[name] == nil {
    student[name] = make(map[string]float64)
}
student[name][sub] = score
```

---

## 4. Zero Value and nil

- Value types (int, float64, bool, struct): non-existent key returns zero value
- Reference types (slice, map, pointer, channel, interface): non-existent key returns nil
- Use ok to check if key exists; checking if value is nil only works for reference types

---

## 5. Concurrency Safety

- Go maps are not concurrency-safe; multiple goroutines reading/writing require locking, otherwise panic

### 5.1 Locking Methods

```go
var mu sync.Mutex
mu.Lock()
m[key] = value
v := m[key]
mu.Unlock()
```

For read-heavy scenarios, use RWMutex:

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

## 6. Common Mistakes

- Writing to a nil map causes panic
- Concurrent read/write causes panic
- Writing to uninitialized nested map causes panic

---

## 7. Advanced Usage: Encapsulated Management Structure

Thread-safe map manager example:

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

## 8. Practical Summary & Experience

- Maps must be initialized after declaration (make or literal)
- Initialize inner map before writing to nested map
- Use ok to check if key exists
- In concurrent scenarios, always use locks or sync.Map

---

## 9. Advanced Tips & Underlying Principles

### 9.1 Key Type Restrictions

- Only comparable types (int, string, bool, pointer, array, etc.) can be used as keys
- Cannot use slice, map, or function as keys
- Using float as key has precision risks; recommend converting to string or integer

### 9.2 Map Iteration Order

- Iteration order with range is random
- For fixed order, sort the key list

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

### 9.3 Capacity & Performance Optimization

- make(map[string]int, n) can specify initial capacity to reduce reallocation overhead

### 9.4 sync.Map Usage

Go 1.9+ provides thread-safe map:

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

### 9.5 Map Memory Management & GC

- After deleting a key, the value reference is GC'd, but bucket structure is not immediately released
- For large maps, after deleting many elements, create a new map and copy valid data to reduce memory usage

---

## 10. Summary

- Map is an efficient key-value data structure; pay attention to initialization and concurrency safety
- Initialize inner map before writing to nested maps
- Use locks or sync.Map for concurrent read/write
- Deleting a key only releases the value reference; bucket structure still occupies memory
- For large maps, after batch deletion, create a new map to optimize memory
- Key type must