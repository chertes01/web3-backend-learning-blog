# SQL View (视图) 学习笔记

视图（View）是一种虚拟存在的表，它本质上是一个预先定义好的 `SELECT` 查询，作为一个命名的对象存储在数据库中。本文详细介绍了视图的基本操作、核心概念以及高级用法。

---

## 目录

1. [什么是视图？](#1-什么是视图)
2. [视图的基本操作](#2-视图的基本操作)
   - [创建视图](#21-创建视图)
   - [查询视图](#22-查询视图)
   - [修改视图](#23-修改视图)
   - [删除视图](#24-删除视图)
3. [核心概念：视图检查选项 (WITH CHECK OPTION)](#3-核心概念视图检查选项-with-check-option)
   - [CASCADED (级联检查)](#31-cascaded-级联检查)
   - [LOCAL (本地检查)](#32-local-本地检查)
   - [CASCADED vs LOCAL 对比总结](#33-cascaded-vs-local-对比总结)
4. [核心概念：视图的更新性](#4-核心概念视图的更新性)
5. [进阶知识：视图算法 (ALGORITHM)](#5-进阶知识视图算法-algorithm)
6. [视图的作用与优势](#6-视图的作用与优势)
7. [实战练习](#7-实战练习)

---

## 1. 什么是视图？

视图是一种虚拟表，它的核心特点包括：

- **虚拟性：** 视图本身不包含任何数据，数据来自定义视图时的基表（Base Table），并在每次查询时动态生成。
- **逻辑封装：** 视图只保存查询的 SQL 逻辑，不保存查询结果。
- **操作简化：** 视图可以像普通表一样查询，有时也可以进行增删改操作。

---

## 2. 视图的基本操作

### 2.1 创建视图

使用 `CREATE VIEW` 语句创建视图，`OR REPLACE` 子句表示如果已存在同名视图，则替换它。

```sql
-- 语法
CREATE OR REPLACE VIEW view_name AS
SELECT column1, column2, ...
FROM table_name
WHERE [condition];

-- 示例：创建一个只包含 id <= 10 的学生的视图
CREATE OR REPLACE VIEW student_v_1 AS
SELECT id, name
FROM student
WHERE id <= 10;
```

---

### 2.2 查询视图

查询视图的方法与查询普通表完全一致。

```sql
-- 查看视图的创建语句
SHOW CREATE VIEW student_v_1;

-- 查询视图的全部数据
SELECT * FROM student_v_1;

-- 对视图进行条件查询
SELECT * FROM student_v_1 WHERE id < 3;
```

---

### 2.3 修改视图

修改视图通常指修改视图底层的 `SELECT` 查询逻辑。

```sql
-- 方法一：使用 CREATE OR REPLACE VIEW 重新定义
CREATE OR REPLACE VIEW student_v_1 AS
SELECT id, name, gender -- 添加 gender 字段
FROM student
WHERE id <= 10;

-- 方法二：使用 ALTER VIEW
ALTER VIEW student_v_1 AS
SELECT id, name -- 移除 gender 字段
FROM student
WHERE id <= 10;
```

---

### 2.4 删除视图

```sql
DROP VIEW IF EXISTS student_v_1;
```

---

## 3. 核心概念：视图检查选项 (WITH CHECK OPTION)

`WITH CHECK OPTION` 确保通过视图进行的 `INSERT` 或 `UPDATE` 操作必须满足视图的 `WHERE` 条件。

### 3.1 CASCADED (级联检查)

`CASCADED` 会递归检查本视图及所有基于的底层视图的 `WHERE` 条件。

```sql
-- 创建带 CASCADED 检查选项的视图
CREATE OR REPLACE VIEW stu_v_1 AS
SELECT id, name FROM student
WHERE id <= 10
WITH CASCADED CHECK OPTION;

-- 成功：满足条件
INSERT INTO stu_v_1 VALUES(7, 'Tom');

-- 失败：不满足条件
INSERT INTO stu_v_1 VALUES(30, 'Tommy');
```

---

### 3.2 LOCAL (本地检查)

`LOCAL` 只检查本视图的 `WHERE` 条件，底层视图的条件只有在定义了 `CHECK OPTION` 时才会检查。

```sql
-- 创建带 LOCAL 检查选项的视图
CREATE OR REPLACE VIEW stu_v_4 AS
SELECT id, name FROM stu_v_1 WHERE id >= 5
WITH LOCAL CHECK OPTION;

-- 成功：只检查 stu_v_4 的条件
INSERT INTO stu_v_4 VALUES(12, 'Hany');
```

---

### 3.3 CASCADED vs LOCAL 对比总结

| 检查选项 | 是否检查本视图的 WHERE 条件？ | 如何检查底层视图的 WHERE 条件？ |
|----------|-----------------------------|--------------------------------|
| CASCADED | 是                          | 总是检查，无论底层视图是否定义了 CHECK OPTION |
| LOCAL    | 是                          | 仅当底层视图定义了 CHECK OPTION 时才检查 |

---

## 4. 核心概念：视图的更新性

视图中的行与基表中的行必须存在一对一的关系，视图才可更新。如果视图包含以下内容，则不可更新：

- 聚合函数（如 `SUM()`、`COUNT()` 等）
- `DISTINCT`
- `GROUP BY`
- `HAVING`
- `UNION` 或 `UNION ALL`
- 子查询
- 不可更新的视图
- 连接查询（部分内连接可更新）

```sql
-- 示例：包含聚合函数的视图不可更新
CREATE VIEW stu_v_count AS SELECT COUNT(*) FROM student;

-- 插入操作会失败
INSERT INTO stu_v_count VALUES(10);
```

---

## 5. 进阶知识：视图算法 (ALGORITHM)

创建视图时，可以指定 MySQL 处理视图的算法：

- **MERGE：** 将视图的查询语句与定义语句合并，直接查询基表。效率最高，且可更新。
- **TEMPTABLE：** 先创建临时表，再查询临时表。效率较低。
- **UNDEFINED：** 默认值，MySQL 自动选择算法。

```sql
CREATE ALGORITHM = MERGE VIEW view_name AS
SELECT ...;
```

---

## 6. 视图的作用与优势

1. **简化操作：**  
   视图隐藏了复杂的查询逻辑，用户只需查询视图即可获取结果。

2. **增强安全性：**  
   - **列级别安全：** 只暴露部分字段。
   - **行级别安全：** 只暴露符合条件的数据。

3. **数据独立性：**  
   即使基表结构发生变化，只要视图定义保持不变，用户查询无需修改。

---

## 7. 实战练习

### 练习 1

**需求：** 屏蔽用户表中的手机号和邮箱字段，仅暴露基本信息。

```sql
CREATE VIEW tb_user_v AS
SELECT id, name, profession, age, gender, status, createtime
FROM tb_user;
```

---

### 练习 2

**需求：** 查询每个学生所选修的课程（多表联查）。

```sql
CREATE VIEW stu_course_v AS
SELECT
    s.name AS student_name,
    c.name AS course_name
FROM
    student s
LEFT JOIN
    student_course sc ON s.id = sc.student_id
LEFT JOIN
    course c ON sc.course_id = c.id;

-- 使用视图
SELECT * FROM stu_course_v WHERE student_name = '张三';
```

---