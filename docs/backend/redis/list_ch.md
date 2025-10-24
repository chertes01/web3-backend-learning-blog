# 🚀 Redis 学习笔记：List 类型详解与 Golang 应用

Redis 的 List 类型是一个有序的字符串元素集合，其底层实现是双向链表 (Linked List)（在 Redis 3.2 以前是 `ziplist` 或 `linkedlist`，新版本使用 `quicklist` 优化了内存和性能）。List 类型的设计非常精巧，它既可以作为栈 (Stack)，也可以作为队列 (Queue)，满足了大量应用场景的需求。

## 1. List 类型概览

### 1.1 数据结构特性
- **有序性**： 元素按照插入顺序排列。
- **重复性**： 列表中可以包含重复的元素。
- **双端操作**： 支持从列表的左侧（头部）和右侧（尾部）进行元素的添加和移除，操作速度均为 $O(1)$。

### 1.2 Golang/Java 类比

| Redis List 操作 | Golang 概念 | Java 概念 |
| :--- | :--- | :--- |
| `RPUSH`/`LPOP` | FIFO 队列 (Queue) | `Queue`/`LinkedList` |
| `LPUSH`/`LPOP` | LIFO 栈 (Stack) | `Stack`/`Deque` |
| `LRANGE`/`LINDEX` | 普通列表 (List) | `List`/`Slice` |

## 2. List 类型常用命令

| 命令格式 | 功能描述 | 示例 | 栈/队列角色 |
| :--- | :--- | :--- | :--- |
| `RPUSH key value [v...]` | 从右边（尾部）往集合中添加值。 | `RPUSH Skill html golang` | 队列入队/列表追加 |
| `LPUSH key value [v...]` | 从左边（头部）往集合中添加值。 | `LPUSH Skill JavaScript` | 栈入栈/列表前置 |
| `LPOP key` | 弹出并返回集合中最左边（头部）的数据。 | `LPOP Skill` | 队列出队/栈出栈 |
| `RPOP key` | 弹出并返回集合中最右边（尾部）的数据。 | `RPOP Skill` | |
| `LRANGE key start stop` | 获取指定范围内的元素。 | `LRANGE Skill 0 -1` | 列表查询 |
| `LLEN key` | 获取列表的长度。 | `LLEN Skill` | |

## 3. List 类型非常用命令

| 命令格式 | 功能描述 | 示例 | 备注 |
| :--- | :--- | :--- | :--- |
| `LINSERT key BEFORE\|AFTER pivot value` | 在列表中指定 `pivot` 元素之前/之后插入 `value`。 | `LINSERT Skill before golang CSS` | 如果 `pivot` 不存在，不执行操作。 |
| `LSET key index value` | 更新索引 `index` 位置的值为 `value`。 | `LSET Skill 5 Solidity` | $O(N)$ 操作，应避免对大列表频繁使用。 |
| `LREM key count value` | 从列表中移除 `count` 个等于 `value` 的元素。 | `LREM Skill 3 Solidity` | `count > 0` (从左移除)，`count < 0` (从右移除)，`count = 0` (移除所有)。 |
| `LTRIM key start stop` | 截取列表，只保留 `start` 到 `stop` 范围内的元素。 | `LTRIM Skill 3 5` | **重要**：常用于限制列表长度（如最新N条数据）。 |
| `LINDEX key index` | 获取索引为 `index` 位置的数据。 | `LINDEX Skill 0` | $O(N)$ 操作，应避免对大列表频繁使用。 |

## 4. 经典应用场景

### 4.1 消息队列 (Message Queue)

List 是实现简单消息队列的天然选择，具有 **原子性** 和 **阻塞操作** 的特性。

- **生产者 (Producer)**： 使用 `LPUSH` 或 `RPUSH` 将消息推入列表。
- **消费者 (Consumer)**： 使用 `LPOP` 或 `RPOP` 取出消息。
- **阻塞队列**： 使用 `BLPOP` / `BRPOP` 实现阻塞式弹出。当队列为空时，消费者会阻塞等待新消息，直到超时或有新消息到来，有效避免了 CPU 空轮询。

**命令示例（阻塞消费）：**
```redis
// BRPOP key [key ...] timeout
BRPOP task:queue 0  // 阻塞等待，直到有消息或服务器关闭
```

### 4.2 栈与最新列表 (Feed/Timeline)

List 可用于存储用户动态或最新消息流，如微博 Timeline。

- **场景**： 用户收藏文章列表、最新购买记录、操作日志。
- **实现**： 使用 `LPUSH` 保证最新添加的元素始终在列表的头部（索引 0）。

**命令示例：**
```redis
// 添加最新文章 ID 到列表头部
LPUSH user:favor:article:1001 aid_300
LPUSH user:favor:article:1001 aid_299

// 查看最近的 10 篇文章 ID
LRANGE user:favor:article:1001 0 9

// 配合 LTRIM 限制列表长度，只保留最新的 1000 条
LTRIM user:favor:article:1001 0 999
```

## 5. Golang 客户端操作 List 示例

在 Golang 中，使用 `go-redis/redis/v9` 库操作 List 类型：

```go
package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
// 假设 rdb 已经初始化并连接成功

func listExample(rdb *redis.Client) {
    listKey := "task:queue:sync"
    
    // 1. RPUSH/LPUSH 添加入队
    rdb.RPush(ctx, listKey, "task_1", "task_2", "task_3")
    rdb.LPush(ctx, listKey, "task_0") 
    // 队列当前状态: ["task_0", "task_1", "task_2", "task_3"]

    // 2. LRANGE 查看列表所有元素
    all, _ := rdb.LRange(ctx, listKey, 0, -1).Result()
    fmt.Printf("当前队列: %v\n", all)

    // 3. LPOP/RPOP 出队
    leftPop, _ := rdb.LPop(ctx, listKey).Result()
    rightPop, _ := rdb.RPop(ctx, listKey).Result()
    fmt.Printf("LPOP (左出): %s, RPOP (右出): %s\n", leftPop, rightPop)
    // 队列当前状态: ["task_1", "task_2"]

    // 4. BLPOP 阻塞弹出 (演示：在实际应用中通常是单独的消费者服务调用)
    // rdb.BLPop(ctx, 5*time.Second, listKey) 
    // 阻塞 5 秒等待新任务
    
    // 5. LTRIM 限制长度
    rdb.RPush(ctx, listKey, "a", "b", "c", "d", "e") // 队列变长
    rdb.LTrim(ctx, listKey, 0, 3) // 只保留前 4 个元素 (0, 1, 2, 3)
    trimmed, _ := rdb.LRange(ctx, listKey, 0, -1).Result()
    fmt.Printf("LTRIM 后队列: %v\n", trimmed)
}
```
