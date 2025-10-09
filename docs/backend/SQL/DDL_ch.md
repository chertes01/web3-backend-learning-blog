# ⚡ MySQL DDL 操作速查清单

---

## 📌 数据库操作

```sql
SHOW DATABASES;          -- 查看所有数据库
USE testdb;              -- 切换数据库
SELECT DATABASE();       -- 查看当前数据库
```

---

## 📌 表操作

```sql
SHOW TABLES;             -- 查看所有表
DESC 表名;               -- 查看表结构
SHOW TABLES LIKE 'xxx';  -- 模糊匹配表
```

---

## 📌 创建 / 删除表

```sql
CREATE TABLE 表名 (
  字段名 类型 [约束],
  ...
);

DROP TABLE 表名;         -- 删除表
TRUNCATE TABLE 表名;     -- 清空表数据（结构保留）
```

---

## 📌 修改表结构

```sql
ALTER TABLE 表名 ADD 字段 类型;                 -- 新增字段
ALTER TABLE 表名 DROP 字段;                     -- 删除字段
ALTER TABLE 表名 MODIFY 字段 类型;              -- 改字段类型
ALTER TABLE 表名 CHANGE 旧字段 新字段 新类型;   -- 改字段名+类型
ALTER TABLE 表名 RENAME TO 新表名;              -- 改表名
```

---

## 📌 查看修改结果

```sql
DESC 表名;              -- 查看表最新结构
```

---

## ✨ 用法口诀

- 建表 CREATE
- 删表 DROP
- 清空 TRUNCATE
- 改表 ALTER（加 ADD、删 DROP、改 MODIFY/CHANGE、改表名 RENAME）
- 看结果 DESC

---

## 🧭 MySQL DDL 操作流程图

```text
┌───────────────────────────┐
│        创建数据库          │
│  CREATE DATABASE testdb;  │
└──────────────┬────────────┘
               │
               ▼
┌───────────────────────────┐
│      切换到目标数据库      │
│        USE testdb;         │
└──────────────┬────────────┘
               │
               ▼
┌───────────────────────────┐
│         创建数据表         │
│ CREATE TABLE user (        │
│   id INT PRIMARY KEY,      │
│   name VARCHAR(20),        │
│   age TINYINT,             │
│   gender CHAR(1)           │
│ );                         │
└──────────────┬────────────┘
               │
               ▼
┌───────────────────────────┐
│       查看表结构/状态      │
│  SHOW TABLES;              │
│  DESC user;                │
└──────────────┬────────────┘
               │
               ▼
┌───────────────────────────┐
│        修改表结构          │
│ ALTER TABLE user ADD addr VARCHAR(50);     │
│ ALTER TABLE user MODIFY age TINYINT UNSIGNED; │
│ ALTER TABLE user CHANGE name username VARCHAR(30); │
│ ALTER TABLE user RENAME TO userInfo;       │
└──────────────┬────────────┘
               │
               ▼
┌───────────────────────────┐
│        清空表数据          │
│     TRUNCATE TABLE user;   │
└──────────────┬────────────┘
               │
               ▼
┌───────────────────────────┐
│         删除表格           │
│       DROP TABLE user;     │
└──────────────┬────────────┘
               │
               ▼
┌───────────────────────────┐
│        删除数据库          │
│     DROP DATABASE testdb;  │
└───────────────────────────┘
```

---

## 🧩 操作顺序说明

- **CREATE DATABASE / USE**  
  👉 创建并选择目标数据库。

- **CREATE TABLE**  
  👉 定义表结构（字段名、类型、约束）。

- **SHOW / DESC**  
  👉 检查表是否创建成功。

- **ALTER TABLE**  
  👉 修改字段（新增、删除、改名、改类型、改表名）。

- **TRUNCATE TABLE**  
  👉 清空表中数据但保留结构。

- **DROP TABLE / DROP DATABASE**  
  👉 删除表或整个数据库（慎用 ⚠️）。

---