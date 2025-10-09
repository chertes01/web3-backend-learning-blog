# 📘 MySQL DDL & DML 学习笔记对照表

---

## 1️⃣ DDL（Data Definition Language，数据定义语言）

> 主要用于数据库/表的结构定义与修改。

| 操作           | 语法                                              | 示例                                                                 |
|----------------|---------------------------------------------------|----------------------------------------------------------------------|
| 创建数据库     | `CREATE DATABASE dbname;`                         | `CREATE DATABASE testdb;`                                            |
| 删除数据库     | `DROP DATABASE dbname;`                           | `DROP DATABASE testdb;`                                              |
| 选择数据库     | `USE dbname;`                                     | `USE testdb;`                                                        |
| 创建表         | `CREATE TABLE table_name (列名 类型 约束,...);`   | `CREATE TABLE userInformation ( ID INT PRIMARY KEY, Name VARCHAR(50), Age INT, Gender VARCHAR(10), idCare BIGINT );` |
| 新增列         | `ALTER TABLE table_name ADD 列名 类型;`           | `ALTER TABLE userInformation ADD Entrydate DATE;`                    |
| 修改列类型     | `ALTER TABLE table_name MODIFY 列名 新类型;`      | `ALTER TABLE userInformation MODIFY Age SMALLINT;`                   |
| 修改列名       | `ALTER TABLE table_name CHANGE 旧列名 新列名 类型;`| `ALTER TABLE userInformation CHANGE idCare idCard BIGINT;`           |
| 删除列         | `ALTER TABLE table_name DROP 列名;`               | `ALTER TABLE userInformation DROP Entrydate;`                        |
| 删除表         | `DROP TABLE table_name;`                          | `DROP TABLE userInformation;`                                        |
| 清空表         | `TRUNCATE TABLE table_name;`                      | `TRUNCATE TABLE userInformation;`                                    |

---

## 2️⃣ DML（Data Manipulation Language，数据操作语言）

> 主要用于表中数据的增删改查。

| 操作           | 语法                                              | 示例                                                                 |
|----------------|---------------------------------------------------|----------------------------------------------------------------------|
| 插入单行数据   | `INSERT INTO table VALUES (...);`                 | `INSERT INTO userInformation VALUES (1,'Alice',30,'female',440303733174533778);` |
| 插入多行数据   | `INSERT INTO table VALUES (...),(...);`           | `INSERT INTO userInformation VALUES (2,'Frank',38,'male',440303733317477835),(3,'bob',25,'male',440303774773533318);` |
| 指定列插入     | `INSERT INTO table (col1,col2,...) VALUES (...);` | `INSERT INTO userInformation (ID,Name) VALUES (1,'Hany');`           |
| SET 语法插入   | `INSERT INTO table SET col1=val1, col2=val2;`     | `INSERT INTO userInformation SET ID=1,Name='Alice',Age=30,Gender='female',idCare=440303733174533778,Entrydate='2020-01-10';` |
| 查询所有数据   | `SELECT * FROM table;`                            | `SELECT * FROM userInformation;`                                     |
| 查询指定列     | `SELECT col1,col2 FROM table;`                    | `SELECT Name,Age FROM userInformation;`                              |
| 条件查询       | `SELECT * FROM table WHERE 条件;`                 | `SELECT * FROM userInformation WHERE Gender='female';`               |
| 更新单列       | `UPDATE table SET col=val WHERE 条件;`            | `UPDATE userInformation SET Name='Bob' WHERE Name='bob';`            |
| 更新多列       | `UPDATE table SET col1=val1,col2=val2 WHERE 条件;`| `UPDATE userInformation SET ID=4,Age=23,Gender='female',idCare=440303733783374715 WHERE Name='Hany';` |
| 更新所有记录   | `UPDATE table SET col=val;`                       | `UPDATE userInformation SET Entrydate='2020-01-10';`                 |
| 删除指定行     | `DELETE FROM table WHERE 条件;`                   | `DELETE FROM userInformation WHERE Name='Alice';`                    |
| 删除所有记录   | `DELETE FROM table;`                              | `DELETE FROM userInformation;`                                       |

---

## 3️⃣ 小结

- **DDL**：改结构（建库建表/改表/删表）。
- **DML**：改数据（增删改查）。

> ⚠ 更新与删除时一定要加 WHERE，否则会影响整张表。