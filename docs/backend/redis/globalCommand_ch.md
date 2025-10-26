# 🚀Redis 全局命令学习笔记

🧩 一、概述

全局命令的作用范围覆盖所有数据库中的 key，对 Redis 的整体状态进行查看、修改或维护。
它们常用于：

- 系统监控（如查看 key 数量、数据库信息）

- 资源管理（如清空数据库、移动 key）

- 配置调整（如查看配置信息）

- 测试与调试（如 echo、randomkey）

🗝️ 二、常用全局命令速查表
| 命令格式 | 功能描述 | 示例 |
| :--- | :--- | :--- |
| `keys pattern` | 按照 pattern 匹配规则列出所有 key（⚠️ 不建议生产使用） | `keys user:*` |
| `exists key` | 判断指定 key 是否存在 | `exists name` |
| `expire key seconds` | 设置 key 的过期时间（秒） | `expire name 10` |
| `persist key` | 取消 key 的过期时间，使其永久有效 | `persist name` |
| `select index` | 切换数据库（默认 0，共 16 个：0～15） | `select 1` |
| `move key db` | 将 key 从当前数据库移动到指定数据库 | `move player 1` |
| `randomkey` | 随机返回一个 key | `randomkey` |
| `rename key newkey` | 将 key 改名为 newkey | `rename name newname` |
| `echo message` | 打印消息（常用于测试连接） | `echo Hello Redis` |
| `dbsize` | 查看当前数据库的 key 总数 | `dbsize` |
| `info` | 查看 Redis 运行状态与信息 | `info` |
| `config get *` | 查看所有 Redis 配置信息 | `config get *` |
| `flushdb` | 清空当前数据库 | `flushdb` |
| `flushall` | 清空所有数据库 | `flushall` |

⚠️ 三、命令详解与注意事项
1️⃣ **`keys pattern`**
```shell
keys *
keys s*
```

按 pattern 匹配所有 key。

- `*` 表示匹配任意字符，`?` 表示匹配一个字符，`[]` 表示匹配范围。

- **不推荐在生产环境使用**：数据量大时耗时长，会阻塞其他命令执行。
- ✅ 建议使用 `SCAN` 命令替代。

2️⃣ **`exists key`**
```shell
exists player
```

判断 key 是否存在。

返回：

- `1`：存在
- `0`：不存在

3️⃣ **`expire` / `persist`**
```shell
expire player 100
persist player
```

- `expire`：设置过期时间（秒），到期自动删除 key。
- `persist`：取消过期时间，让 key 永久有效。

4️⃣ **`select index`**
```shell
select 1
```

切换到指定数据库（默认 0 号库）。

- Redis 默认提供 16 个逻辑数据库：编号范围 0～15。
- 每个数据库之间相互独立。

5️⃣ **`move key db`**
```shell
move player 1
```

将当前数据库中的 key 移动到指定数据库。

- 如果目标数据库已有同名 key，移动失败（返回 0）。

6️⃣ **`randomkey`**
```shell
randomkey
```

随机返回一个 key。

- 若数据库为空，返回 `nil`。
- 常用于测试和抽样分析。

7️⃣ **`rename key newkey`**
```shell
rename player newkey
```

修改 key 的名称。

- 若新 key 已存在，则覆盖原内容。
- 原子操作，执行期间不会被中断。

8️⃣ **`echo message`**
```shell
echo "Hello Redis"
```

返回输入的字符串。

- 常用于测试连接是否正常。

9️⃣ **`dbsize`**
```shell
dbsize
```

返回当前数据库的 key 数量。

- 统计速度快，不会阻塞其他命令。

🔟 **`info`**
```shell
info
```

查看 Redis 的运行状态信息。

可加参数查看特定模块：

- `info memory`：查看内存信息
- `info stats`：查看统计信息
- `info replication`：查看主从状态

1️⃣1️⃣ **`config get *`**
```shell
config get *
```

查看所有 Redis 配置项。

- 可用于定位运行参数、内存策略、持久化设置等。

1️⃣2️⃣ **`flushdb` / `flushall`**
```shell
flushdb
flushall
```

- `flushdb`：清空当前数据库。
- `flushall`：清空所有数据库。
- ⚠️ **危险操作！**执行后数据不可恢复。
- 生产环境建议开启 `--protected-mode` 或增加 `requirepass` 防护。

🧠 四、实战示例：Redis 数据库管理流程
```shell
# 1. 切换到第1个数据库
select 1

# 2. 创建一个键
set player Tom

# 3. 查看所有键
keys *

# 4. 判断键是否存在
exists player

# 5. 设置过期时间
expire player 120

# 6. 查看数据库 key 数量
dbsize

# 7. 清空数据库
flushdb
```

🧾 五、总结
| 场景 | 推荐命令 | 说明 |
| :--- | :--- | :--- |
| 查看数据库状态 | `info`, `dbsize` | 快速获取运行信息 |
| 管理 key 生命周期 | `expire`, `persist` | 控制数据过期 |
| 调试与测试 | `echo`, `randomkey` | 验证连接与数据 |
| 配置管理 | `config get *` | 查看运行参数 |
| 数据清理 | `flushdb`, `flushall` | 清空数据库（慎用） |