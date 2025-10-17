# MySQL 锁机制学习笔记

锁是数据库并发控制的基石，用于协调多个进程或线程并发访问资源，保证数据一致性和有效性。锁冲突也是影响数据库并发性能的重要因素。本文系统梳理 MySQL（InnoDB 存储引擎）中的锁机制，包括锁的分类、全局锁、表级锁、行级锁及其实现细节。

---

## 1. 锁的分类（按粒度）

- **全局锁（Global Lock）**：锁定整个数据库实例，影响范围最大。
- **表级锁（Table-Level Lock）**：每次操作锁住整张表，粒度较大。
- **行级锁（Row-Level Lock）**：每次操作锁住对应的行数据，粒度最细。

---

## 2. 全局锁（Global Lock）

全局锁会让整个实例进入只读状态，所有 DML、DDL 及事务提交都会被阻塞。

**典型场景**：全库逻辑备份（如 mysqldump），保证备份数据一致性。

**SQL 语法**：

```sql
-- 添加全局锁
FLUSH TABLES WITH READ LOCK;

-- 释放全局锁
UNLOCK TABLES;
```

**缺点与替代方案**：

- 阻塞主库：如果在主库上备份，那么在备份期间所有更新操作都会被阻塞，业务基本上处于停摆状态。
- 主从延迟：如果在从库上备份，那么在备份期间从库不能执行主库同步过来的二进制日志 (binlog)，会导致主从延迟。
- 推荐使用 InnoDB 的 `--single-transaction` 参数，利用 MVCC 快照实现一致性备份，无需加锁。

---

## 3. 表级锁（Table-Level Lock）

表级锁粒度大，冲突概率高，并发度低。

### 3.1 表锁（LOCK TABLES）

- **表共享读锁（READ）**：阻塞写，不阻塞读。
- **表独占写锁（WRITE）**：阻塞读和写。

**示例：表读锁**

```sql
-- 客户端A
LOCK TABLES tb_user READ;
SELECT * FROM tb_user; -- 可查询
UPDATE tb_user SET phone='111111' WHERE id=1; -- 报错，不可更新

-- 客户端B
SELECT * FROM tb_user; -- 可查询
UPDATE tb_user SET phone='222222' WHERE id=1; -- 阻塞

-- 客户端A释放锁后，B的更新才会继续
UNLOCK TABLES;
```

**示例：表写锁**

```sql
-- 客户端A
LOCK TABLES tb_user WRITE;
SELECT * FROM tb_user; -- 可查询
UPDATE tb_user SET phone='333333' WHERE id=1; -- 可更新

-- 客户端B
SELECT * FROM tb_user; -- 阻塞
UPDATE tb_user SET phone='444444' WHERE id=1; -- 阻塞

-- 客户端A释放锁后，B的操作才会继续
UNLOCK TABLES;
```

### 3.2 元数据锁（Meta Data Lock, MDL）

- MDL 锁自动加上，维护表元数据一致性，防止有活动事务时变更表结构。
- DML 操作（SELECT, INSERT ...）：MDL 读锁（共享）
- DDL 操作（ALTER TABLE ...）：MDL 写锁（排他）

**冲突场景**：

- 长事务持有 MDL 读锁，DDL 操作被阻塞，后续所有访问也被阻塞。

### 3.3 意向锁（Intention Lock）

意向锁是表级锁，用于协调行锁和表锁关系，提升性能。

- **意向共享锁（IS Lock）**：如 `SELECT ... LOCK IN SHARE MODE`
- **意向独占锁（IX Lock）**：如 `INSERT, UPDATE, DELETE, SELECT ... FOR UPDATE`

**锁兼容性矩阵**：

| 已持有/请求 | IS | IX | S（行锁） | X（行锁） |
|-------------|----|----|----------|----------|
| IS          | ✅ | ✅ | ✅       | ❌       |
| IX          | ✅ | ✅ | ❌       | ❌       |
| S           | ✅ | ❌ | ✅       | ❌       |
| X           | ❌ | ❌ | ❌       | ❌       |

- 意向锁之间完全兼容，冲突主要在行级别。
- 意向锁与表级读/写锁的冲突，是其核心价值所在。

---

## 4. 行级锁（Row-Level Lock）

InnoDB 的行级锁优势在于粒度细、冲突概率低、并发度高。

**前提**：InnoDB 行锁通过索引项加锁实现，未使用索引检索时会退化为表锁。

### 4.1 行锁类型

- **共享锁（S Lock）**：读锁，允许读，阻止写。
- **排他锁（X Lock）**：写锁，允许写，阻止读和写。

**SQL 与锁类型关系**：

| SQL 操作                       | 行锁类型      | 是否自动加锁 |
|--------------------------------|--------------|--------------|
| INSERT ...                     | 排他锁 (X)   | 自动         |
| UPDATE ...                     | 排他锁 (X)   | 自动         |
| DELETE ...                     | 排他锁 (X)   | 自动         |
| SELECT ...（普通查询）         | 无           | -            |
| SELECT ... LOCK IN SHARE MODE  | 共享锁 (S)   | 手动         |
| SELECT ... FOR UPDATE          | 排他锁 (X)   | 手动         |

### 4.2 快照读 vs 当前读

- **快照读（Snapshot Read）**：普通 SELECT，基于 MVCC，不加锁，读取事务开始时的数据快照。
- **当前读（Current Read）**：如 `SELECT ... FOR UPDATE`、`SELECT ... LOCK IN SHARE MODE`，读取最新已提交版本并加锁。所有 INSERT/UPDATE/DELETE 也是当前读。

### 4.3 行锁实现算法

在 RR（Repeatable Read）隔离级别下，主要有三种锁：

- **行锁（Record Lock）**：锁定单个索引记录。
    ```sql
    -- 对 id=5 这一条索引记录加 X 锁
    UPDATE tb_user SET phone = '555555' WHERE id = 5;
    ```
- **间隙锁（Gap Lock）**：锁定索引记录之间的间隙，不包含记录本身，防止间隙插入。
    ```sql
    -- 假设表中有 id=10 和 id=20，无 id=15
    -- 会对 (10, 20) 区间加间隙锁，禁止插入新记录
    SELECT * FROM tb_user WHERE id = 15 FOR UPDATE;
    ```
- **临键锁（Next-Key Lock）**：锁定索引记录及其前间隙，是 Record Lock + Gap Lock 的组合。
    ```sql
    -- 表中有 id=10, 20, 30
    -- 锁住 (10, 20] 区间，即锁定 id=20 及 10-20 间隙
    SELECT * FROM tb_user WHERE id = 20 FOR UPDATE;
    ```

### 4.4 行锁核心三要素

- **行锁算法**：Record Lock（行锁）、Gap Lock（间隙锁）、Next-Key Lock（临键锁）。
- **加锁时机**：根据索引类型（主键、唯一、普通）和查询条件决定。
- **加锁规则**：InnoDB 的行锁是通过对索引上的索引项加锁来实现的。

### 4.5 行锁核心原则

- 锁只加在用到的索引上：一个查询只会选择一个最优索引来执行，锁也只会施加在这个被选中的索引上，与其他无关的索引没有任何关系。
- 无索引，锁全表：如果不通过索引条件检索数据，InnoDB 会锁定表中的所有记录，行锁会升级为表锁。

**锁的退化与升级**：

- 唯一索引等值查询：记录存在时，Next-Key Lock 退化为 Record Lock；不存在时，退化为 Gap Lock。
- 普通索引等值查询：Next-Key Lock 可能退化为 Gap Lock。
- 无索引查询：行锁升级为表锁，锁定整张表所有记录，应避免。

---

## 5. 行锁核心原则与场景解析

### 5.1 行锁核心原则

- 锁只加在用到的索引上：一个查询只会选择一个最优索引来执行，锁也只会施加在这个被选中的索引上，与其他无关的索引没有任何关系。
- 无索引，锁全表：如果不通过索引条件检索数据，InnoDB 会锁定表中的所有记录，行锁会升级为表锁。

### 5.2 行锁场景深度解析（RR 隔离级别）

### 场景一：基于主键/唯一索引的锁定

```sql
-- SQL: SELECT * FROM tb_user WHERE id = 17 LOCK IN SHARE MODE;
```
- 加锁步骤：
  1. 在 tb_user 表上加意向共享锁（IS Lock）。
  2. 在主键索引上加行级锁。
- 具体行锁类型：
  - 情况A：id=17 记录存在，锁类型为行锁（Record Lock），仅锁定 id=17 这一行，不影响 id=16 或 id=18 的插入。
  - 情况B：id=17 记录不存在（假设只存在10和20），锁类型为间隙锁（Gap Lock），锁定 (10, 20) 区间，防止任何记录插入到这个间隙中。

### 场景二：基于非主键索引（二级索引）的锁定

```sql
-- SQL: SELECT * FROM tb_user WHERE age = 30 FOR UPDATE; -- 假设 age 是普通索引
```
- 这是一个“两步走”的加锁过程：
  1. 锁定二级索引（idx_age），锁类型为临键锁（Next-Key Lock），锁定 age=30 的索引项以及相关间隙，防止其他事务插入新的 age=30 的记录，避免幻读。
  2. 根据在 idx_age 上找到的记录的主键值（比如 id=20 和 id=30），回到主键索引中加行锁（Record Lock），防止这几行数据被其他事务通过主键进行修改或删除，保证数据一致性。

**总结**：对二级索引的锁定，是一个先通过二级索引锁定一个“范围”以防幻读，再通过主键索引锁定具体“目标”以保护数据的严谨过程。
---

## 6. 总结

- MySQL 锁机制分为全局锁、表级锁、行级锁，粒度依次减小，并发能力依次增强。
- 行级锁依赖索引实现，合理设计索引和 SQL 能有效提升并发性能，避免锁升级。
- 了解各种锁的实现和冲突场景，是高并发数据库系统设计和优化的基础。

---