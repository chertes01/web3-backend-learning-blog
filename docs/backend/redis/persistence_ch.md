# Redis 持久化机制详解（Redis 7.0.15）

> 📘 **作者**：Chertes von
> 🕒 **更新时间**：2025-10
> 💡 **说明**：本文基于 Redis 7.0.15，演示默认持久化行为、配置解析与数据丢失原因。

---

## 🧩 一、引言：为什么要了解持久化

Redis 是一个基于内存的高性能数据库。所有数据默认存放在内存中，因此一旦重启或断电，数据就可能全部丢失。

为了在重启后还能恢复数据，Redis 提供了三种持久化机制：

- **RDB（快照）**
- **AOF（追加日志）**
- **RDB + AOF 混合持久化**（Redis 4.0+）

---

## 🧪 二、小实验：数据为何会“随机消失”

#### 实验步骤

1.  设置两个键：

    ```shell
    127.0.0.1:6379> SET aa aa
    OK
    127.0.0.1:6379> SET bb bb
    OK
    127.0.0.1:6379> KEYS *
    1) "aa"
    2) "bb"
    ```

2.  然后重启 Redis：

    ```bash
    sudo systemctl restart redis
    ```

3.  再次查看：

    ```shell
    127.0.0.1:6379> KEYS *
    (empty array)
    ```

#### 实验结果可能有以下几种：

| 现象 | 说明 |
| :--- | :--- |
| ✅ `aa`、`bb` 都在 | 数据已被持久化 |
| ⚠️ `aa` 在，`bb` 不在 | 部分数据未写入磁盘 |
| ❌ 全部不在 | Redis 没有执行任何持久化操作 |

**思考**：为什么行为“随机”？其实它并不随机，而是由持久化机制的触发条件决定的。

---

## ⚙️ 三、Redis 7.0.15 的默认持久化配置

你可以通过命令查看当前配置：

- **RDB 配置**：

  ```shell
  127.0.0.1:6379> CONFIG GET save
  1) "save"
  2) "3600 1 300 100 60 10000"
  ```

  这表示 RDB 持久化是**开启**的。

- **AOF 配置**：

  ```shell
  127.0.0.1:6379> CONFIG GET appendonly
  1) "appendonly"
  2) "no"
  ```

  这表示 AOF 持久化是**关闭**的。

---

## 📁 四、默认配置解析

| 配置规则 | 触发条件 | 含义 |
| :--- | :--- | :--- |
| `save 3600 1` | 1 小时内至少 1 次写操作 | 每小时至少有变更时保存快照 |
| `save 300 100` | 5 分钟内至少 100 次写操作 | 写入频繁时提前保存 |
| `save 60 10000` | 1 分钟内至少 10000 次写操作 | 写入极多时快速保存 |
| `appendonly no` | 不启用 AOF | 不追加写入日志 |

#### 🔍 文件位置

默认情况下：

- `dbfilename dump.rdb`
- `dir /var/lib/redis/`

所以 RDB 文件一般位于：`/var/lib/redis/dump.rdb`

---

## 🧠 五、为什么会“随机丢数据”

| 场景 | 持久化方式 | 丢失原因 |
| :--- | :--- | :--- |
| 写完立刻重启 | RDB | 快照未触发保存 |
| 短时间大量写入 | RDB | 不满足写次数条件 |
| 异常宕机 | AOF 关闭 | 没有追加日志 |
| 部分键存在 | RDB/AOF 混合 | AOF 文件未完全写入或未加载 |

Redis 在重启时优先加载顺序如下：

1.  `appendonly.aof`（若存在）
2.  `dump.rdb`（若无 AOF）
3.  空数据库（都不存在）

---

## 🧾 六、验证快照状态

使用 `INFO Persistence` 命令查看 RDB 相关状态：

```shell
redis-cli INFO Persistence | grep rdb
```

**输出示例**：

```
rdb_changes_since_last_save:10
rdb_bgsave_in_progress:0
rdb_last_save_time:1735192443
rdb_last_bgsave_status:ok
```

| 字段 | 含义 |
| :--- | :--- |
| `rdb_last_save_time` | 上次保存时间（Unix 时间戳） |
| `rdb_changes_since_last_save` | 自上次保存后修改的 key 数量 |
| `rdb_last_bgsave_status` | 最近一次快照结果 |
| `rdb_bgsave_in_progress` | 当前是否在执行快照 |

---

## 🧮 七、三种持久化机制对比

| 持久化方式 | 文件名 | 持久化时机 | 性能 | 数据安全性 | Redis 7 默认状态 |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **RDB** | `dump.rdb` | 定期生成快照 | 快 | ⚠️ 可能丢失最近数据 | ✅ **开启** |
| **AOF** | `appendonly.aof` | 每次写命令追加日志 | 中 | ✅ 可实时保存 | ❌ **关闭** |
| **RDB + AOF 混合** | 同时存在 | 综合两者优点 | 较快 | ✅ 高可靠 | ❌ **关闭** |

---

## 💡 八、推荐配置（开发/生产通用）

```conf
# 开启 RDB 快照
save 900 1
save 300 10
save 60 10000

# 开启 AOF 持久化
appendonly yes
appendfsync everysec     # 每秒刷盘，性能与安全平衡
auto-aof-rewrite-percentage 100
auto-aof-rewrite-min-size 64mb
```

**📌 建议**：

- **开发环境**：可以只用 RDB。
- **生产环境**：**RDB + AOF 混合** 最安全。

---

## 🧭 九、总结

| 特性 | Redis 默认行为（7.0.15） |
| :--- | :--- |
| **RDB** | ✅ **开启**，周期性触发（3600s/300s/60s） |
| **AOF** | ❌ **默认关闭** |
| **混合模式** | ❌ **默认关闭** |
| **数据丢失** | ⚠️ 快照未触发时，重启后数据丢失 |
| **快照文件位置** | `/var/lib/redis/dump.rdb` |

**🧠 结论**：
Redis 7 默认启用了 RDB 快照机制，但保存周期较长。如果你写入数据后立即重启 Redis，很可能看起来“随机丢数据”，实际上是快照机制还没触发。

---

## 📚 十、附录：快速命令参考

| 操作 | 命令 |
| :--- | :--- |
| 手动保存快照 | `SAVE`（阻塞） / `BGSAVE`（后台异步） |
| 手动触发 AOF 重写 | `BGREWRITEAOF` |
| 查看当前持久化配置 | `CONFIG GET save` / `CONFIG GET appendonly` |
| 查看快照状态 | `INFO Persistence` |
| 禁用 RDB | `CONFIG SET save ""` |

---

## 🎨 十一、推荐图示（可选）

如果要发布到 `README.md`，可附加下图说明持久化机制关系：

```
内存数据变化
      ↓
  ┌─────────┐
  │  RDB快照 │──► dump.rdb（周期保存）
  └─────────┘
      ↓
  ┌─────────┐
  │  AOF日志 │──► appendonly.aof（命令记录）
  └─────────┘
      ↓
  Redis重启时按优先级恢复
      (AOF > RDB)
```
