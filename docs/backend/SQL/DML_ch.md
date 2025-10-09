📘 MySQL DML 学习笔记
1️⃣ 插入数据（INSERT）

插入单行数据：

INSERT INTO userInformation 
VALUES (1,'Alice',30,'female',440303733174533778);


插入多行数据：

INSERT INTO userInformation 
VALUES
(2,'Frank',38,'male',440303733317477835),
(3,'bob',25,'male',440303774773533318);


指定列插入（推荐，避免因表结构变化出错）：

INSERT INTO userInformation (ID, Name)
VALUES (1,'Hany');


使用 SET 语法插入：

INSERT INTO userInformation 
SET ID=1, Name='Alice', Age=30, Gender='female', 
    idCare=440303733174533778, Entrydate='2020-01-10';


⚠ 注意：

一般建议用 指定列方式，避免表结构有新增列时报错。

插入时应确保主键或唯一键不重复。

2️⃣ 查询数据（SELECT）

查询所有数据：

SELECT * FROM userInformation;


查询指定列：

SELECT Name, Age FROM userInformation;


条件查询：

SELECT * FROM userInformation WHERE Gender='female';

3️⃣ 更新数据（UPDATE）

更新单个字段：

UPDATE userInformation 
SET Name='Bob' 
WHERE Name='bob';


更新多个字段：

UPDATE userInformation 
SET ID=4, Age=23, Gender='female', 
    idCare=440303733783374715 
WHERE Name='Hany';


更新主键（注意唯一性）：

UPDATE userInformation 
SET ID=3 
WHERE Name='Bob';


更新所有行（慎用）：

UPDATE userInformation 
SET Entrydate='2020-01-10';


⚠ 注意：

必须带 WHERE 条件，否则会更新整张表！

更新主键时要确认不会违反唯一约束。

4️⃣ 删除数据（DELETE）

删除指定行：

DELETE FROM userInformation 
WHERE Name='Alice';


删除所有数据（保留表结构）：

DELETE FROM userInformation;


更高效清空整表：

TRUNCATE TABLE userInformation;


⚠ 注意：

DELETE 是逐行删除，效率较低；

TRUNCATE 直接清空表，速度快，但不能回滚。

5️⃣ 小结

INSERT：新增记录（可单行、多行、指定列）。

SELECT：查询数据（可配合条件、投影）。

UPDATE：修改数据（一定要注意 WHERE 条件）。

DELETE：删除记录（区分 DELETE 与 TRUNCATE）。

