# 🚀 Redis 学习笔记：Hash 类型详解与 Golang 应用

Redis 的 Hash 类型是用来存储多个 Field-Value 对的集合。它非常类似于编程语言中的 Map 或字典结构，特别适合用来存储对象（Object）或记录（Record）。

## 1. Hash 类型概览

### 1.1 数据结构

Hash 类型本质上是 Key-Value 存储结构中的 Value 部分，这个 Value 又是一个包含了多个 Field-Value 对的 Map 结构。

- **外部 Key**： 用于标识整个 Hash 结构（例如 `user:1001`）。
- **内部 Field-Value 对**： 存储对象的属性和值（例如 `name: Alice`, `age: 18`）。

| 外部 Key | 内部 Field | 内部 Value |
| :--- | :--- | :--- |
| `user:1001` | `name` | `dafei` |
| | `age` | `18` |

### 1.2 Golang/Java 类比

在概念上，Redis 的 Hash 类型类似于 Golang 中的 `map[string]map[string]string` 或 Java 中的 `Map<String, Map<String, String>>`。

### 1.3 内存优势

相较于将整个对象序列化为 JSON 字符串存储在 Redis String 类型中，使用 Hash 类型存储对象通常能占用更少的内存空间，因为 Redis 内部会根据存储的 Field-Value 数量和长度，选择不同的编码方式（如 `ziplist` 或 `hashtable`）进行优化。

## 2. Hash 类型常用命令

| 命令格式 | 功能描述 | 示例 | 备注 |
| :--- | :--- | :--- | :--- |
| `HSET key field value [field value ...]` | 将一个或多个 Field-Value 对缓存到 Hash 中。 | `HSET user name Alice age 100` | Redis 4.0+ 支持同时设置多个 Field。 |
| `HGET key field` | 从 Hash 列表中获取指定 Field 的值。 | `HGET user name` | |
| `HEXISTS key field` | 判断 Hash 列表中是否存在指定 Field 字段。 | `HEXISTS user age` | 返回 1 (存在) 或 0 (不存在)。 |
| `HDEL key field [field ...]` | 删除 Hash 列表中一个或多个 Field 字段。 | `HDEL user age` | |
| `HINCRBY key field increment` | 给 Hash 列表中指定 Field 字段原子性地增加指定增量。 | `HINCRBY user age 10` | 字段值必须是整数，常用于对象属性计数。 |
| `HLEN key` | 查看 Hash 列表中 Field 的数量。 | `HLEN user` | |
| `HKEYS key` | 获取 Hash 列表中的所有 Field 名称。 | `HKEYS user` | |
| `HVALS key` | 获取 Hash 列表中所有 Field 对应的 Value 值。 | `HVALS user` | |
| `HGETALL key` | 获取 Hash 列表中所有的 Field 及其对应的 Value 值。 | `HGETALL user` | **注意**： Hash 很大时，`HGETALL` 可能会阻塞服务器。 |

## 3. Hash 类型应用场景

Hash 结构最主要的优势在于可以独立地对对象的单个属性进行操作，无需读取和重写整个对象。

### 3.1 存储对象信息（首选方案）

这是 Hash 类型最常见的应用。与 String 类型存储 JSON 字符串相比，Hash 结构在更新对象局部属性时效率更高。

| 方案 | 存储方式 | 局部更新操作 | 性能/便捷性 |
| :--- | :--- | :--- | :--- |
| **方案 1 (String + JSON)** | `user_token` → `"{name:dafei, age:18, password:666}"` | 1. `GET` 整个 JSON。<br>2. 解析 JSON。<br>3. 修改属性。<br>4. 序列化为 JSON。<br>5. `SET` 整个 JSON。 | 侧重于查；改非常麻烦且耗性能。 |
| **方案 2 (Hash)** | `user_token` → `{name:dafei, age:18, password:666}` | 1. `HSET user_token age 19`。 | 侧重于改；查询整个对象时需用 `HGETALL`。 |

**结论**： 对于需要频繁更新对象单个属性（例如用户积分、年龄、登录状态）的场景，Hash 是更优的选择。

### 3.2 共享 Session 设计（进阶）

Hash 类型是实现集中式 Session 管理的理想选择，因为它能够轻松地管理用户对象的多个属性。

**登录缓存示例 (Golang Struct 对应 Hash)：**

```go
// Golang Struct 结构
type User struct {
    Username string
    Password string
    Age      int
}

// 对应 Redis Hash 结构
// Key: user:token:<TOKEN_ID>
// Fields: "username", "password", "age"
```

**操作流程：**

- **登录成功**：
  ```redis
  // 使用 HSET 存储用户对象的所有属性
  HSET user:token:xxx username dafei password 666 age 18
  ```
- **更新年龄（用户生日）**： 只需要操作一个 Field，无需读取和重写整个 Session。
  ```redis
  HINCRBY user:token:xxx age 1
  ```
- **获取所有信息**：
  ```redis
  HGETALL user:token:xxx
  ```

## 4. Golang 客户端操作 Hash 示例

在 Golang 中，使用 `go-redis/redis/v9` 库操作 Hash 类型：

```go
package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
// 假设 rdb 已经初始化并连接成功

func hashExample(rdb *redis.Client) {
    userKey := "user:profile:alice"

    // 1. HSET 批量设置字段
    rdb.HSet(ctx, userKey, map[string]interface{}{
        "name":     "Alice",
        "age":      25,
        "location": "Wonderland",
    })
    fmt.Println("HSET 成功设置 Alice 的 profile")

    // 2. HGET 获取单个字段
    age, err := rdb.HGet(ctx, userKey, "age").Result()
    if err != nil {
        panic(err)
    }
    fmt.Printf("HGET Alice 的年龄: %s\n", age) // 25

    // 3. HINCRBY 增加年龄
    newAge, err := rdb.HIncrBy(ctx, userKey, "age", 1).Result()
    if err != nil {
        panic(err)
    }
    fmt.Printf("HINCRBY 后的年龄: %d\n", newAge) // 26

    // 4. HGETALL 获取所有字段及其值
    data, err := rdb.HGetAll(ctx, userKey).Result()
    if err != nil {
        panic(err)
    }
    fmt.Printf("HGETALL 所有数据: %+v\n", data) 
    // 输出: map[age:26 location:Wonderland name:Alice]

    // 5. HDEL 删除字段
    rdb.HDel(ctx, userKey, "location")
    fmt.Println("HDEL 删除了 location 字段")
}
```