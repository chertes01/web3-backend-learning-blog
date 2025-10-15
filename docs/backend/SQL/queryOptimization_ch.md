# MySQL 查询优化学习笔记

这是一份关于 MySQL 数据库性能优化的学习笔记，涵盖了从数据导入、索引设计到 SQL 查询优化的多个核心主题。

---

## 目录

1. [数据导入 (LOAD DATA INFILE)](#1-数据导入-load-data-infile)
2. [主键优化](#2-主键优化)
3. [ORDER BY 优化](#3-order-by-优化)
4. [GROUP BY 优化](#4-group-by-优化)
5. [LIMIT 分页优化](#5-limit-分页优化)
6. [COUNT 优化](#6-count-优化)
7. [UPDATE 优化与锁机制](#7-update-优化与锁机制)

---

## 1. 数据导入 (LOAD DATA INFILE)

对于大批量数据的插入，使用 `LOAD DATA INFILE` 命令比 `INSERT` 语句快得多，因为它直接读取文件并批量加载，减少了 SQL 解析和网络通信的开销。

### 操作步骤

#### 第一步：检查并开启文件加载权限

```sql
-- 查看当前配置 (0代表关闭, 1代表开启)
SELECT @@local_infile;

-- 开启全局配置 (需要有SUPER权限)
SET GLOBAL local_infile = 1;
```

> **安全警告：** 开启 `local_infile` 存在一定的安全风险，因为它允许服务器访问客户端的文件。请确保只在信任的环境和数据源下使用。

#### 第二步：创建目标表

```sql
CREATE TABLE `tb_user` (
  `id` INT(11) NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(50) NOT NULL,
  `password` VARCHAR(50) NOT NULL,
  `name` VARCHAR(20) NOT NULL,
  `birthday` DATE DEFAULT NULL,
  `sex` CHAR(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_user_username` (`username`)
) ENGINE=INNODB DEFAULT CHARSET=utf8;
```

#### 第三步：执行导入命令

假设数据文件 `tb_user.sql` 存放在 `/home/von/testDataMySQL/` 目录下。

```sql
LOAD DATA LOCAL INFILE '/home/von/testDataMySQL/tb_user.sql'
INTO TABLE tb_user
FIELDS TERMINATED BY ',' -- 指定字段之间的分隔符
LINES TERMINATED BY '\n'; -- 指定行与行之间的分隔符
```

---

## 2. 主键优化

### 2.1 数据组织方式

在 InnoDB 存储引擎中，表本身就是一个索引组织表 (Index-Organized Table)。表中的数据根据主键的顺序进行物理存储，这个主键索引也被称为聚簇索引。

### 2.2 页分裂与页合并

- **页分裂 (Page Split)：**  
  当数据页已满且插入的记录主键值位于页中间时，InnoDB 会将该页拆分成两个新页，导致额外的 I/O 和磁盘碎片。

- **页合并 (Page Merge)：**  
  删除数据后，若页中“已删除”记录比例达到一定阈值（默认为 50%），InnoDB 会尝试合并相邻页以回收空间。

### 2.3 主键设计原则

1. **降低长度：** 主键长度应尽可能短，避免二级索引变得更大。
2. **顺序插入：** 使用单调递增的主键（如 `AUTO_INCREMENT`），避免页分裂。
3. **避免无序主键：** 不建议使用 UUID 或身份证号作为主键。
4. **避免修改主键：** 修改主键的成本非常高，相当于删除旧记录并插入新记录。

---

## 3. ORDER BY 优化

`ORDER BY` 的目标是获取有序的结果集。MySQL 有两种排序方式：

1. **Using index：** 利用索引的有序性直接返回数据，无需额外排序，效率极高。
2. **Using filesort：** 无法利用索引时，MySQL 需要在内存或磁盘上进行排序，成本较高。

### 3.1 优化场景分析

假设创建了联合索引：

```sql
CREATE INDEX idx_tb_user_age_ph ON tb_user(age, phone);
```

| 查询语句                              | EXPLAIN 结果 (Extra) | 分析                                                                 |
|---------------------------------------|----------------------|----------------------------------------------------------------------|
| `ORDER BY age, phone`                 | Using index          | 排序顺序与索引一致，直接通过索引返回有序结果。                       |
| `ORDER BY age DESC, phone DESC`       | Backward index scan; Using index | 排序方向与索引相反，MySQL 可反向扫描索引，效率高。                  |
| `SELECT * FROM tb_user ORDER BY age, phone` | Using filesort       | 查询了 `*`，索引不覆盖所有字段，需回表查询，导致 filesort。         |
| `ORDER BY age ASC, phone DESC`        | Using index; Using filesort | 一个字段升序，一个字段降序，部分字段需 filesort 排序。              |

### 3.2 升降序混合排序优化

MySQL 8.0+ 支持在创建索引时指定排序规则：

```sql
CREATE INDEX idx_tb_user_age_ph_ad ON tb_user(age ASC, phone DESC);

-- 查询时可完美利用索引
EXPLAIN SELECT id, age, phone FROM tb_user ORDER BY age ASC, phone DESC;
```

### 3.3 ORDER BY 优化总结

1. **建立合适索引：** 根据排序字段建立索引，遵循最左前缀法则。
2. **使用覆盖索引：** 查询字段尽量包含在索引中，避免回表导致 filesort。
3. **注意排序方向：** 确保查询排序方向与索引一致，必要时创建匹配索引。
4. **调整排序缓冲区：** 若无法避免 filesort，可适当增大 `sort_buffer_size`。

---

## 4. GROUP BY 优化

`GROUP BY` 的优化思路与 `ORDER BY` 类似，因为它隐含了排序操作。

1. **Using index：** 直接利用索引完成分组，效率最高。
2. **Using temporary：** 创建临时表存储中间结果，效率较低。

### 4.1 优化场景分析

假设创建了索引：

```sql
CREATE INDEX idx_tb_user_pro_age_sta ON tb_user(profession, age, status);
```

| 查询语句                              | EXPLAIN 结果 (Extra) | 分析                                                                 |
|---------------------------------------|----------------------|----------------------------------------------------------------------|
| `GROUP BY profession`                 | Using index          | `profession` 是索引的最左前缀，可直接利用索引分组。                  |
| `GROUP BY age`                        | Using index; Using temporary | 未遵循最左前缀法则，需创建临时表。                                  |
| `WHERE profession='产品经理' GROUP BY age` | Using index          | WHERE 条件锁定了第一个字段，`age` 成为连续字段，可利用索引分组。    |

---

## 5. LIMIT 分页优化

在深度分页场景下，`LIMIT offset, rows` 的性能会急剧下降。例如 `LIMIT 2000000, 10`，MySQL 需要扫描并丢弃前 2,000,000 条记录。

### 优化方案：索引覆盖 + 内连接

```sql
-- 优化前：扫描 2,000,010 条记录
SELECT * FROM tb_sku LIMIT 2000000, 10;

-- 优化后：
SELECT s.*
FROM tb_sku s
INNER JOIN (
    SELECT id FROM tb_sku LIMIT 2000000, 10
) AS tmp ON s.id = tmp.id;
```

---

## 6. COUNT 优化

### 6.1 存储引擎差异

- **MyISAM：** 表总行数存储在磁盘上，`COUNT(*)` 直接返回，速度极快。
- **InnoDB：** 不存储总行数，`COUNT(*)` 需全表扫描或扫描索引，速度较慢。

### 6.2 COUNT 函数的用法比较

| 用法          | 效率排序       | 执行过程 (InnoDB)                                              |
|---------------|----------------|---------------------------------------------------------------|
| `COUNT(*)`    | 最高           | 不读取字段，直接按行累加。                                     |
| `COUNT(1)`    | 约等于 `COUNT(*)` | 遍历表，服务层对每行放入数字“1”，然后累加。                  |
| `COUNT(主键)` | 稍低           | 遍历聚簇索引，取主键值返回服务层，主键不为 NULL，直接累加。    |
| `COUNT(字段)` | 最低           | 遍历表，取字段值，服务层判断是否为 NULL，不为 NULL 才累加。    |

---

## 7. UPDATE 优化与锁机制

`UPDATE` 的性能与并发能力取决于 WHERE 条件是否使用索引，这决定了 MySQL 使用行锁还是表锁。

- **行锁：** 只锁定需要修改的行，并发性能高。
- **表锁：** 锁定整张表，其他会话无法写入，性能较差。

### 场景分析

#### 场景一：WHERE 条件是索引列

```sql
START TRANSACTION;
UPDATE course SET name = 'Java' WHERE id = 1; -- 使用行锁
COMMIT;
```

#### 场景二：WHERE 条件是非索引列

```sql
START TRANSACTION;
UPDATE course SET name = 'Python' WHERE course_name = 'Python程序设计'; -- 使用表锁
COMMIT;
```

**优化：** 为 `course_name` 字段创建索引：

```sql
CREATE INDEX idx_course_name ON course(course_name);
```

---