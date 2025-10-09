# MySQL DQL 学习笔记

---

## DQL（Data Query Language）简介

- DQL 主要用于数据查询，核心语句是 `SELECT`。
- 常用子句：`WHERE`、`GROUP BY`、`HAVING`、`ORDER BY`、`LIMIT`。

---

## 1️⃣ 基本查询

**查询指定字段**
```sql
SELECT Name, Age, Workplace FROM employee;
```
- 只返回指定列，建议明确字段，便于维护。

**查询所有字段**
```sql
SELECT * FROM employee;
SELECT ID, Name, Age, Gender, idCare, Entrydate, Workplace FROM employee;
```
- `*` 虽然方便，但性能和安全性较差。

---

## 2️⃣ 去重查询

**字段别名**
```sql
SELECT Workplace AS '工作地点' FROM employee;
```
- `AS` 可为结果集字段起别名。

**去重**
```sql
SELECT DISTINCT Workplace FROM employee;
```
- `DISTINCT` 去除重复值。

---

## 3️⃣ 条件查询（WHERE 子句）

**等值条件**
```sql
SELECT * FROM employee WHERE Age = 28;
```

**比较运算**
```sql
SELECT * FROM employee WHERE Age <= 30;
```

**逻辑运算**
```sql
SELECT * FROM employee WHERE Age <= 30 AND Gender = 'male';
SELECT * FROM employee WHERE Age <= 30 && Gender = 'male';
```
- `AND` 与 `&&` 效果一样。

**范围查询**
```sql
SELECT * FROM employee WHERE Age BETWEEN 30 AND 40;
```
- 包含边界，等价于 `Age >= 30 AND Age <= 40`。

**空值判断**
```sql
SELECT * FROM employee WHERE idCare IS NULL;
```
- 不能用 `= NULL`。

**多值匹配**
```sql
SELECT * FROM employee WHERE Age IN (31, 28, 29);
```
- `IN` 比多个 `OR` 简洁。

**模糊查询**
```sql
SELECT * FROM employee WHERE idCare LIKE '%8999';
```
- `%` 表示任意长度字符串。

---

## 4️⃣ 聚合函数（统计）

**基本聚合**
```sql
SELECT COUNT(*) AS '员工数量' FROM employee;
SELECT AVG(Age) AS '平均年龄' FROM employee;
SELECT MIN(Age) AS '年龄最小的员工' FROM employee;
SELECT MAX(Age) AS '最年长员工' FROM employee;
SELECT SUM(Age) FROM employee WHERE Workplace = 'Shenzhen';
```
- `AVG`、`SUM` 会自动忽略 `NULL`。

---

## 5️⃣ 分组查询（GROUP BY + HAVING）

**分组统计**
```sql
SELECT Gender, COUNT(*) FROM employee GROUP BY Gender;
SELECT Gender, COUNT(*) AS '人数', AVG(Age) AS '平均年龄' FROM employee GROUP BY Gender;
```

**分组条件（HAVING）**
```sql
SELECT Workplace, COUNT(*) AS '人数'
FROM employee
WHERE Age < 40
GROUP BY Workplace
HAVING COUNT(*) > 2;
```
- `WHERE` 过滤行，`HAVING` 过滤分组。

**执行顺序**  
FROM → WHERE → GROUP BY → HAVING → SELECT → ORDER BY → LIMIT

---

## 6️⃣ 排序查询（ORDER BY）

**基本排序**
```sql
SELECT * FROM employee ORDER BY Age ASC;
SELECT * FROM employee ORDER BY Entrydate DESC;
```

**多条件排序**
```sql
SELECT * FROM employee ORDER BY Age ASC, Entrydate DESC;
```
- 先按年龄升序，相同年龄再按入职时间降序。

---

## 7️⃣ 分页查询（LIMIT）

**第一页，每页10条**
```sql
SELECT * FROM employee LIMIT 0, 10;
```

**第二页，每页10条**
```sql
SELECT * FROM employee LIMIT 10, 10;
```
- 语法：`LIMIT 偏移量, 条数`
- 第一页 `LIMIT 0,10`；第二页 `LIMIT 10,10`；第三页 `LIMIT 20,10`…

---

## 8️⃣ 综合练习

**多条件查询**
```sql
SELECT * FROM employee WHERE Age IN (24, 30, 28) AND Gender = 'female';
```

**多条件 + 字符长度**
```sql
SELECT * FROM employee 
WHERE Age <= 30 AND Gender = 'male' AND CHAR_LENGTH(Name) <= 5;
```

**分组统计**
```sql
SELECT Gender, COUNT(*) AS '人数'
FROM employee
WHERE Age <= 30
GROUP BY Gender;
```

**排序 + 限制**
```sql
SELECT ID, Name, Age, Entrydate
FROM employee
WHERE Age <= 35
ORDER BY Age ASC, Entrydate ASC;
```

**男性，年龄25–40之间，取前5条**
```sql
SELECT * FROM employee
WHERE Age BETWEEN 25 AND 40 AND Gender = 'male'
ORDER BY Age ASC, Entrydate ASC
LIMIT 0, 5;
```

---

## 9️⃣ 易错点总结

- `IS NULL` / `IS NOT NULL` 不能写成 `= NULL`
- `BETWEEN` 包含边界，顺序不能写反
- `WHERE` 行级过滤，`HAVING` 组级过滤
- `LIMIT` 第二个参数是条数，不是结束位置
- `LIMIT 10,10` 表示从第11条开始，取10条
- 多条件建议用 `AND`、`OR`，`&&`、`||` 兼容性差
- `ORDER BY` 注意排序字段是否正确

---

## DQL 语法对照表

| 知识点         | 语法                                      | 示例                                                         | 注意点 / 易错点                   |
|----------------|-------------------------------------------|--------------------------------------------------------------|-----------------------------------|
| 基本查询       | SELECT 字段列表 FROM 表名;                | SELECT Name, Age, Workplace FROM employee;                   | 建议写明字段而不是用 *            |
| 查询所有字段   | SELECT * FROM 表;                         | SELECT * FROM employee;                                      | * 会返回所有列，性能略差           |
| 别名           | SELECT 字段 AS 别名                       | SELECT Workplace AS '工作地点' FROM employee;                | AS 可省略，但加上更清晰           |
| 去重           | SELECT DISTINCT 字段 FROM 表;             | SELECT DISTINCT Workplace FROM employee;                     | DISTINCT 去掉重复值                |
| 条件查询       | WHERE 条件                                | SELECT * FROM employee WHERE Age = 28;                       | 执行顺序：FROM → WHERE            |
| 比较运算       | =, <, >, <=, >=, <>                       | SELECT * FROM employee WHERE Age <= 30;                      | <> 表示不等于                     |
| 逻辑运算       | AND / OR / NOT                            | SELECT * FROM employee WHERE Age <= 30 AND Gender = 'male';  | MySQL 中 && / \|\| 也可用          |
| 范围查询       | BETWEEN a AND b                           | SELECT * FROM employee WHERE Age BETWEEN 30 AND 40;          | 含边界，等价于 >=30 AND <=40      |
| 多值匹配       | IN (值1,值2,...)                          | SELECT * FROM employee WHERE Age IN (31,28,29);              | 更简洁，避免多个 OR               |
| 模糊查询       | LIKE '匹配模式'                           | SELECT * FROM employee WHERE idCare LIKE '%8999';            | % 任意长度，_ 单个字符            |
| 空值判断       | IS NULL / IS NOT NULL                     | SELECT * FROM employee WHERE idCare IS NULL;                 | 不能写 = NULL                     |
| 聚合函数       | COUNT(), AVG(), SUM(), MIN(), MAX()       | SELECT COUNT(*) FROM employee;                               | 自动忽略 NULL 值                  |
| 平均值         | AVG(列)                                   | SELECT AVG(Age) FROM employee;                               |                                   |
| 求和           | SUM(列)                                   | SELECT SUM(Age) FROM employee WHERE Workplace='Shenzhen';    | 统计时要配合 WHERE                |
| 最小值         | MIN(列)                                   | SELECT MIN(Age) FROM employee;                               |                                   |
| 最大值         | MAX(列)                                   | SELECT MAX(Age) FROM employee;                               |                                   |
| 分组           | GROUP BY 字段                             | SELECT Gender, COUNT(*) FROM employee GROUP BY Gender;       | 分组后只能查分组字段和聚合函数    |
| 分组+统计      | GROUP BY + 聚合函数                       | SELECT Gender, AVG(Age) FROM employee GROUP BY Gender;       |                                   |
| 分组过滤       | HAVING 条件                               | SELECT Workplace, COUNT(*) FROM employee GROUP BY Workplace HAVING COUNT(*) > 2; | HAVING 用在分组后，区别于 WHERE   |
| 排序           | ORDER BY 字段 [ASC/DESC]                  | SELECT * FROM employee ORDER BY Age ASC;                     | 默认升序（ASC）                   |
| 多字段排序     | ORDER BY 字段1, 字段2                      | SELECT * FROM employee ORDER BY Age ASC, Entrydate DESC;     | 先按第一个字段，再按第二个         |
| 分页           | LIMIT 偏移量, 条数                        | SELECT * FROM employee LIMIT 0,10;                           | 第一页：LIMIT 0,10，第二页：LIMIT 10,10 |
| 综合条件查询   | 多条件组合                                | SELECT * FROM employee WHERE Age <= 30 AND Gender = 'male' AND CHAR_LENGTH(Name) <= 5; | 可配合函数（如 CHAR_LENGTH）       |
| 执行顺序       | 1. FROM<br>2. WHERE<br>3. GROUP BY<br>4. HAVING<br>5. SELECT<br>6. ORDER BY<br>7. LIMIT | - | 理解执行顺序能帮你调试复杂 SQL |

---

## 📌 重点总结

- `WHERE` 是行级过滤，`HAVING` 是组级过滤
- `BETWEEN` 包含边界
- `IS NULL` 必须这样写，不能用 `= NULL`
- `LIMIT` 的第二个参数是条数，不是结束位置
- SQL 执行顺序和书写顺序不同

---