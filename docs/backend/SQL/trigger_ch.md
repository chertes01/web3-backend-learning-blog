# MySQL 学习笔记：触发器 (Triggers)

## 1. 什么是触发器？

触发器（Trigger）是与表有关的数据库对象，它指在 `INSERT`、`UPDATE`、`DELETE` 事件之前(BEFORE)或之后(AFTER)，自动触发并执行触发器中定义的 SQL 语句集合。

触发器的这种特性可以协助应用在数据库端确保数据的完整性、执行日志记录、数据校验等操作。

## 2. 核心关键字：NEW 与 OLD

在触发器中，我们可以使用别名 `OLD` 和 `NEW` 来引用发生变化的记录内容。MySQL 目前仅支持行级触发（`FOR EACH ROW`），不支持语句级触发。

| 触发器类型 | NEW 和 OLD 的作用                                                              |
| :--------- | :------------------------------------------------------------------------------- |
| `INSERT`   | `NEW` 表示将要或已经新增的数据                                                   |
| `UPDATE`   | `OLD` 表示修改之前的数据 <br> `NEW` 表示将要或已经修改后的数据                  |
| `DELETE`   | `OLD` 表示将要或者已经删除的数据                                                 |

## 3. 实战案例：记录 tb_user 表的操作日志

我们的目标是：当 `tb_user` 表发生 `INSERT`、`UPDATE`、`DELETE` 操作时，自动将操作记录插入到 `user_logs` 日志表中。

### 步骤一：创建日志表

首先，我们创建用于存放日志的 `user_logs` 表。

```sql
USE indexTest;

-- 创建日志表
CREATE TABLE user_logs(
    id INT(10) NOT NULL AUTO_INCREMENT PRIMARY KEY ,
    opratetion VARCHAR(10) NOT NULL COMMENT '操作类型，insert/update/delate',
    operate_time DATETIME NOT NULL COMMENT '操作时间',
    operate_id INT(10) NOT NULL COMMENT '操作ID',
    operate_params VARCHAR(500) COMMENT '操作参数'
) ENGINE=INNODB DEFAULT CHARSET=UTF8;
```

### 步骤二：创建触发器

> ⚠️ **重要：关于 DELIMITER**
> 在创建触发器时，我们必须使用 `BEGIN...END` 语句块，而 `BEGIN...END` 内部的 SQL 语句需要用分号 `;` 结束。这会导致一个问题：MySQL 命令行默认遇到第一个分号时，会认为语句已经结束，从而导致 `CREATE TRIGGER` 语法错误。
> **解决方法**： 使用 `DELIMITER` 命令临时将语句结束符从 `;` 修改为其他符号（如 `//` 或 `$$`）。

#### A. INSERT 触发器

创建 `tb_user_insert_trigger`，在向 `tb_user` 表插入数据之后 (`AFTER INSERT`)，向 `user_logs` 插入日志。

```sql
-- 1. 临时修改结束符
DELIMITER //

CREATE TRIGGER tb_user_insert_trigger
    AFTER INSERT ON tb_user FOR EACH ROW
BEGIN
    -- 注意：CONCAT 函数内的每个参数都必须用英文逗号 , 隔开
    INSERT INTO user_logs(id, opratetion, operate_time, operate_id, operate_params)
    VALUES(NULL, 'insert', NOW(), NEW.id,
           CONCAT('id=', NEW.id, ',name=', NEW.name, ', phone=', NEW.phone, ', email=', NEW.email, ',profession=', NEW.profession)
          );
END//

-- 2. 将结束符改回分号
DELIMITER ;
```

**测试 INSERT 触发器：**

```sql
INSERT INTO tb_user(id, name, phone, email, profession, age, gender, status, createtime)
VALUES(15, '丁针', '13524746921', 'xuebaobizhui@frw.com', '香烟品鉴', 18, 'male', '1', NOW());

-- 执行后，可以去 user_logs 表中查看是否已自动生成日志
SELECT * FROM user_logs WHERE operate_id = 15;
```

#### B. UPDATE 触发器

创建 `tb_user_update_trigger`，在 `tb_user` 表更新数据之后 (`AFTER UPDATE`)，记录旧数据（`OLD`）和新数据（`NEW`）。

```sql
DELIMITER //

CREATE TRIGGER tb_user_update_trigger
    AFTER UPDATE ON tb_user FOR EACH ROW
BEGIN
    INSERT INTO user_logs(id, opratetion, operate_time, operate_id, operate_params)
    VALUES(NULL, 'update', NOW(), NEW.id,
           -- 将 OLD 和 NEW 的值都记录下来，用 | 分隔
           CONCAT('OLD: id=', OLD.id, ',name=', OLD.name, ', phone=', OLD.phone, ', email=', OLD.email, ',profession=', OLD.profession,
                  ' | NEW: id=', NEW.id, ',name=', NEW.name, ', phone=', NEW.phone, ', email=', NEW.email, ',profession=', NEW.profession)
          );
END//

DELIMITER ;
```

**测试 UPDATE 触发器：**

```sql
UPDATE tb_user SET name='丁针珍珠' WHERE id=15;

-- 再次查询日志表
SELECT * FROM user_logs WHERE operate_id = 15 AND opratetion = 'update';
```

#### C. DELETE 触发器

创建 `tb_user_delete_trigger`，在 `tb_user` 表删除数据之后 (`AFTER DELETE`)，记录被删除的数据（`OLD`）。

```sql
DELIMITER //

CREATE TRIGGER tb_user_delete_trigger
    AFTER DELETE ON tb_user FOR EACH ROW
BEGIN
    INSERT INTO user_logs(id, opratetion, operate_time, operate_id, operate_params)
    VALUES(NULL, 'delete', NOW(), OLD.id,
           CONCAT('id=', OLD.id, ',name=', OLD.name, ', phone=', OLD.phone, ', email=', OLD.email, ',profession=', OLD.profession)
          );
END//

DELIMITER ;
```

**测试 DELETE 触发器：**

```sql
DELETE FROM tb_user WHERE id=15;

-- 再次查询日志表
SELECT * FROM user_logs WHERE operate_id = 15 AND opratetion = 'delete';
```

## 4. 如何管理触发器

### A. 查看触发器

查看当前数据库中所有的触发器：

```sql
SHOW TRIGGERS;
```

### B. 删除触发器

使用 `DROP TRIGGER` 命令删除触发器。推荐使用 `IF EXISTS` 来避免在触发器不存在时报错。

```sql
DROP TRIGGER IF EXISTS tb_user_delete_trigger;
```

### C. 更新触发器

MySQL 没有提供 `ALTER TRIGGER` 语句来直接修改触发器。更新一个触发器的标准做法是 “先删除，再重新创建”。

例如，如果你想修改 `tb_user_insert_trigger` 的逻辑：

```sql
-- 1. 先删除旧的
DROP TRIGGER IF EXISTS tb_user_insert_trigger;

-- 2. 使用 DELIMITER 重新创建新的
DELIMITER //
CREATE TRIGGER tb_user_insert_trigger
    AFTER INSERT ON tb_user FOR EACH ROW
BEGIN
    -- ... 在这里写入你更新后的新逻辑 ...
    INSERT INTO user_logs(id, opratetion, operate_time, operate_id, operate_params)
    VALUES(NULL, '新增', NOW(), NEW.id,  -- 比如把 'insert' 改为 '新增'
           CONCAT('id=', NEW.id, ',name=', NEW.name)
          );
END//
DELIMITER ;
```