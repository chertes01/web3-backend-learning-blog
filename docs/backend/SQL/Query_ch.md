# SQL进阶查询学习笔记

这份笔记旨在整理和讲解SQL中常用且重要的进阶查询技巧，包括多表连接、联合查询和各种类型的子查询。所有示例都基于预设的数据表，方便读者直接运行和理解。

---

## 目录

- 前提准备：数据表结构与数据
- 多表连接查询 (JOIN)
  - 内连接 (INNER JOIN)
  - 左外连接 (LEFT JOIN)
  - 右外连接 (RIGHT JOIN)
  - 全外连接 (FULL OUTER JOIN)
  - 自连接 (SELF JOIN)
- 联合查询 (UNION)
- 子查询/嵌套查询 (Subquery)
  - 标量子查询 (Scalar Subquery)
  - 列子查询 (Column Subquery)
  - 行子查询 (Row Subquery)
  - 表子查询 (Table Subquery)
- 综合练习

---

## 前提准备：数据表结构与数据

### 核心表

- `student`: 学生表
- `course`: 课程表
- `student_course`: 学生和课程的关联表 (多对多)

### 练习用表

- `emp`: 员工表
- `dept`: 部门表
- `salgrade`: 薪资等级表

### SQL建表及数据

```sql
-- 学生表
CREATE TABLE student (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    gender ENUM('male', 'female')
) COMMENT '学生表';

-- 课程表
CREATE TABLE course (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cour VARCHAR(20)
) COMMENT '课程表';

-- 学生选课中间表
CREATE TABLE student_course (
    id_stu INT,
    id_cour INT,
    PRIMARY KEY (id_stu, id_cour),
    FOREIGN KEY (id_stu) REFERENCES student(id),
    FOREIGN KEY (id_cour) REFERENCES course(id)
) COMMENT '学生选课中间表';

-- 插入数据
INSERT INTO student (name, gender) VALUES
('Alice', 'female'), ('Bob', 'male'), ('Cindy', 'female'),
('David', 'male'), ('Eva', 'female'), ('Frank', 'male');

INSERT INTO course (cour) VALUES
('数据库系统'), ('Python程序设计'), ('高等数学'),
('计算机网络'), ('线性代代'), ('Go语言');

INSERT INTO student_course (id_stu, id_cour) VALUES
(1, 1), (1, 2), (1, 5),
(2, 2), (2, 3),
(3, 1), (3, 4),
(4, 2), (4, 3), (4, 4),
(5, 5);

-- 部门表
CREATE TABLE dept (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(20) NOT NULL
) COMMENT '部门表';

-- 员工表
CREATE TABLE emp (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(10) NOT NULL,
    age INT,
    job VARCHAR(20),
    salary INT,
    entrydate DATE,
    managerid INT,
    dept_id INT
) COMMENT '员工表';

-- 薪资等级表
CREATE TABLE salgrade (
    grade INT,
    losal INT,
    hisal INT
) COMMENT '薪资等级表';

-- 插入练习数据
INSERT INTO dept (name) VALUES ('研发部'), ('市场部'), ('财务部'), ('销售部');

INSERT INTO emp (name, age, job, salary, entrydate, managerid, dept_id) VALUES
('丁敏君', 28, 'DBA', 12000, '2012-03-01', 2, 1),
('赵敏', 35, '项目经理', 18000, '2008-08-08', NULL, 1),
('员工A', 45, 'Java开发', 14000, '2005-01-01', 2, 1),
('员工B', 52, '市场专员', 8000, '2010-06-01', 4, 2),
('员工C', 30, '财务', 9000, '2011-09-01', 4, 3),
('领导B', 48, '市场总监', 15000, '2006-07-01', NULL, 2);

INSERT INTO salgrade VALUES (1, 0, 3000), (2, 3001, 5000), (3, 5001, 8000),
(4, 8001, 10000), (5, 10001, 15000), (6, 15001, 20000),
(7, 20001, 25000), (8, 25001, 30000);
```

---

## 多表连接查询 (JOIN)

当需要的数据分散在多个表中时，需要使用JOIN将这些表按照指定的关联条件连接起来。

### 内连接 (INNER JOIN)

**定义**：返回两个表中连接字段相匹配的行。

**语法**：

- 隐式内连接（WHERE子句）：
  ```sql
  SELECT a.col, b.col FROM tableA a, tableB b WHERE a.id = b.id;
  ```
- 显式内连接（推荐）：
  ```sql
  SELECT a.col, b.col FROM tableA a INNER JOIN tableB b ON a.id = b.id;
  ```

**示例**：查询员工的姓名、年龄、职位、部门信息。
```sql
-- 隐式内连接
SELECT emp.name, emp.age, emp.job, dept.name
FROM emp, dept
WHERE dept.id = emp.dept_id;

-- 显式内连接
SELECT emp.name, emp.age, emp.job, dept.name
FROM emp JOIN dept ON emp.dept_id = dept.id;
```

---

### 左外连接 (LEFT JOIN)

**定义**：返回左表中的所有记录，以及右表中连接字段相匹配的记录。如果右表中没有匹配项，则结果为 NULL。

**语法**：
```sql
SELECT a.col, b.col FROM tableA a LEFT JOIN tableB b ON a.id = b.id;
```

**示例**：查询所有学生及其选课情况（未选课的学生也需显示）。
```sql
SELECT s.id, s.name, c.cour
FROM student s
LEFT JOIN student_course sc ON s.id = sc.id_stu
LEFT JOIN course c ON sc.id_cour = c.id
ORDER BY s.id, c.id;
```

---

### 右外连接 (RIGHT JOIN)

**定义**：与LEFT JOIN相反，返回右表中的所有记录，以及左表中连接字段相匹配的记录。如果左表中没有匹配项，则结果为 NULL。

**语法**：
```sql
SELECT a.col, b.col FROM tableA a RIGHT JOIN tableB b ON a.id = b.id;
```

**示例**：查询所有课程及其被选情况（没有学生选的课程也需显示）。
```sql
SELECT s.name, c.cour
FROM student s
RIGHT JOIN student_course sc ON s.id = sc.id_stu
RIGHT JOIN course c ON sc.id_cour = c.id;
```

---

### 全外连接 (FULL OUTER JOIN)

**定义**：返回左表和右表中的所有记录。当某行在另一表中没有匹配时，另一表的选择列表列包含 NULL。

**注意**：MySQL 不直接支持 FULL OUTER JOIN，可以通过 LEFT JOIN 和 RIGHT JOIN 的 UNION 来实现。

**模拟语法**：
```sql
SELECT * FROM tableA a LEFT JOIN tableB b ON a.id = b.id
UNION
SELECT * FROM tableA a RIGHT JOIN tableB b ON a.id = b.id;
```

---

### 自连接 (SELF JOIN)

**定义**：一张表与它自身进行连接。常用于处理表内具有层级关系的数据（如员工与经理）。

**示例**：查询emp表中员工及其所属领导。
```sql
-- 仅显示有领导的员工 (内连接)
SELECT a.id AS '员工ID', a.name AS '员工姓名', b.name AS '领导姓名'
FROM emp a
JOIN emp b ON a.managerid = b.id;

-- 显示所有员工，没有领导的员工领导姓名为NULL (左连接)
SELECT a.id AS '员工ID', a.name AS '员工姓名', b.name AS '领导姓名'
FROM emp a
LEFT JOIN emp b ON a.managerid = b.id;
```

---

## 联合查询 (UNION)

**定义**：用于合并两个或多个 SELECT 语句的结果集。

**规则**：

- 所有 SELECT 语句必须拥有相同数量的列。
- 列也必须拥有相似的数据类型。
- 每条 SELECT 语句中的列的顺序必须相同。

**语法**：

- `UNION`：合并结果集，并自动去除重复的行。
- `UNION ALL`：合并结果集，但保留所有行，包括重复行，效率更高。

**示例**：将emp表中薪资低于10000的员工，和年龄大于50岁的员工全部查询出来。
```sql
-- UNION ALL: 包含重复数据
SELECT * FROM emp WHERE salary < 10000
UNION ALL
SELECT * FROM emp WHERE age > 50;

-- UNION: 自动去重
SELECT * FROM emp WHERE salary < 10000
UNION
SELECT * FROM emp WHERE age > 50;
```

---

## 子查询/嵌套查询 (Subquery)

**定义**：嵌套在其他SQL语句（如SELECT, INSERT, UPDATE, DELETE）中的查询。

### 标量子查询 (Scalar Subquery)

**定义**：子查询返回的结果是单个值（一行一列）。

**示例**：查询"研发部"的所有员工信息。
```sql
SELECT * FROM emp WHERE dept_id = (SELECT id FROM dept WHERE name = '研发部');
```

**示例2**：查询在"丁敏君"入职之后的员工信息。
```sql
SELECT * FROM emp
WHERE entrydate > (SELECT entrydate FROM emp WHERE name = '丁敏君');
```

---

### 列子查询 (Column Subquery)

**定义**：子查询返回的结果是一列多行。

**常用操作符**：

- `IN`：等于列表中的任意一个。
- `ANY`：与子查询返回的任意一个值进行比较。
- `ALL`：与子查询返回的所有值进行比较。

**示例1**：查询“市场部”和“财务部”的所有员工信息。
```sql
SELECT * FROM emp
WHERE dept_id IN (SELECT id FROM dept WHERE name IN ('市场部', '财务部'));
```

**示例2**：查询比研发部所有人工资都高的员工信息。
```sql
SELECT * FROM emp
WHERE salary > ALL (SELECT salary FROM emp WHERE dept_id = (SELECT id FROM dept WHERE name = '研发部'));
-- 等价于: salary > (SELECT MAX(salary) FROM ... )
```

**示例3**：查询比研发部任意一人工资高的员工信息。
```sql
SELECT * FROM emp
WHERE salary > ANY (SELECT salary FROM emp WHERE dept_id = (SELECT id FROM dept WHERE name = '研发部'));
-- 等价于: salary > (SELECT MIN(salary) FROM ... )
```

---

### 行子查询 (Row Subquery)

**定义**：子查询返回的结果是一行多列。

**示例**：查询入职时间和工资均与“赵敏”相同的员工。
```sql
SELECT * FROM emp
WHERE (salary, entrydate) = (SELECT salary, entrydate FROM emp WHERE name = '赵敏');
```

---

### 表子查询 (Table Subquery)

**定义**：子查询返回的结果是多行多列，即一个虚拟表。

**关键**：返回的虚拟表通常放在FROM子句中，并且必须为其指定一个别名。

**示例**：查询入职日期是2010-06-01之后的员工，及其部门信息。
```sql
SELECT e.*, d.name AS dept_name
FROM (SELECT * FROM emp WHERE entrydate > '2010-06-01') e
JOIN dept d ON e.dept_id = d.id;
```

---

## 综合练习

**练习1**：查询所有学生的选课情况，展示出学生名称，学号，课程名称。
```sql
SELECT s.name, s.id, c.cour
FROM student s
LEFT JOIN student_course sc ON s.id = sc.id_stu
LEFT JOIN course c ON sc.id_cour = c.id;
```

**练习2**：统计每门课程的选课人数。
```sql
SELECT c.cour, COUNT(sc.id_stu) AS '选课人数'
FROM course c
LEFT JOIN student_course sc ON c.id = sc.id_cour
GROUP BY c.id, c.cour;
```

**练习3**：查询所有员工的工资等级。
```sql
SELECT e.name, e.salary, s.grade
FROM emp e
JOIN salgrade s ON e.salary BETWEEN s.losal AND s.hisal;
```

**练习4**：查询"研发部"所有员工的信息及工资等级。
```sql
SELECT e.*, s.grade
FROM emp e
JOIN salgrade s ON e.salary BETWEEN s.losal AND s.hisal
WHERE e.dept_id = (SELECT id FROM dept WHERE name = '研发部');
```

**练习5**：查询工资比"丁敏君"高的员工信息。
```sql
SELECT * FROM emp WHERE salary > (SELECT salary FROM emp WHERE name = '丁敏君');
```

**练习6**：查询低于本部门平均工资的员工信息。
```sql
SELECT *
FROM emp e1
WHERE e1.salary < (SELECT AVG(e2.salary) FROM emp e2 WHERE e2.dept_id = e1.dept_id);
```

**练习7**：查询所有的部门信息，并统计每个部门的员工人数。
```sql
-- 使用标量子查询
SELECT
    d.id,
    d.name,
    (SELECT COUNT(*) FROM emp e WHERE e.dept_id = d.id) AS '员工人数'
FROM dept d;

-- 使用JOIN和GROUP BY (效率通常更高)
SELECT d.id, d.name, COUNT(e.id) AS '员工人数'
FROM dept d
LEFT JOIN emp e ON d.id = e.dept_id
GROUP BY d.id, d.name;
```

---