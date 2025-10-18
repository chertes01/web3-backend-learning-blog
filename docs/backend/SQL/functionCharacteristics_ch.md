# MySQL 自定义函数特性（Function Characteristics）学习笔记

在 MySQL 中创建自定义函数时，可以为其指定一些“特性”（Characteristics），这些特性能帮助 MySQL 优化器理解函数行为。合理声明函数特性不仅是良好的编程习惯，还能提升性能和数据一致性。

---

## 为什么要声明函数特性？

1. **优化器提示**  
   优化器可根据函数特性选择更优执行计划。例如，`DETERMINISTIC`（确定性）函数的结果可被缓存，重复调用时直接复用结果。

2. **复制安全**  
   在主从复制环境下，`DETERMINISTIC` 特性尤为重要。对于 `STATEMENT` 格式的二进制日志，声明为确定性的函数可保证主从一致。

3. **安全性与权限**  
   特性声明了函数的数据访问级别。例如，`READS SQL DATA` 只需 `SELECT` 权限，无需 `UPDATE` 或 `DELETE` 权限。

---

## 函数特性详解

### 1. DETERMINISTIC（确定性）

- **定义**：相同输入参数，总是产生相同结果。
- **示例**：

```sql
-- 计算从 1 到 n 的整数和，结果只依赖于输入参数 n
CREATE FUNCTION fun1(n INT)
RETURNS INT
DETERMINISTIC
BEGIN
    DECLARE total INT DEFAULT 0; -- 初始化累加变量
    WHILE n > 0 DO
        SET total := total + n; -- 累加当前 n
        SET n := n - 1;         -- n 递减
    END WHILE;
    RETURN total;               -- 返回最终结果
END;

-- 调用
SELECT fun1(100); -- 结果总是 5050
```

**字段讲解**  
- `RETURNS INT`：声明函数返回值类型为整数，需与 RETURN 语句类型一致。  
- `DETERMINISTIC`：声明函数为确定性，只要输入参数相同，返回值必定相同。  
- `RETURN total`：结束函数并返回最终结果。

> 注意：声明为 DETERMINISTIC 的函数若包含非确定性逻辑（如 NOW()、RAND()），可能导致复制警告或不一致。

---

### 2. NO SQL

- **定义**：函数体不包含任何 SQL 语句，仅逻辑或数学运算。
- **示例**：

```sql
-- 将字符串 s 重复 n 次拼接
CREATE FUNCTION repeat_string(s VARCHAR(255), n INT)
RETURNS TEXT
DETERMINISTIC
NO SQL
BEGIN
    DECLARE result TEXT DEFAULT ''; -- 初始化结果字符串
    DECLARE i INT DEFAULT 0;        -- 循环计数器
    IF n <= 0 THEN
        RETURN '';                  -- 若 n 不合法，直接返回空串
    END IF;
    WHILE i < n DO
        SET result = CONCAT(result, s); -- 拼接字符串
        SET i = i + 1;                 -- 计数器递增
    END WHILE;
    RETURN result;                     -- 返回拼接结果
END;

-- 调用
SELECT repeat_string('DB_', 3); -- 结果: DB_DB_DB_
```

**字段讲解**  
- `RETURNS TEXT`：声明返回类型为文本字符串。  
- `DETERMINISTIC`：确定性，输入相同结果相同。  
- `NO SQL`：函数体不包含任何 SQL 语句，仅做字符串拼接和循环。

---

### 3. READS SQL DATA

- **定义**：函数体包含读取数据的语句（如 SELECT），但不修改数据。
- **示例**：

```sql
-- 创建员工表并插入数据
CREATE TABLE employees (
    id INT PRIMARY KEY,
    name VARCHAR(50),
    salary DECIMAL(10, 2)
);

INSERT INTO employees (id, name, salary) VALUES
(1, 'Alice', 50000.00),
(2, 'Bob', 65000.00);

-- 根据员工ID查询薪水
CREATE FUNCTION get_employee_salary(emp_id INT)
RETURNS DECIMAL(10, 2)
READS SQL DATA
BEGIN
    DECLARE emp_salary DECIMAL(10, 2); -- 临时变量存储薪水
    SELECT salary INTO emp_salary FROM employees WHERE id = emp_id; -- 查询薪水
    RETURN emp_salary; -- 返回结果
END;

-- 调用
SELECT get_employee_salary(2); -- 结果: 65000.00
```

**字段讲解**  
- `RETURNS DECIMAL(10, 2)`：声明返回类型为高精度十进制，适合存储薪水等金额。  
- `READS SQL DATA`：声明函数只读取数据，不修改。  
- `DECLARE emp_salary DECIMAL(10, 2)`：声明临时变量用于存储查询结果。  
- `SELECT ... INTO ...`：将查询结果赋值给变量，要求只返回一行。

> 包含 READS SQL DATA 的函数通常是非确定性的，因为表数据可能被其他会话修改。

---

### 4. CONTAINS SQL

- **定义**：函数体包含 SQL 语句，但不读写表数据。例如 SET @x = 1。
- **示例**：

```sql
-- 返回当前数据库版本号
CREATE FUNCTION get_db_version()
RETURNS VARCHAR(255)
CONTAINS SQL
BEGIN
    RETURN VERSION(); -- 调用内置函数，不涉及表数据
END;

-- 调用
SELECT get_db_version(); -- 结果示例: 8.0.32
```

**字段讲解**  
- `RETURNS VARCHAR(255)`：声明返回类型为字符串。  
- `CONTAINS SQL`：函数体包含 SQL，但不涉及表数据。  
- `RETURN VERSION()`：返回 MySQL 版本号。

> 未指定特性时，MySQL 默认使用 CONTAINS SQL。

---

### 5. MODIFIES SQL DATA

- **定义**：函数体包含修改数据的语句（如 INSERT、UPDATE、DELETE）。
- **示例**：

```sql
-- 给指定员工涨薪，并返回新薪水
CREATE FUNCTION grant_raise(emp_id INT, raise_amount DECIMAL(10, 2))
RETURNS DECIMAL(10, 2)
MODIFIES SQL DATA
BEGIN
    UPDATE employees
    SET salary = salary + raise_amount
    WHERE id = emp_id; -- 修改员工薪水
    RETURN (SELECT salary FROM employees WHERE id = emp_id); -- 返回新薪水
END;

-- 调用前查询
SELECT get_employee_salary(1); -- 结果: 50000.00

-- 调用涨薪函数
SELECT grant_raise(1, 5000.00); -- 结果: 55000.00

-- 调用后再次查询
SELECT get_employee_salary(1); -- 结果: 55000.00
```

**字段讲解**  
- `RETURNS DECIMAL(10, 2)`：返回类型为高精度金额。  
- `MODIFIES SQL DATA`：声明函数会修改表数据。  
- `UPDATE ...`：修改员工薪水。  
- `RETURN (SELECT ...)`：返回修改后的新薪水。

> 注意：存储函数中直接执行 UPDATE 语句可能受 `log_bin_trust_function_creators` 参数限制。

---

## 总结与对比

| 特性                | 含义                   | SQL 示例                      |
|---------------------|------------------------|-------------------------------|
| NO SQL              | 完全不包含 SQL         | DECLARE, SET var = 1, WHILE   |
| CONTAINS SQL        | 包含 SQL，不读写表数据 | SET @x = NOW(), SELECT VERSION() |
| READS SQL DATA      | 读表数据               | SELECT ... FROM my_table      |
| MODIFIES SQL DATA   | 写表数据               | INSERT, UPDATE, DELETE        |

- 层级关系：NO SQL 最严格，MODIFIES SQL DATA 最宽松。MODIFIES SQL DATA 隐含 READS SQL DATA 和 CONTAINS SQL 能力。
- 默认值：未指定时，MySQL 默认使用 CONTAINS SQL。