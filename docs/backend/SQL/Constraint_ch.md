# 🧭 MySQL 约束学习笔记（Constraint Notes）

约束（Constraint）用于限制表中数据的规则，保证数据的正确性、一致性、完整性。  
常见约束类型包括：

- PRIMARY KEY 主键约束
- AUTO_INCREMENT 自动增长
- NOT NULL 非空约束
- UNIQUE 唯一约束
- DEFAULT 默认值约束
- CHECK 检查约束
- FOREIGN KEY 外键约束

---

## 🧩 一、基本约束类型讲解与示例

### 1️⃣ 主键约束（PRIMARY KEY）

**作用**：唯一标识一条记录，不允许重复，也不允许为 NULL。  
**关键字**：PRIMARY KEY

- 一个表只能有一个主键。
- 主键列值必须唯一且非空。
- 可搭配 AUTO_INCREMENT 实现自动编号。

```sql
CREATE TABLE user(
    id INT PRIMARY KEY AUTO_INCREMENT,
    ...
);
```

> auto_increment 只能用于整数类型。  
> 删除表时主键会自动删除，不需单独解除约束。

---

### 2️⃣ 非空约束（NOT NULL）

**作用**：保证该字段不能为空。  
**关键字**：NOT NULL

```sql
name VARCHAR(10) NOT NULL
```

> 不能在插入时省略该列的值，否则报错。  
> 如果需要可以为空，应去掉 NOT NULL。

---

### 3️⃣ 唯一约束（UNIQUE）

**作用**：保证该字段的值在表中唯一。  
**关键字**：UNIQUE

```sql
name VARCHAR(10) NOT NULL UNIQUE
```

> 与主键不同，唯一约束可以有多个列。  
> 允许存在 NULL（NULL 不参与比较）。

---

### 4️⃣ 默认约束（DEFAULT）

**作用**：插入数据时未指定该列值，自动使用默认值。  
**关键字**：DEFAULT

```sql
status CHAR(1) DEFAULT '1'
```

> 插入时省略该字段会自动填入默认值。  
> 若明确写入 NULL，则不会使用默认值。

---

### 5️⃣ 检查约束（CHECK）

**作用**：限制字段值的范围或格式。  
**关键字**：CHECK (条件表达式)

```sql
age INT CHECK (age > 0 AND age <= 120)
```

> MySQL 5.x 以前 CHECK 不生效；MySQL 8.0+ 已正式支持。  
> 条件中必须是逻辑表达式。

---

### ✅ 综合示例：创建完整 user 表

```sql
CREATE TABLE user(
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(10) NOT NULL UNIQUE,
    age INT CHECK (age > 0 AND age <= 120),
    status CHAR(1) DEFAULT '1',
    gender VARCHAR(6)
);
```

**插入数据示例：**

```sql
INSERT INTO user(name,age,status,gender) VALUES ('Alice',18,1,'female');
INSERT INTO user(name,age,status,gender) VALUES ('Bob',14,1,'male');
```

---

## 🏗️ 二、外键约束（FOREIGN KEY）

### 1️⃣ 概念说明

外键用于建立两个表之间的引用关系，保证数据的参照完整性。

- 外键所在表：子表（如 emp）
- 被引用的表：父表（如 dept）
- 外键字段的值必须来自父表中已存在的主键或唯一值

---

### 2️⃣ 父表与子表示例

**创建父表：部门表 dept**

```sql
CREATE TABLE dept(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

INSERT INTO dept(name)
VALUES ('研发部'),('项目部'),('市场部'),('财务部');
```

**创建子表：员工表 emp**

```sql
CREATE TABLE emp (
    id INT AUTO_INCREMENT PRIMARY KEY COMMENT 'ID',
    name VARCHAR(50) NOT NULL COMMENT '姓名',
    age INT COMMENT '年龄',
    job VARCHAR(20) COMMENT '职位',
    salary INT COMMENT '薪资',
    entrydate DATE COMMENT '入职时间',
    managerid INT COMMENT '直属领导ID',
    dept_id INT COMMENT '部门ID'
) COMMENT '员工表';
```

---

### 3️⃣ 添加外键约束

```sql
ALTER TABLE emp
ADD CONSTRAINT fk_emp_dept_id
FOREIGN KEY (dept_id) REFERENCES dept(id);
```

- fk_emp_dept_id 是外键约束名称，可自定义。
- (dept_id) 是子表中的外键字段。
- references dept(id) 表示引用父表 dept 的主键 id。

---

### 4️⃣ 删除外键约束

```sql
ALTER TABLE emp DROP FOREIGN KEY fk_emp_dept_id;
```

---

### 5️⃣ 外键的级联操作

#### ON UPDATE CASCADE、ON DELETE CASCADE

父表主键被更新/删除时，子表对应数据自动更新或删除。

```sql
ALTER TABLE emp
ADD CONSTRAINT fk_emp_dept_id
FOREIGN KEY (dept_id)
REFERENCES dept(id)
ON UPDATE CASCADE
ON DELETE CASCADE;
```

> 用途：保证数据同步删除，例如删除部门时自动删除员工。

#### ON UPDATE SET NULL、ON DELETE SET NULL

父表主键被更新/删除时，子表外键字段自动设为 NULL。

```sql
ALTER TABLE emp
ADD CONSTRAINT fk_emp_dept_id
FOREIGN KEY (dept_id)
REFERENCES dept(id)
ON UPDATE SET NULL
ON DELETE SET NULL;
```

> 用途：保留子表记录，但取消它与父表的关联。  
> ⚠️ 子表外键字段必须允许为 NULL。

---

## ⚙️ 三、表的删除顺序与依赖关系

1️⃣ 删除子表 → 删除父表

```sql
DROP TABLE emp;
DROP TABLE dept;
```

否则会因为外键依赖关系导致错误。

2️⃣ 如果使用 ON DELETE CASCADE，删除父表记录时会自动删除子表记录。

---

## 🧠 四、执行顺序与易错点总结

| 操作类型   | 正确顺序           | 易错提示                                 |
|------------|--------------------|------------------------------------------|
| 创建外键   | 先创建父表，再创建子表 | 子表字段类型、长度必须与父表主键一致     |
| 删除外键   | 先删除外键约束，再删除父表 | 否则报“Cannot delete or update a parent row” |
| 插入数据   | 先插入父表，再插入子表 | 子表外键值必须存在于父表中               |
| 删除表     | 子表 → 父表         | 否则因外键引用报错                       |
| CHECK条件  | MySQL 8.0以前不生效 | 推荐使用MySQL 8.0+版本                   |

---

## 📚 五、补充知识点

- 查看表结构  
  ```sql
  DESC user;
  DESC emp;
  ```

- 查看所有约束  
  ```sql
  SELECT * FROM information_schema.table_constraints
  WHERE table_name = 'emp';
  ```

- 重命名外键约束  
  ```sql
  ALTER TABLE emp DROP FOREIGN KEY fk_old;
  ALTER TABLE emp ADD CONSTRAINT fk_new FOREIGN KEY (dept_id) REFERENCES dept(id);
  ```

---

## ✨ 六、学习总结

| 约束类型   | 关键字         | 作用           | 示例                                   |
|------------|---------------|----------------|----------------------------------------|
| 主键约束   | PRIMARY KEY    | 唯一标识记录   | id int primary key                     |
| 自增约束   | AUTO_INCREMENT | 自动编号       | id int auto_increment                  |
| 非空约束   | NOT NULL       | 不允许为空     | name varchar(10) not null              |
| 唯一约束   | UNIQUE         | 不允许重复     | name varchar(10) unique                |
| 默认约束   | DEFAULT        | 未指定值时使用默认 | status char(1) default '1'         |
| 检查约束   | CHECK          | 限制数值范围   | check(age>0 and age<=120)              |
| 外键约束   | FOREIGN KEY    | 建立表间关系   | foreign key(dept_id) references dept(id)|

---