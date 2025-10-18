# MySQL 性能分析与优化

本指南通过分步讲解的方式，介绍如何分析并优化 MySQL 中的 SQL 查询性能。涵盖用于监控服务器负载、识别慢查询以及理解查询执行计划的各种工具。

---

## 目录

- [1. 宏观性能分析：SHOW STATUS](#1-宏观性能分析show-status)
- [2. 识别问题查询：慢查询日志](#2-识别问题查询慢查询日志)
- [3. 查询细节剖析：SHOW PROFILE](#3-查询细节剖析show-profile)
- [4. 执行计划：EXPLAIN](#4-执行计划explain)
- [5. 实战案例：索引的力量](#5-实战案例索引的力量)

---

## 1. 宏观性能分析：SHOW STATUS

在深入研究单个查询之前，对服务器的整体工作负载有一个宏观的了解是很有帮助的。`SHOW STATUS` 命令可以提供服务器的状态变量，让你了解不同类型操作的执行频率。这能帮你判断数据库是读密集型、写密集型，还是混合型负载。

### 全局状态 (Global) vs 会话状态 (Session)

- **全局状态 (Global Scope)**：显示自上次服务器重启以来的累计值，适合 DBA 监控服务器整体健康。
- **会话状态 (Session Scope)**：仅显示当前客户端连接的值，适合开发人员调试特定脚本或应用会话所产生的影响。

```sql
-- 查看整个 MySQL 服务器从启动以来的累计状态值
SHOW GLOBAL STATUS LIKE 'Com_______';

-- 查看当前连接会话的状态值
SHOW SESSION STATUS LIKE 'Com_______';
```

| 特性       | SHOW GLOBAL STATUS                     | SHOW SESSION STATUS                 |
|------------|----------------------------------------|-------------------------------------|
| 统计范围   | 整个 MySQL 服务器实例                  | 当前连接                            |
| 生命周期   | 从 MySQL 服务器启动开始累加            | 从客户端连接到断开                  |
| 重置条件   | MySQL 服务器重启                       | 当前连接断开                        |
| 典型用途   | DBA 监控服务器整体健康、负载、趋势等   | 开发人员调试特定应用或脚本性能      |

---

## 2. 识别问题查询：慢查询日志

慢查询日志 (Slow Query Log) 是一个用于捕获执行时间过长的 SQL 查询的基础工具。任何执行时间超过 `long_query_time` 阈值的查询都会被记录。

### 配置慢查询日志

- 检查日志功能是否开启：

```sql
SHOW VARIABLES LIKE 'slow_query_log';
```

- 动态开启/关闭日志：

```sql
SET GLOBAL slow_query_log = 1; -- 1 开启, 0 关闭
```

- 设置时间阈值（单位：秒）：

```sql
SET GLOBAL long_query_time = 2; -- 记录执行时间超过 2 秒的查询
```

- 持久化配置（my.cnf 或 my.ini）：

```ini
[mysqld]
slow_query_log = 1
slow_query_log_file = /var/log/mysql/mysql-slow.log
long_query_time = 2
```

> **注意**：将 `long_query_time` 设置为 0 会记录所有查询，适合调试，但可能导致日志文件暴增。

---

## 3. 查询细节剖析：SHOW PROFILE

`SHOW PROFILE` 命令可以将查询的执行时间分解到各个阶段（如发送数据、执行、解析等），帮助定位性能瓶颈。

> ⚠️ **注意**：`SHOW PROFILE` 从 MySQL 5.7 开始弃用，8.0 中移除，建议使用 Performance Schema 替代。

### 使用流程

1. 检查是否支持 profiling 功能：

```sql
SELECT @@have_profiling;
```

2. 开启 profiling：

```sql
SET profiling = 1;
```

3. 执行查询并查看 profiles：

```sql
SELECT * FROM tb_user;
SHOW PROFILES;
```

4. 查看指定查询的阶段耗时：

```sql
SHOW PROFILE FOR QUERY 5;
SHOW PROFILE CPU FOR QUERY 5;
```

示例输出：

| Query_ID | Duration   | Query                                                            |
|----------|------------|------------------------------------------------------------------|
| 1        | 0.00012000 | SHOW WARNINGS                                                   |
| 2        | 0.00044350 | SELECT * FROM tb_user                                           |
| 3        | 0.00033725 | SELECT * FROM tb_user WHERE id=30                               |
| 4        | 0.00157550 | SELECT * FROM tb_user WHERE id=30 OR id=26 OR name='白骨精'      |
| 5        | 13.6130575 | SELECT count(*) FROM tb_sku ts WHERE id > 0                     |

---

## 4. 执行计划：EXPLAIN

`EXPLAIN` 是查询优化中最重要的工具之一。它显示 MySQL 优化器选择的执行计划，而不实际执行查询。

```sql
EXPLAIN SELECT count(*) FROM tb_sku ts WHERE id > 0;
```

### EXPLAIN 字段解读

| 字段         | 说明                                                         |
|--------------|--------------------------------------------------------------|
| id           | 查询序号，决定执行优先级                                     |
| select_type  | 查询类型（SIMPLE、PRIMARY、SUBQUERY、DERIVED）               |
| table        | 操作的表名或别名                                             |
| partitions   | 匹配的分区                                                   |
| type         | 访问类型（system > const > eq_ref > ref > range > index > ALL）|
| possible_keys| 可能使用的索引                                               |
| key          | 实际使用的索引                                               |
| key_len      | 索引长度（字节），判断复合索引利用程度                       |
| ref          | 与索引比较的字段或常量                                       |
| rows         | 预估扫描的数据行数                                           |
| filtered     | 预估过滤百分比                                               |
| Extra        | 额外信息（Using index, Using where, Using filesort 等）      |

### 常见优化提示

- **type = ALL**：全表扫描，性能瓶颈，应优先优化（如添加索引）。
- **Extra = Using filesort**：表示无法利用索引完成排序，需额外排序步骤，建议优化。

---

## 5. 实战案例：索引的力量

### 场景：未索引列导致慢查询

假设需要根据产品序列号（`sn`）查找产品，但 `sn` 列没有索引。

```sql
-- 全表扫描，性能极差
SELECT * FROM tb_sku WHERE sn = '10000000314500389';
```

**性能表现**：执行时间约 19.4 秒。

**EXPLAIN 输出**：

- `type: ALL`（全表扫描）
- `key: NULL`（未使用索引）
- `rows: 表总行数`

---

### 解决方案：创建索引

```sql
CREATE INDEX idx_tb_sku_sn ON tb_sku(sn);
```

**再次查询**：

```sql
SELECT * FROM tb_sku WHERE sn = '10000000314500389';
```

**性能表现**：执行时间降至约 0.001 秒。

**新的 EXPLAIN 输出**：

- `type: ref` 或 `const`（高效索引查找）
- `key: idx_tb_sku_sn`（已使用新索引）
- `rows: 1`（只需读取 1 行）

---

> **总结**：通过索引优化，查询性能可显著提升。欢迎补充更多实战经验和优化技巧！