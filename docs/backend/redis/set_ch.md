# 🚀 Redis 学习笔记：Set 类型详解与 Golang 应用

Redis 的 Set（集合）类型是一个 **无序** 的字符串元素集合。其核心特点是元素的唯一性，底层通过 **哈希表（Hash Table）** 实现，因此元素的添加、删除、查找操作的时间复杂度均为 $O(1)$。Set 类型是 Redis 中用于处理数学集合运算（交集、并集、差集）最理想的数据结构。

## 1. Set 类型概览

### 1.1 数据结构特性
- **无序性**： 元素没有固定的顺序。
- **唯一性**： 集合中不允许出现重复的成员（Members）。
- **高效操作**： 添加、移除、判断成员是否存在都非常快 ($O(1)$)。

### 1.2 Golang/Java 类比

| Redis Set 操作 | Golang 概念 | Java 概念 |
| :--- | :--- | :--- |
| 整体结构 | `map[string]struct{}` (利用 map key 的唯一性) | `Set` / `HashSet` |
| 成员操作 | `len(set)` | `Set.add()`, `Set.remove()`, `Set.contains()` |

## 2. Set 类型常用命令

| 命令格式 | 功能描述 | 示例 | 备注 |
| :--- | :--- | :--- | :--- |
| `SADD key member [member ...]` | 往 Key 集合中添加一个或多个成员元素。 | `SADD user Alice Frank Tom` | 重复添加的元素会被忽略。 |
| `SMEMBERS key` | 遍历 Key 集合中的所有元素。 | `SMEMBERS user` | **注意**： 集合很大时，可能会阻塞服务器。 |
| `SREM key member [member ...]` | 删除 Key 集合中的一个或多个成员元素。 | `SREM user Tim` | |
| `SPOP key [count]` | 从 Key 集合中随机弹出并返回 `count` 个元素。 | `SPOP user 1` | 弹出的元素将从集合中删除。 |
| `SCARD key` | 获取集合中元素的数量（Cardinality）。 | `SCARD user` | |
| `SISMEMBER key member` | 判断 `member` 元素是否在 Key 集合中。 | `SISMEMBER user Tim` | 返回 1 (是) 或 0 (否)。 |

## 3. Set 类型集合运算命令

Set 类型提供了一套完整的集合运算命令，这是其最强大的功能之一。

| 命令格式 | 功能描述 | 示例 | 场景 |
| :--- | :--- | :--- | :--- |
| `SDIFF key1 key2 [key3...]` | 返回 `key1` 中独有的元素（差集）。 | `SDIFF user student` | 查找“只关注用户但没有关注学生”的人。 |
| `SINTER key1 key2 [key3...]` | 返回所有给定集合的共有元素（交集）。 | `SINTER user student` | 查找“既是用户又是学生”的人（共同关注）。 |
| `SUNION key1 key2 [key3...]` | 返回所有给定集合的合并结果（并集）。 | `SUNION user student` | 查找“关注用户或关注学生”的所有人。 |
| `S*STORE dest key1 key2...` | 执行上述运算（`SDIFFSTORE`, `SINTERSTORE`, `SUNIONSTORE`），并将结果缓存到 `dest` 集合中。 | `SINTERSTORE desk user student` | 缓存集合运算结果，避免重复计算。 |

## 4. Set 类型非常用命令

| 命令格式 | 功能描述 | 示例 | 备注 |
| :--- | :--- | :--- | :--- |
| `SRANDMEMBER key [count]` | 随机获取 `count` 个元素，但不从集合中删除。 | `SRANDMEMBER user 3` | |
| `SMOVE source destination member` | 将 `source` 集合中的 `member` 元素原子性地移动到 `destination` 集合。 | `SMOVE user destination Tim` | 常用于状态转移。 |

## 5. 经典应用场景

### 5.1 数据去重

利用 Set 集合的唯一性，可以轻松实现数据去重。

- **场景**： 统计网站独立 IP 访问量（UV - Unique Visitor）、用户每日打卡去重。
- **实现**： 每次访问或打卡时，将用户的唯一标识（如 IP 或 UserID）使用 `SADD` 加入 Set 集合。最终使用 `SCARD` 获取集合大小，即为去重后的数量。

**命令示例：**
```redis
SADD unique:visitors:20251025 192.168.1.1 192.168.1.5
SCARD unique:visitors:20251025  // 获取今日 UV 数量
```

### 5.2 社交关系与共同好友

集合运算是处理社交关系的核心工具。

- **场景**： 共同关注、好友推荐、用户标签系统。
- **实现**： 每个用户关注的好友列表存储为一个 Set。

**命令示例：**
```redis
// 查找 user:alice 和 user:bob 的共同好友（交集）
SINTER user:alice:follows user:bob:follows

// 查找 user:alice 关注了但 user:bob 没有关注的人（差集）
SDIFF user:alice:follows user:bob:follows
```

### 5.3 随机抽奖

利用 `SPOP` 的随机性和弹出特性，可以高效安全地进行抽奖活动。

- **场景**： 抽取 N 等奖、随机分配任务。
- **实现**： 奖池或参与者列表使用 Set 存储，使用 `SPOP` 抽取。

**命令示例：**
```redis
// 1. 准备抽奖池
SADD luckydraw 1 2 3 ... 100

// 2. 抽取 3 个三等奖（并从奖池中移除）
SPOP luckydraw 3

// 3. 抽取 1 个一等奖
SPOP luckydraw 1
```

## 6. Golang 客户端操作 Set 示例

在 Golang 中，使用 `go-redis/redis/v9` 库操作 Set 类型：

```go
package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
// 假设 rdb 已经初始化并连接成功

func setExample(rdb *redis.Client) {
    userSet := "set:user:followers"
    studentSet := "set:student:members"
    
    // 1. SADD 添加成员
    rdb.SAdd(ctx, userSet, "Alice", "Frank", "Tom", "Tim")
    rdb.SAdd(ctx, studentSet, "Hans", "Trump", "Tom", "Tim")
    
    // 2. SISMEMBER 判断成员是否存在
    isMember, _ := rdb.SIsMember(ctx, userSet, "Alice").Result()
    fmt.Printf("Alice 是否在 userSet 中: %v\n", isMember) // true

    // 3. SINTER 计算交集
    commonMembers, _ := rdb.SInter(ctx, userSet, studentSet).Result()
    fmt.Printf("共同成员 (交集): %v\n", commonMembers) 
    // 输出: [Tom Tim]

    // 4. SDIFF 计算差集
    diffMembers, _ := rdb.SDiff(ctx, userSet, studentSet).Result()
    fmt.Printf("userSet 独有成员 (差集): %v\n", diffMembers) 
    // 输出: [Alice Frank]

    // 5. SPOP 随机弹出
    spopMember, _ := rdb.SPop(ctx, userSet).Result() // 随机弹出一个，并从集合中移除
    fmt.Printf("SPOP 弹出的成员: %s\n", spopMember)

    // 6. SCARD 获取数量
    cardinality, _ := rdb.SCard(ctx, userSet).Result()
    fmt.Printf("userSet 剩余成员数量: %d\n", cardinality)
}
```
