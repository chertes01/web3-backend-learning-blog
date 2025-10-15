# MySQL 存储过程与程序设计学习笔记

这是一份关于 MySQL 存储过程、变量、流程控制和游标的综合学习笔记。旨在帮助开发者理解和掌握在 MySQL 中进行数据库层面的程序设计，以提升应用性能和代码复用性。

---

## 目录

1. [什么是存储过程？](#1-什么是存储过程)
2. [存储过程 vs 函数](#2-存储过程-vs-函数)
3. [基本语法](#3-基本语法)
   - [创建、调用、查看、删除](#创建调用查看删除)
   - [关于 DELIMITER](#关于-delimiter)
4. [变量详解](#4-变量详解)
   - [系统变量](#系统变量)
   - [用户自定义变量](#用户自定义变量)
   - [局部变量](#局部变量)
   - [三种变量对比](#三种变量对比)
5. [参数传递 (IN, OUT, INOUT)](#5-参数传递-in-out-inout)
6. [流程控制](#6-流程控制)
   - [IF 语句](#if-语句)
   - [CASE 语句](#case-语句)
   - [循环语句 (WHILE, REPEAT, LOOP)](#循环语句-while-repeat-loop)
7. [游标 (Cursor)](#7-游标-cursor)
   - [游标的概念](#游标的概念)
   - [使用步骤与示例](#使用步骤与示例)
8. [补充知识点与最佳实践](#8-补充知识点与最佳实践)

---

## 1. 什么是存储过程？

存储过程是事先经过编译并存储在数据库中的一段 SQL 语句的集合。调用存储过程可以简化应用开发人员的工作，减少数据在数据库和应用服务器之间的传输，对于提高数据处理的效率大有裨益。

**特点：**
- **封装与复用：** 将复杂的业务逻辑封装在一起，便于重复使用。
- **灵活的参数：** 可以接收输入参数，也可以返回数据作为输出参数。
- **提升性能：** 减少了客户端与服务器之间的网络交互次数。执行一次存储过程比执行多条独立的 SQL 语句效率更高。

---

## 2. 存储过程 vs 函数

| 特性             | 存储过程 (Procedure)                     | 函数 (Function)                     |
|------------------|------------------------------------------|-------------------------------------|
| **返回值**       | 不能直接 `RETURN` 一个值，需通过 OUT/INOUT 参数返回值 | 必须 `RETURN` 一个单一的值         |
| **调用方式**     | `CALL procedure_name();`                | `SELECT function_name();` 或嵌入到 SQL 语句中 |
| **用途**         | 处理复杂的数据操作、事务、业务逻辑集合    | 主要用于计算并返回一个结果，常嵌入到 SELECT 语句中 |
| **SELECT 语句**  | 内部可以使用 SELECT 语句                 | 内部不允许直接使用 SELECT 语句返回结果集 |

---

## 3. 基本语法

### 创建、调用、查看、删除

```sql
-- 准备数据库
USE testdb;

-- 创建一个简单的存储过程
CREATE PROCEDURE p1()
BEGIN
    SELECT COUNT(*) FROM student;
END;

-- 调用存储过程
CALL p1();

-- 查看存储过程
-- 方式一：从 information_schema 中查询
SELECT * FROM information_schema.ROUTINES WHERE ROUTINE_SCHEMA = 'testdb';

-- 方式二：查看创建语句
SHOW CREATE PROCEDURE p1;

-- 删除存储过程
DROP PROCEDURE IF EXISTS p1;
```

### 关于 DELIMITER

在命令行或图形化工具中创建存储过程时，必须临时修改 SQL 的结束符。因为存储过程内部可能包含多个以分号 `;` 结束的语句，客户端读到第一个分号时会误以为创建语句已结束。

```sql
-- 1. 将结束符改为 $$
DELIMITER $$

CREATE PROCEDURE p1()
BEGIN
    SELECT COUNT(*) FROM student;
    SELECT 'Hello';
END$$ -- 2. 使用新的结束符

-- 3. 将结束符恢复为 ;
DELIMITER ;
```

---

## 4. 变量详解

### 系统变量

由 MySQL 服务器提供，非用户定义，分为 GLOBAL（全局）和 SESSION（会话）两个级别。

```sql
-- 查看会话级别的系统变量 (两种方式)
SHOW SESSION VARIABLES LIKE 'autocommit';
SELECT @@session.autocommit;

-- 查看全局级别的系统变量 (两种方式)
SHOW GLOBAL VARIABLES LIKE 'autocommit';
SELECT @@global.autocommit;

-- 设置系统变量
SET SESSION autocommit = 1;
SET GLOBAL autocommit = 0; -- 此操作不会影响当前已存在的会话
```

**注意：** 通过 `SET GLOBAL` 修改的参数会在 MySQL 服务重新启动后失效。若要永久生效，需要写入 `my.cnf` (或 `my.ini`) 配置文件。

### 用户自定义变量

用户根据需要定义的变量，以 `@` 开头。它不需要提前声明，作用域为当前整个连接会话。

```sql
-- 赋值方式
SET @myname = 'Alice';
SET @myage := 18; -- := 也是赋值运算符
SET @mygender = 'female', @myhobby = 'Golang';

SELECT @mycolor := 'Yellow';
SELECT COUNT(*) INTO @mycount FROM student;

-- 读取变量
SELECT @myname, @myage, @mygender, @myhobby;
SELECT @mycolor, @mycount;
```

### 局部变量

在 `BEGIN...END` 代码块内生效的变量，必须先用 `DECLARE` 声明才能使用。

```sql
DELIMITER $$
CREATE PROCEDURE p2()
BEGIN
    -- 声明局部变量，并给予默认值
    DECLARE stu_count INT DEFAULT 0;
    
    -- 赋值
    SELECT COUNT(*) INTO stu_count FROM student;
    
    -- 查询
    SELECT stu_count;
END$$
DELIMITER ;

CALL p2();
```

### 三种变量对比

| 类型         | 语法                | 作用域           | 是否需要声明 |
|--------------|---------------------|------------------|--------------|
| **系统变量** | `@@scope.var_name` | GLOBAL / SESSION | 否           |
| **用户变量** | `@var_name`        | 当前会话         | 否           |
| **局部变量** | `var_name`         | `BEGIN...END` 块 | 是           |

---

## 5. 参数传递 (IN, OUT, INOUT)

- **IN：** 输入参数。将值传入程序，内部修改不影响外部。（默认模式）
- **OUT：** 输出参数。将程序内部的值返回给调用者。
- **INOUT：** 输入输出参数。既可传入，也可传出。

```sql
DELIMITER $$
-- OUT 示例：根据分数输出等级
CREATE PROCEDURE p4(IN score FLOAT, OUT result VARCHAR(10))
BEGIN
    IF score >= 85 THEN
        SET result = '优秀';
    ELSEIF score >= 60 THEN
        SET result = '及格';
    ELSE
        SET result = '不及格';
    END IF;
END$$

-- INOUT 示例：将 200 分制换算为百分制
CREATE PROCEDURE p5(INOUT score FLOAT)
BEGIN
    SET score = score / 2;
END$$
DELIMITER ;

-- 调用 OUT 参数的程序
CALL p4(99, @result); -- 使用用户变量 @result 来接收输出
SELECT @result;

-- 调用 INOUT 参数的程序
SET @score = 155;
CALL p5(@score); -- @score 既是输入也是输出
SELECT @score;
```

---

## 6. 流程控制

### IF 语句

```sql
DELIMITER $$
CREATE PROCEDURE p3_fixed(IN score FLOAT)
BEGIN
    DECLARE result VARCHAR(10);
    IF score >= 85 THEN
        SET result = '优秀';
    ELSEIF score >= 60 THEN
        SET result = '及格';
    ELSE
        SET result = '不及格';
    END IF;
    SELECT result;
END$$
DELIMITER ;

CALL p3_fixed(66);
```

### CASE 语句

```sql
DELIMITER $$
CREATE PROCEDURE p6(IN month INT)
BEGIN
    DECLARE result VARCHAR(10);
    CASE
        WHEN month >= 1 AND month <= 3 THEN SET result = '第一季度';
        WHEN month >= 4 AND month <= 6 THEN SET result = '第二季度';
        WHEN month >= 7 AND month <= 9 THEN SET result = '第三季度';
        WHEN month >= 10 AND month <= 12 THEN SET result = '第四季度';
        ELSE SET result = '非法月份';
    END CASE;
    SELECT CONCAT('输入月份：', month, ', 所属季度：', result) AS '结果';
END$$
DELIMITER ;

CALL p6(6);
```

---

## 7. 游标 (Cursor)

### 游标的概念

游标用于逐行处理查询返回的结果集。当需要对查询出的每一行数据进行独立的、复杂的处理时，游标非常有用。

### 使用步骤与示例

1. **声明游标：** 定义游标要处理的 SELECT 语句。
2. **声明 NOT FOUND 处理程序：** 定义在游标读取完所有行后要执行的动作。
3. **开启游标：** `OPEN`
4. **提取数据：** `FETCH`，通常在循环中进行。
5. **关闭游标：** `CLOSE`

```sql
DELIMITER $$
CREATE PROCEDURE p11_process_users(IN uage INT)
BEGIN
    DECLARE uname VARCHAR(100);
    DECLARE upro VARCHAR(100);
    DECLARE done INT DEFAULT FALSE;

    DECLARE u_cursor CURSOR FOR
        SELECT name, profession FROM tb_user WHERE age <= uage;

    DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;

    CREATE TABLE IF NOT EXISTS tb_user_pro (
        id INT PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(100),
        profession VARCHAR(100)
    );

    OPEN u_cursor;

    read_loop: LOOP
        FETCH u_cursor INTO uname, upro;
        IF done THEN
            LEAVE read_loop;
        END IF;
        INSERT INTO tb_user_pro (name, profession) VALUES (uname, upro);
    END LOOP;

    CLOSE u_cursor;

    SELECT '数据处理完成！' AS '状态';
END$$
DELIMITER ;

CALL p11_process_users(29);
```

---

## 8. 补充知识点与最佳实践

### SQL 安全性 (DEFINER vs INVOKER)

- **DEFINER：** 存储过程以创建者的权限执行。
- **INVOKER：** 存储过程以调用者的权限执行。

```sql
CREATE SQL SECURITY INVOKER PROCEDURE ...
```

### 最佳实践

1. **命名规范：** 存储过程以 `sp_` 开头，参数以 `p_` 开头，变量以 `v_` 开头。
2. **添加注释：** 对复杂的逻辑进行说明。
3. **错误处理：** 使用 `DECLARE ... HANDLER` 捕获和处理潜在错误。
4. **避免 DDL 操作：** 尽量将数据定义和数据操作分离。
5. **事务管理：** 对一系列的 `INSERT/UPDATE/DELETE` 操作，使用 `START TRANSACTION`、`COMMIT`、`ROLLBACK` 确保数据一致性。