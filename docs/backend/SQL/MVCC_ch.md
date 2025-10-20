# MySQL MVCC 学习笔记

## MVCC (Multi-Version Concurrency Control)

### 1 基本概念

#### 1. 当前读 (Current Read)
- 读取的是记录的最新版本。
- 读取时需要加锁，以保证其他并发事务不能修改当前记录。
- 日常操作中的以下语句都是当前读：
  - `select... lock in share mode` (共享锁)
  - `select... for update` (排他锁)
  - `update`、`insert`、`delete` (排他锁)
- 即使在默认的 Repeatable Read (RR) 隔离级别下，如果查询语句后面加上了共享锁或排他锁，也会执行当前读操作，读取到事务最新提交的内容。

#### 2. 快照读 (Snapshot Read)
- 简单的 `select` (不加锁) 就是快照读。
- 快照读读取的是记录数据的可见版本，可能是历史数据。
- 快照读不加锁，是非阻塞读。
- **不同隔离级别下快照读的生成时机：**
  - **Read Committed (RC):** 每次 `select` 都会生成一个快照读。
  - **Repeatable Read (RR):** 开启事务后第一个 `select` 语句才是生成快照读的地方，后续复用该快照。
  - **Serializable:** 快照读会退化为当前读。

#### 3. MVCC 简介
- **全称:** Multi-Version Concurrency Control，多版本并发控制。
- 指维护一个数据的多个版本，使得读写操作没有冲突。
- 快照读为 MVCC 提供了非阻塞读功能。
- **MVCC 的具体实现依赖于以下三个要素:**
  - 数据库记录中的三个隐式字段
  - undo log 日志（版本链）
  - ReadView

### 2 隐藏字段

InnoDB 表中除了用户定义的字段外，还会自动添加以下隐藏字段：

| 隐藏字段      | 含义                 | 备注                                           |
| ------------- | -------------------- | ---------------------------------------------- |
| `DB_TRX_ID`   | 最近修改事务 ID      | 记录插入或最后一次修改该记录的事务 ID。        |
| `DB_ROLL_PTR` | 回滚指针             | 指向这条记录的上一个版本，用于配合 undo log。    |
| `DB_ROW_ID`   | 隐藏主键             | 只有在表结构没有指定主键时，才会生成该隐藏字段。 |

**注意:** `DB_TRX_ID` 和 `DB_ROLL_PTR` 是肯定会添加的。

### 3 undo log (回滚日志)

- **介绍:** 用于记录数据被修改前的信息。
- **作用:**
  - 提供回滚（保证事务的原子性）。
  - 实现 MVCC（多版本并发控制）。
- `undo log` 是逻辑日志，与 `redo log` 记录物理日志不同。例如，`delete` 记录时，`undo log` 记录一条对应的 `insert` 记录。
- 当执行 `rollback` 时，从 `undo log` 中读取逻辑记录并进行回滚。

#### undo log 的删除
- 事务提交时，并不会立即删除 `undo log`，因为这些日志可能还用于 MVCC。
- `insert` 产生的 `undo log` 只在回滚时需要，事务提交后可被立即删除。
- `update`、`delete` 产生的 `undo log` 不仅在回滚时需要，在快照读时也需要，因此不会立即被删除。

#### 版本链 (Version Chain)
- 不同事务或相同事务对同一条记录进行修改，会导致该记录的 `undo log` 生成一条记录版本链表。
- 版本链的头部是最新的旧记录，链表尾部是最早的旧记录。
- 记录中的 `DB_ROLL_PTR`（回滚指针）指向这条记录的上一个版本，从而将这些不同版本的旧数据串联成链表。

### 4 ReadView (读视图)

- ReadView 是快照读 SQL 执行时 MVCC 提取数据的依据。
- 它记录并维护系统当前活跃的事务 (未提交的) ID。

#### 核心字段

| 字段             | 含义                                                     |
| ---------------- | -------------------------------------------------------- |
| `m_ids`          | 当前活跃的事务 ID 集合。                                 |
| `min_trx_id`     | 最小活跃事务 ID。                                        |
| `max_trx_id`     | 预分配事务 ID，等于当前最大事务 ID + 1 (事务 ID 是自增的)。 |
| `creator_trx_id` | ReadView 创建者的事务 ID。                               |

#### 版本链数据访问规则

ReadView 规定了快照读访问 `undo log` 版本链数据的规则（`trx_id` 代表 `undo log` 版本链对应事务 ID）：

| 条件                               | 是否可以访问 | 说明                                     |
| ---------------------------------- | ------------ | ---------------------------------------- |
| `trx_id == creator_trx_id`         | 可以         | 该版本数据是当前这个事务更改的。         |
| `trx_id < min_trx_id`              | 可以         | 该版本数据已经提交了。                   |
| `trx_id > max_trx_id`              | 不可以       | 该事务是在 ReadView 生成后才开启的。     |
| `min_trx_id <= trx_id <= max_trx_id` | 如果 `trx_id` 不在 `m_ids` 中，可以访问 | 该版本数据已经提交。                     |

#### ReadView 的生成时机
- **READ COMMITTED (RC):** 在事务中每一次执行快照读时生成 ReadView。
- **REPEATABLE READ (RR):** 仅在事务中第一次执行快照读时生成 ReadView，后续复用该 ReadView。

### 5 原理总结

- **MVCC 的实现原理:** 通过 InnoDB 表的隐藏字段、UndoLog 版本链、ReadView 来实现。
- **实现事务的隔离性:** MVCC 加上锁实现了事务的隔离性。
- **保证一致性:** Redo log 与 undo log 保证了事务的一致性。