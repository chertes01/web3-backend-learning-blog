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
CREATE FUNCTION fun1(n INT)
RETURNS INT
DETERMINISTIC
BEGIN
    DECLARE total INT DEFAULT 0;
    WHILE n > 0 DO
        SET total := total + n;
        SET n := n - 1;
    END WHILE;
    RETURN total;
END;

-- 调用
SELECT fun1(100); -- 结果总是 5050
```

> 注意：声明为 DETERMINISTIC 的函数若包含非确定性逻辑（如 NOW()、RAND()），可能导致复制警告或不一致。

---

### 2. NO SQL

- **定义**：函数体不包含任何 SQL 语句，仅逻辑或数学运算。
- **示例**：

```sql
CREATE FUNCTION repeat_string(s VARCHAR(255), n INT)
RETURNS TEXT
DETERMINISTIC
NO SQL
BEGIN
    DECLARE result TEXT DEFAULT '';
    DECLARE i INT DEFAULT 0;
    IF n <= 0 THEN
        RETURN '';
    END IF;
    WHILE i < n DO
        SET result = CONCAT(result, s);
        SET i = i + 1;
    END WHILE;
    RETURN result;
END;

-- 调用
SELECT repeat_string('DB_', 3); -- 结果: DB_DB_DB_
```

---

### 3. READS SQL DATA

- **定义**：函数体包含读取数据的语句（如 SELECT），但不修改数据。
- **示例**：

```sql
CREATE TABLE employees (
    id INT PRIMARY KEY,
    name VARCHAR(50),
    salary DECIMAL(10, 2)
);
INSERT INTO employees (id, name, salary) VALUES
(1, 'Alice', 50000.00),
(2, 'Bob', 65000.00);

CREATE FUNCTION get_employee_salary(emp_id INT)
RETURNS DECIMAL(10, 2)
READS SQL DATA
BEGIN
    DECLARE emp_salary DECIMAL(10, 2);
    SELECT salary INTO emp_salary FROM employees WHERE id = emp_id;
    RETURN emp_salary;
END;

-- 调用
SELECT get_employee_salary(2); -- 结果: 65000.00
```

> 包含 READS SQL DATA 的函数通常是非确定性的，因为表数据可能被其他会话修改。

---

### 4. CONTAINS SQL

- **定义**：函数体包含 SQL 语句，但不读写表数据。例如 SET @x = 1。
- **示例**：

```sql
CREATE FUNCTION get_db_version()
RETURNS VARCHAR(255)
CONTAINS SQL
BEGIN
    RETURN VERSION();
END;

-- 调用
SELECT get_db_version(); -- 结果示例: 8.0.32
```

> 未指定特性时，MySQL 默认使用 CONTAINS SQL。

---

### 5. MODIFIES SQL DATA

- **定义**：函数体包含修改数据的语句（如 INSERT、UPDATE、DELETE）。
- **示例**：

```sql
CREATE FUNCTION grant_raise(emp_id INT, raise_amount DECIMAL(10, 2))
RETURNS DECIMAL(10, 2)
MODIFIES SQL DATA
BEGIN
    UPDATE employees
    SET salary = salary + raise_amount
    WHERE id = emp_id;
    RETURN (SELECT salary FROM employees WHERE id = emp_id);
END;

-- 调用前查询
SELECT get_employee_salary(1); -- 结果: 50000.00

-- 调用涨薪函数
SELECT grant_raise(1, 5000.00); -- 结果: 55000.00

-- 调用后再次查询
SELECT get_employee_salary(1); -- 结果: 55000.00
```

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