# Advanced SQL Query Study Notes

This note aims to organize and explain commonly used and important advanced SQL query techniques, including multi-table joins, union queries, and various types of subqueries. All examples are based on preset data tables for easy execution and understanding.

---

## Table of Contents

- Preparation: Table Structures and Data
- Multi-table Join Queries (JOIN)
  - Inner Join (INNER JOIN)
  - Left Outer Join (LEFT JOIN)
  - Right Outer Join (RIGHT JOIN)
  - Full Outer Join (FULL OUTER JOIN)
  - Self Join (SELF JOIN)
- Union Queries (UNION)
- Subqueries/Nested Queries (Subquery)
  - Scalar Subquery
  - Column Subquery
  - Row Subquery
  - Table Subquery
- Comprehensive Exercises

---

## Preparation: Table Structures and Data

### Core Tables

- `student`: Student table
- `course`: Course table
- `student_course`: Association table between students and courses (many-to-many)

### Practice Tables

- `emp`: Employee table
- `dept`: Department table
- `salgrade`: Salary grade table

### SQL Table Creation and Data

```sql
-- Student table
CREATE TABLE student (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(10) NOT NULL,
    gender ENUM('male', 'female')
) COMMENT 'Student table';

-- Course table
CREATE TABLE course (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cour VARCHAR(20)
) COMMENT 'Course table';

-- Student-course association table
CREATE TABLE student_course (
    id_stu INT,
    id_cour INT,
    PRIMARY KEY (id_stu, id_cour),
    FOREIGN KEY (id_stu) REFERENCES student(id),
    FOREIGN KEY (id_cour) REFERENCES course(id)
) COMMENT 'Student-course association table';

-- Insert data
INSERT INTO student (name, gender) VALUES
('Alice', 'female'), ('Bob', 'male'), ('Cindy', 'female'),
('David', 'male'), ('Eva', 'female'), ('Frank', 'male');

INSERT INTO course (cour) VALUES
('Database Systems'), ('Python Programming'), ('Advanced Mathematics'),
('Computer Networks'), ('Linear Algebra'), ('Go Language');

INSERT INTO student_course (id_stu, id_cour) VALUES
(1, 1), (1, 2), (1, 5),
(2, 2), (2, 3),
(3, 1), (3, 4),
(4, 2), (4, 3), (4, 4),
(5, 5);

-- Department table
CREATE TABLE dept (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(20) NOT NULL
) COMMENT 'Department table';

-- Employee table
CREATE TABLE emp (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(10) NOT NULL,
    age INT,
    job VARCHAR(20),
    salary INT,
    entrydate DATE,
    managerid INT,
    dept_id INT
) COMMENT 'Employee table';

-- Salary grade table
CREATE TABLE salgrade (
    grade INT,
    losal INT,
    hisal INT
) COMMENT 'Salary grade table';

-- Insert practice data
INSERT INTO dept (name) VALUES ('R&D'), ('Marketing'), ('Finance'), ('Sales');

INSERT INTO emp (name, age, job, salary, entrydate, managerid, dept_id) VALUES
('Ding Minjun', 28, 'DBA', 12000, '2012-03-01', 2, 1),
('Zhao Min', 35, 'Project Manager', 18000, '2008-08-08', NULL, 1),
('Employee A', 45, 'Java Developer', 14000, '2005-01-01', 2, 1),
('Employee B', 52, 'Marketing Specialist', 8000, '2010-06-01', 4, 2),
('Employee C', 30, 'Finance', 9000, '2011-09-01', 4, 3),
('Leader B', 48, 'Marketing Director', 15000, '2006-07-01', NULL, 2);

INSERT INTO salgrade VALUES (1, 0, 3000), (2, 3001, 5000), (3, 5001, 8000),
(4, 8001, 10000), (5, 10001, 15000), (6, 15001, 20000),
(7, 20001, 25000), (8, 25001, 30000);
```

---

## Multi-table Join Queries (JOIN)

When the required data is distributed across multiple tables, you need to use JOIN to connect these tables according to specified association conditions.

### Inner Join (INNER JOIN)

**Definition**: Returns rows where the join fields match in both tables.

**Syntax**:

- Implicit inner join (WHERE clause):
  ```sql
  SELECT a.col, b.col FROM tableA a, tableB b WHERE a.id = b.id;
  ```
- Explicit inner join (recommended):
  ```sql
  SELECT a.col, b.col FROM tableA a INNER JOIN tableB b ON a.id = b.id;
  ```

**Example**: Query employee name, age, job, and department information.
```sql
-- Implicit inner join
SELECT emp.name, emp.age, emp.job, dept.name
FROM emp, dept
WHERE dept.id = emp.dept_id;

-- Explicit inner join
SELECT emp.name, emp.age, emp.job, dept.name
FROM emp JOIN dept ON emp.dept_id = dept.id;
```

---

### Left Outer Join (LEFT JOIN)

**Definition**: Returns all records from the left table and the matched records from the right table. If there is no match, the result is NULL.

**Syntax**:
```sql
SELECT a.col, b.col FROM tableA a LEFT JOIN tableB b ON a.id = b.id;
```

**Example**: Query all students and their course selection (students who have not selected courses are also displayed).
```sql
SELECT s.id, s.name, c.cour
FROM student s
LEFT JOIN student_course sc ON s.id = sc.id_stu
LEFT JOIN course c ON sc.id_cour = c.id
ORDER BY s.id, c.id;
```

---

### Right Outer Join (RIGHT JOIN)

**Definition**: Opposite of LEFT JOIN, returns all records from the right table and the matched records from the left table. If there is no match, the result is NULL.

**Syntax**:
```sql
SELECT a.col, b.col FROM tableA a RIGHT JOIN tableB b ON a.id = b.id;
```

**Example**: Query all courses and their selection status (courses not selected by any student are also displayed).
```sql
SELECT s.name, c.cour
FROM student s
RIGHT JOIN student_course sc ON s.id = sc.id_stu
RIGHT JOIN course c ON sc.id_cour = c.id;
```

---

### Full Outer Join (FULL OUTER JOIN)

**Definition**: Returns all records from both left and right tables. If a row does not have a match in the other table, the columns from the other table are NULL.

**Note**: MySQL does not directly support FULL OUTER JOIN; it can be simulated using UNION of LEFT JOIN and RIGHT JOIN.

**Simulated Syntax**:
```sql
SELECT * FROM tableA a LEFT JOIN tableB b ON a.id = b.id
UNION
SELECT * FROM tableA a RIGHT JOIN tableB b ON a.id = b.id;
```

---

### Self Join (SELF JOIN)

**Definition**: A table joins with itself. Commonly used for hierarchical data (e.g., employee and manager).

**Example**: Query employees and their managers in the emp table.
```sql
-- Only show employees with managers (inner join)
SELECT a.id AS 'Employee ID', a.name AS 'Employee Name', b.name AS 'Manager Name'
FROM emp a
JOIN emp b ON a.managerid = b.id;

-- Show all employees, manager name is NULL if no manager (left join)
SELECT a.id AS 'Employee ID', a.name AS 'Employee Name', b.name AS 'Manager Name'
FROM emp a
LEFT JOIN emp b ON a.managerid = b.id;
```

---

## Union Queries (UNION)

**Definition**: Used to combine the result sets of two or more SELECT statements.

**Rules**:

- All SELECT statements must have the same number of columns.
- Columns must have similar data types.
- The order of columns in each SELECT statement must be the same.

**Syntax**:

- `UNION`: Combines result sets and removes duplicate rows automatically.
- `UNION ALL`: Combines result sets and keeps all rows, including duplicates; more efficient.

**Example**: Query employees with salary less than 10,000 and employees older than 50.
```sql
-- UNION ALL: includes duplicates
SELECT * FROM emp WHERE salary < 10000
UNION ALL
SELECT * FROM emp WHERE age > 50;

-- UNION: removes duplicates
SELECT * FROM emp WHERE salary < 10000
UNION
SELECT * FROM emp WHERE age > 50;
```

---

## Subqueries/Nested Queries (Subquery)

**Definition**: A query nested within another SQL statement (such as SELECT, INSERT, UPDATE, DELETE).

### Scalar Subquery

**Definition**: The subquery returns a single value (one row, one column).

**Example**: Query all employees in the "R&D" department.
```sql
SELECT * FROM emp WHERE dept_id = (SELECT id FROM dept WHERE name = 'R&D');
```

**Example 2**: Query employees who joined after "Ding Minjun".
```sql
SELECT * FROM emp
WHERE entrydate > (SELECT entrydate FROM emp WHERE name = 'Ding Minjun');
```

---

### Column Subquery

**Definition**: The subquery returns one column with multiple rows.

**Common Operators**:

- `IN`: Equals any value in the list.
- `ANY`: Compares with any value returned by the subquery.
- `ALL`: Compares with all values returned by the subquery.

**Example 1**: Query all employees in "Marketing" and "Finance" departments.
```sql
SELECT * FROM emp
WHERE dept_id IN (SELECT id FROM dept WHERE name IN ('Marketing', 'Finance'));
```

**Example 2**: Query employees whose salary is higher than all employees in the R&D department.
```sql
SELECT * FROM emp
WHERE salary > ALL (SELECT salary FROM emp WHERE dept_id = (SELECT id FROM dept WHERE name = 'R&D'));
-- Equivalent to: salary > (SELECT MAX(salary) FROM ...)
```

**Example 3**: Query employees whose salary is higher than any employee in the R&D department.
```sql
SELECT * FROM emp
WHERE salary > ANY (SELECT salary FROM emp WHERE dept_id = (SELECT id FROM dept WHERE name = 'R&D'));
-- Equivalent to: salary > (SELECT MIN(salary) FROM ...)
```

---

### Row Subquery

**Definition**: The subquery returns one row with multiple columns.

**Example**: Query employees whose entry date and salary are the same as "Zhao Min".
```sql
SELECT * FROM emp
WHERE (salary, entrydate) = (SELECT salary, entrydate FROM emp WHERE name = 'Zhao Min');
```

---

### Table Subquery

**Definition**: The subquery returns multiple rows and columns, i.e., a virtual table.

**Key Point**: The returned virtual table is usually placed in the FROM clause and must be given an alias.

**Example**: Query employees who joined after 2010-06-01 and their department information.
```sql
SELECT e.*, d.name AS dept_name
FROM (SELECT * FROM emp WHERE entrydate > '2010-06-01') e
JOIN dept d ON e.dept_id = d.id;
```

---

## Comprehensive Exercises

**Exercise 1**: Query all students' course selection, showing student name, student ID, and course name.
```sql
SELECT s.name, s.id, c.cour
FROM student s
LEFT JOIN student_course sc ON s.id = sc.id_stu
LEFT JOIN course c ON sc.id_cour = c.id;
```

**Exercise 2**: Count the number of students selecting each course.
```sql
SELECT c.cour, COUNT(sc.id_stu) AS 'Number of Students'
FROM course c
LEFT JOIN student_course sc ON c.id = sc.id_cour
GROUP BY c.id, c.cour;
```

**Exercise 3**: Query the salary grade of all employees.
```sql
SELECT e.name, e.salary, s.grade
FROM emp e
JOIN salgrade s ON e.salary BETWEEN s.losal AND s.hisal;
```

**Exercise 4**: Query all employees in the "R&D" department and their salary grades.
```sql
SELECT e.*, s.grade
FROM emp e
JOIN salgrade s ON e.salary BETWEEN s.losal AND s.hisal
WHERE e.dept_id = (SELECT id FROM dept WHERE name = 'R&D');
```

**Exercise 5**: Query employees whose salary is higher than "Ding Minjun".
```sql
SELECT * FROM emp WHERE salary > (SELECT salary FROM emp WHERE name = 'Ding Minjun');
```

**Exercise 6**: Query employees whose salary is lower than the average salary of their department.
```sql
SELECT *
FROM emp e1
WHERE e1.salary < (SELECT AVG(e2.salary) FROM emp e2 WHERE e2.dept_id = e1.dept_id);
```

**Exercise 7**: Query all department information and count the number of employees in each department.
```sql
-- Using scalar subquery
SELECT
    d.id,
    d.name,
    (SELECT COUNT(*) FROM emp e WHERE e.dept_id = d.id) AS 'Number of Employees'
FROM dept d;

-- Using JOIN and GROUP BY (usually more efficient)
SELECT d.id, d.name, COUNT(e.id) AS 'Number of Employees'
FROM dept d
LEFT JOIN emp e ON d.id = e.dept_id
GROUP BY d.id, d.name;
```

---