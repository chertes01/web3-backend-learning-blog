# MySQL 索引优化核心笔记 🚀

本文档是一份全面的 MySQL 索引学习指南，从核心原理到性能优化，再到避坑技巧，旨在帮助你彻底掌握 MySQL 索引，写出高性能的 SQL 查询。

---

## 目录

- [📖 核心原理篇](#-核心原理篇)
  - [联合索引与最左前缀法则](#联合索引与最左前缀法则)
  - [索引底层结构与“回表”](#索引底层结构与回表)
- [✅ 性能优化篇](#-性能优化篇)
  - [覆盖索引 (Covering Index)](#覆盖索引-covering-index)
  - [索引与排序 (ORDER BY) 优化](#索引与排序-order-by-优化)
  - [函数/表达式索引 (MySQL 8.0+)](#函数表达式索引-mysql-80)
  - [前缀索引 (Prefix Index)](#前缀索引-prefix-index)
  - [SQL 提示 (SQL Hints)](#sql-提示-sql-hints)
- [❌ 避坑指南篇](#-避坑指南篇)
  - [常见索引失效场景](#常见索引失效场景)
  - [优化器选择与数据分布](#优化器选择与数据分布)
- [📜 终极总结篇](#-终极总结篇)
  - [索引设计原则](#索引设计原则)

---

## 📖 核心原理篇

### 联合索引与最左前缀法则

**法则定义：**  
查询条件必须从索引的最左边的列开始，并且不能跳过中间的列。如果跳过某一列，那么该列右侧的所有列都无法使用索引进行优化。

假设联合索引 `idx_user_pro_age_sta` 的顺序为 `(profession, age, status)`：

1. **完美匹配 - 索引完全生效**

```sql
EXPLAIN SELECT * FROM tb_user WHERE profession='产品经理' AND age=29 AND status='1';
EXPLAIN SELECT * FROM tb_user WHERE status='1' AND age=29 AND profession='产品经理';
```

**执行计划亮点：**
- `type: ref`（非唯一性索引查找）
- `key: idx_user_pro_age_sta`（实际使用的索引）
- `key_len: 47`（索引的所有部分都被用上）

2. **跳过中间列 - 索引部分失效**

```sql
EXPLAIN SELECT * FROM tb_user WHERE profession='产品经理' AND status='1';
```

**执行计划亮点：**
- `key_len: 43`（只有 `profession` 部分被用到，`status` 因 `age` 被跳过而失效）

3. **从中间列开始 - 索引完全失效**

```sql
EXPLAIN SELECT * FROM tb_user WHERE age=29 AND status='1';
```

**执行计划亮点：**
- `type: ALL`（全表扫描，索引完全失效）

---

### 索引底层结构与“回表”

- **聚集索引 (Clustered Index)：** 表数据按主键顺序存放在 B+Tree 结构中，一个表只能有一个聚集索引。
- **二级索引 (Secondary Index)：** 普通索引的叶子节点存储的是主键值，而不是数据行的物理地址。

**回表定义：**  
通过二级索引找到主键值，再通过主键值去聚集索引查找完整数据的过程，称为“回表”（Back to Table）。

---

## ✅ 性能优化篇

### 覆盖索引 (Covering Index)

**定义：**  
查询中需要的字段全部包含在索引中，无需回表即可返回数据。

- **未使用覆盖索引（需要回表）：**

```sql
EXPLAIN SELECT id, profession, age, status, email FROM tb_user WHERE profession='产品经理';
```

- **使用覆盖索引（无需回表）：**

```sql
EXPLAIN SELECT id, profession, age, status FROM tb_user WHERE profession='产品经理';
```

**Extra 字段：**  
`Using index` 表示成功使用覆盖索引。

---

### 索引与排序 (ORDER BY) 优化

**优化条件：**  
`ORDER BY` 的字段顺序和排序方向必须与索引保持一致，并且查询条件遵循最左前缀法则。

```sql
-- ✅ 避免了文件排序
EXPLAIN SELECT id, age, status FROM tb_user 
WHERE profession = '产品经理' ORDER BY age ASC, status ASC;

-- ❌ 触发文件排序
EXPLAIN SELECT id, age, status FROM tb_user 
WHERE profession = '产品经理' ORDER BY status ASC, age ASC;
```

---

### 函数/表达式索引 (MySQL 8.0+)

**问题：**  
`WHERE YEAR(create_time) = 2025` 无法使用 `create_time` 上的索引。

**解决方案：**  
直接在函数或表达式上创建索引。

```sql
CREATE INDEX idx_createtime_year ON tb_user((YEAR(create_time)));
EXPLAIN SELECT id FROM tb_user WHERE YEAR(create_time) = 2025;
```

---

### 前缀索引 (Prefix Index)

**适用场景：**  
长字符串列可以使用前缀索引，只索引字符串的一部分，以节省空间、提高效率。

```sql
-- 计算不同前缀长度的选择性
SELECT 
    COUNT(DISTINCT email) / COUNT(*) AS sel_full,
    COUNT(DISTINCT LEFT(email, 8)) / COUNT(*) AS sel_prefix_8
FROM tb_user;

-- 创建前缀索引
CREATE INDEX idx_email_prefix ON tb_user(email(10));
```

---

### SQL 提示 (SQL Hints)

- `USE INDEX`：建议使用指定索引。
- `IGNORE INDEX`：忽略指定索引。
- `FORCE INDEX`：强制使用指定索引。

```sql
EXPLAIN SELECT * FROM tb_user USE INDEX(idx_tb_user_pro) WHERE profession='产品经理';
EXPLAIN SELECT * FROM tb_user FORCE INDEX(idx_user_pro_age_sta) WHERE profession='产品经理';
```

---

## ❌ 避坑指南篇

### 常见索引失效场景

| 失效场景           | SQL 示例                                         | 原因分析                                                         |
|--------------------|--------------------------------------------------|------------------------------------------------------------------|
| 范围查询右侧失效   | WHERE profession='产品经理' AND age > 29 AND status='1' | 范围查询（>, <, BETWEEN）之后的列会失效                          |
| 索引列上运算       | WHERE substring(phone, 10, 2) = '13'             | 对索引列使用函数或计算，导致索引失效                              |
| 隐式类型转换       | WHERE phone = 13998877665                        | `phone` 为 `varchar`，查询用数字，发生隐式转换                   |
| 头模糊匹配         | WHERE profession LIKE '%产品'                    | LIKE 以 `%` 开头，不符合最左前缀法则                              |
| OR 连接非索引列    | WHERE id=15 OR age=28 (age 无索引)               | OR 条件中只要有一个列无索引，整个查询放弃索引                    |

---

### 优化器选择与数据分布

**场景 1：索引选择性差**

```sql
EXPLAIN SELECT * FROM tb_user WHERE phone >= '13677778888'; -- 全表扫描
EXPLAIN SELECT * FROM tb_user WHERE phone >= '13812345678'; -- 使用索引
```

**场景 2：IS NOT NULL**

```sql
EXPLAIN SELECT * FROM tb_user WHERE profession IS NOT NULL; -- 全表扫描
```

---

## 📜 终极总结篇

### 索引设计原则

1. **为高频查询的大表建立索引：** 小表或不常查询的表无需索引。
2. **为常用于 WHERE、ORDER BY、GROUP BY 的列建立索引：** 索引最能发挥作用的场景。
3. **选择区分度高的列：** 性别这种低基数列不适合做索引，而用户 ID、手机号等高基数列是理想的索引列。
4. **尽量使用联合索引：** 合理设计联合索引可以实现覆盖索引，避免回表。
5. **控制索引数量：** 索引并非越多越好，过多索引会增加写操作的开销。
6. **使用前缀索引：** 对于长字符串列，使用前缀索引可以节省空间。
7. **NOT NULL 约束：** 如果索引列不可能为 NULL，使用 NOT NULL 约束。

---