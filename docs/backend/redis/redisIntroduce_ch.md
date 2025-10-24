# 🚀 Redis 学习笔记：高性能键值存储与 Golang 应用

## 1. 概览：NoSQL 数据库分类与 Redis 定位

Redis（Remote Dictionary Server）是一个开源的、基于内存的键值（Key-Value）存储数据库，它常被用作缓存、消息队列和数据库。

### 1.1 NoSQL 数据库特性对比（基于图片内容）

| 分类 | 示例（Examples） | 典型应用场景 | 数据模型 | 优点 | 缺点 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **键值 (Key-Value) 数据库** | Redis, Tokyo Cabinet, Voldemort, Oracle BDB | 内容缓存、高访问负载、日志系统 | Key 指向 Value 的键值对，常用 Hash Table 实现 | 查找速度快 | 数据非结构化，通常被视为字符串或二进制数据 |
| **列存储数据库** | Cassandra, HBase, Riak | 分布式文件系统 | 以列式存储，将同一列数据存在一起 | 查找速度快，易于分布式扩展 | 功能局部受限 |
| **文档型数据库** | CouchDB, MongoDB | Web 应用 | Key-Value 对应，Value 是结构化数据，无需预定义结构 | 架构灵活，查询灵活 | 查询性能不高，缺乏统一查询语法 |
| **图形 (Graph) 数据库** | Neo4J, InfoGrid, Infinite Graph | 社交网络、推荐系统，专注于构建关系图谱 | 图结构 | 利用图结构相关算法（如最短路径、N 度关系查询） | 很多计算需要遍历整个图，不适合做分布式集群 |

### 1.2 Redis 简介与定位（基于文字与图片内容）

- **Redis 定位**：Redis 是以 Key-Value 形式存储的非关系型 (NoSQL) 数据库。它通常被定位为缓存，用于提高数据读写速度，减轻对后端数据库（如 MySQL）的存储与访问压力。
- **核心特性**：
    - 非关系型、分布式、开源、水平可拓展。
    - **基于内存存储**：实现对数据高并发读写，高效率存储和访问。
    - **单线程模型 (Redis 6 以前)**：保证了操作的原子性，避免了并发问题带来的锁开销（Redis 4.0 引入了多线程来处理I/O操作，Redis 6.0 以后正式引入了多线程 I/O 模型来提高性能，但核心命令的执行仍是单线程原子操作）。
    - 高可用性与可拓展性。

- **优缺点**：

| 优点 (Advantages) | 缺点 (Disadvantages) |
| :--- | :--- |
| 对数据高并发读写（基于内存） | ACID 处理简单，无法支持太复杂的关系数据库模型 |
| 对海量数据的高效率存储和访问 | |
| 数据的可拓展性和高可用性 | |
| 操作原子性（核心命令执行层面） | |

## 2. Redis 支持的数据类型

Redis 支持丰富的数据结构，使得它能应用于多种场景。

### 2.1 常用数据类型 (Common Types)

- **String (字符串)**
    - **概述**：最简单类型，一个 key 对应一个 value。Value 可以是字符串、整数或浮点数。
    - **典型应用场景**：缓存对象、计数器（如点赞数、访问量）、简单键值存储。
- **Hash (哈希)**
    - **概述**：键值对集合，适合存储结构化对象。一个 key 对应一个 field-value 的集合。
    - **典型应用场景**：存储用户信息（如用户ID作为Key，姓名、年龄作为Field）。
- **List (列表)**
    - **概述**：字符串列表，按插入顺序排序。可以从头部或尾部添加或移除元素。底层实现是链表。
    - **典型应用场景**：消息队列（LPUSH/RPOP 实现生产者/消费者模式）、时间线/最新消息列表。
- **Set (集合)**
    - **概述**：无序的字符串集合，元素唯一。支持集合间的操作，如交集、并集、差集。
    - **典型应用场景**：标签（Tags）系统、抽奖活动（随机从集合中取元素）、社交关系（共同好友）。
- **ZSet (Sorted Set - 有序集合)**
    - **概述**：类似于 Set，但每个成员都关联一个分数 (score)，元素按分数从小到大排序，分数相同时按成员的字典序排序。元素唯一。
    - **典型应用场景**：排行榜（如游戏积分榜、热搜榜）、带权重的任务队列。

### 2.2 其他数据类型 (Less Common Types)

- **HyperLogLog**：用于基数估计算法（如统计 UV）。
- **Bitmap (位图)**：用于布尔值存储和计算（如用户签到、活跃用户统计）。
- **Geospatial (地理位置)**：用于存储地理坐标并进行范围查找。
- **Streams (流信息)**：类似于消息队列，支持多消费者组。

## 3. Redis 命令格式与 Golang 类比

Redis 命令的基本格式：`类型命令 key 参数`
```
# 示例：
set name dafei
```
**操作建议**：Redis 的操作类似于 Golang 中的 `map[string]interface{}` 集合，都是 key-value 形式存储数据。Key 大部分情况下为 String 类型。Value 则根据存储的数据结构，可以是 String、Hash、List、Set、ZSet 等类型。

### 3.1 Golang 中的 Redis 客户端

在 Golang 中，我们通常使用成熟的第三方库来操作 Redis，例如 `go-redis/redis/v9`。

**示例：基本连接**
```go
import (
	"context"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func initRedis() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 地址
		Password: "",               // 密码 (如果有)
		DB:       0,                // 数据库索引
		PoolSize: 20,               // 连接池大小
	})
    // 检查连接是否成功
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}

// 示例：String 类型操作
func stringExample(rdb *redis.Client) {
    // 设置值
	err := rdb.Set(ctx, "mykey", "hello redis", 0).Err() 
    if err != nil {
        panic(err)
    }

    // 获取值
	val, err := rdb.Get(ctx, "mykey").Result()
	if err == redis.Nil {
		// Key 不存在
	} else if err != nil {
		panic(err)
	} else {
		// 获取成功
        // ... 使用 val
	}
}
```

**进阶应用**：Golang 客户端通常支持 Redis 的高级特性，例如：
- **Pipeline (管道)**：批量发送命令，减少网络延迟。
- **Transactions (事务)**：使用 `MULTI` 和 `EXEC` 保证一组命令原子性执行（非传统数据库意义上的事务）。
- **Pub/Sub (发布/订阅)**：实现消息传递。
- **Lua Scripting (Lua 脚本)**：保证复杂操作的原子性并减少网络开销。
- **Cluster/Sentinel**：支持高可用和分布式部署模式。
