# 🚀 Redis 学习笔记：Sorted Set (ZSet) 类型详解与应用

Sorted Set（有序集合，简称 ZSet）是 Redis 的一种复合数据结构，它结合了 Set 的唯一性和 List 的有序性。它要求集合中的每个成员（Member）都关联一个 double 类型分数（Score），并根据这个分数从小到大排序。Sorted Set 是 Redis 中实现排行榜功能最核心、最高效的数据结构。

## 1. Sorted Set (ZSet) 类型概览

### 1.1 数据结构特性
- **有序性**： 成员按其关联的 Score 值从小到大排序。Score 相同时，按成员的字典序（ lexicographical order）排序。
- **唯一性**： 集合中不允许有重复的成员，但允许分数重复。
- **高效操作**： 添加、删除、查找、获取成员的分数、获取成员的排名等操作的平均时间复杂度都是 $O(\log N)$（底层通过跳跃表和哈希表实现）。

### 1.2 数据缓存结构

| 外部 Key | 内部 Score | 内部 Member |
| :--- | :--- | :--- |
| `player` | 100 | `Alice` |
| | 120 | `Frank` |
| | 150 | `Tom` |
| | 200 | `Bob` |

## 2. Sorted Set 常用命令

| 命令格式 | 功能描述 | 示例 | 排序方式 |
| :--- | :--- | :--- | :--- |
| `ZADD key score member [s m ...]` | 往集合中添加成员元素及其分数。 | `ZADD player 200 Bob 120 Frank` | |
| `ZINCRBY key increment member` | 将指定成员的分数原子性地增加 `increment`。 | `ZINCRBY player 50 Alice` | 排行榜积分更新。 |
| `ZRANGE key start stop [WITHSCORES]` | 按分数升序返回指定范围内的成员 [及分数]。 | `ZRANGE player 0 -1 WITHSCORES` | 正向（低分在前） |
| `ZREVRANGE key start stop [WITHSCORES]` | 按分数降序返回指定范围内的成员 [及分数]。 | `ZREVRANGE player 0 -1 WITHSCORES` | 反向（高分在前） |
| `ZRANK key member` | 返回成员在集合中的正序排名（从 0 开始）。 | `ZRANK player Alice` | |
| `ZREVRANK key member` | 返回成员在集合中的倒序排名（从 0 开始）。 | `ZREVRANK player Frank` | 用于显示排行榜名次。 |
| `ZCARD key` | 返回集合中元素的个数。 | `ZCARD player` | |

## 3. Sorted Set 非常用命令

| 命令格式 | 功能描述 | 示例 | 备注 |
| :--- | :--- | :--- | :--- |
| `ZRANGEBYSCORE key min max [WITHSCORES]` | 按分数范围 $\in [min, max]$ 升序返回成员。 | `ZRANGEBYSCORE player 120 200 WITHSC` | 可选参数 `LIMIT offset count` 进行分页。 |
| `ZREVRANGEBYSCORE key max min [WITHSCORES]` | 按分数范围 $\in [max, min]$ 降序返回成员。 | `ZREVRANGEBYSCORE player 200 120` | **注意**： `max` 在前，`min` 在后。 |
| `ZREM key member [member ...]` | 删除集合中的一个或多个成员及其分数。 | `ZREM player Tim` | |
| `ZREMRANGEBYSCORE key min max` | 按分数范围 $\in [min, max]$ 删除集合中元素。 | `ZREMRANGEBYSCORE player 120 200` | 用于清理过期或低分数据。 |
| `ZREMRANGEBYRANK key start stop` | 按正序排名范围 $\in [start, stop]$ 删除元素。 | `ZREMRANGEBYRANK player 2 4` | 用于保持排行榜的固定长度。 |
| `ZCOUNT key min max` | 统计分数范围 $\in [min, max]$ 内的元素个数。 | `ZCOUNT player 120 200` | |

## 4. 经典应用场景：排行榜

Sorted Set 是实现各种类型排行榜的完美选择，因为它同时解决了存储、去重、排序和动态更新四大问题。

- **场景**： 游戏积分榜、视频播放量排行榜、热点文章榜、用户等级。
- **实现**：
    - **Key**： 排行榜名称（如 `leaderboard:game_score`）。
    - **Score**： 积分、播放量等排序依据。
    - **Member**： 用户ID、视频ID 等唯一标识。

**操作示例：**

| 需求 | Redis 命令 |
| :--- | :--- |
| 用户积分更新 | `ZINCRBY leaderboard:game_score 100 user:101` |
| 获取前 10 名 | `ZREVRANGE leaderboard:game_score 0 9 WITHSCORES` |
| 获取用户 Bob 的排名 | `ZREVRANK leaderboard:game_score user:bob` |
| 获取分数在 1000 到 2000 之间的成员 | `ZRANGEBYSCORE leaderboard:game_score 1000 2000` |

## 5. 总结：Redis 在项目中的 Key-Value 设计思考

在项目中引入 Redis 并不仅是执行命令，更关键的是设计高效的 Key-Value 结构。

| 思考点 | 决策/设计原则 | 对应的 Redis 类型 |
| :--- | :--- | :--- |
| 1. 业务是否需要缓存？ | **决策**： 是 (用于加速读写、减轻数据库压力)。 | |
| 2. 是否选用 Redis？ | **决策**： 是 (高性能、多数据结构、高可用)。 | |
| 3. Key-Value 对设计 | **原则**： 关注数据的形态、关系和操作。 | |
| 简单 K-V 存储/计数器 | 存储单个值、原子计数。 | **String** |
| 存储对象/记录 | 存储对象的多个字段，且需要局部更新。 | **Hash** |
| 栈、队列、最新列表 | 有序、允许重复、需要从两端操作。 | **List** |
| 去重、集合运算（交/并/差） | 无序、元素唯一、需要处理社交关系。 | **Set** |
| 排行榜、带权重的列表 | 有序、元素唯一、需要按 Score 排序或按范围查找。 | **Sorted Set (ZSet)** |
