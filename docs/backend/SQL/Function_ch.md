# 🧠 MySQL 函数学习笔记（函数篇）

---

## 一、字符串函数（String Functions）

字符串函数用于处理文本数据，比如拼接、截取、填充、大小写转换等。

| 函数 | 说明 | 示例 | 结果 |
|------|------|------|------|
| CONCAT(str1, str2, ...) | 拼接字符串 | SELECT CONCAT('Hello','World'); | HelloWorld |
| LOWER(str) | 转换为小写 | SELECT LOWER('HELLO MYSQL'); | hello mysql |
| UPPER(str) | 转换为大写 | SELECT UPPER('hello mysql'); | HELLO MYSQL |
| LPAD(str, n, padstr) | 左侧填充至指定长度 | SELECT LPAD('567',5,'-'); | --567 |
| RPAD(str, n, padstr) | 右侧填充至指定长度 | SELECT RPAD('567',5,'_'); | 567__ |
| TRIM(str) | 删除字符串开头和结尾的空格 | SELECT TRIM(' Guten Tag! '); | Guten Tag! |
| SUBSTRING(str, pos, len) | 从字符串中截取子串 | SELECT SUBSTRING('hello mysql',6,6); | mysql |

💡 **易错点：**
- LPAD() 和 RPAD() 的返回值是字符串类型，如果用于数字字段（如 id），需先 ALTER 将列改为 VARCHAR 类型。
- MySQL 中 SUBSTR() 是 SUBSTRING() 的简写，二者等价。

📘 **实际应用示例：**
```sql
ALTER TABLE employee MODIFY id VARCHAR(10);
UPDATE employee SET id = LPAD(id, 5, '0');
```
将员工编号左补0至5位，例如 1 → 00001。

📘 **补充函数：**

| 函数 | 作用 | 示例 |
|------|------|------|
| REPLACE(str, from_str, to_str) | 替换字符串 | SELECT REPLACE('good morning','morning','night'); → good night |
| LEFT(str, n) / RIGHT(str, n) | 从左/右取前 n 个字符 | SELECT LEFT('abcdef',3); → abc |

---

## 二、数值函数（Numeric Functions）

数值函数用于处理数字计算、取整、取模、随机数等。

| 函数 | 说明 | 示例 | 结果 |
|------|------|------|------|
| CEIL(x) / CEILING(x) | 向上取整 | SELECT CEIL(3.14); | 4 |
| FLOOR(x) | 向下取整 | SELECT FLOOR(3.14); | 3 |
| MOD(x,y) | 取模（x 除以 y 的余数） | SELECT MOD(5,3); | 2 |
| RAND() | 生成 [0,1) 随机数 | SELECT RAND(); | 0.6485 |
| ROUND(x, d) | 四舍五入保留 d 位小数 | SELECT ROUND(3.1415926,5); | 3.14159 |

📘 **生成六位数验证码**
```sql
SELECT LPAD(FLOOR(RAND()*100000), 6, '0') AS '随机验证码';
```
注意：RAND() 每次执行都会生成不同结果。

📘 **补充函数：**

| 函数 | 说明 | 示例 |
|------|------|------|
| TRUNCATE(x, d) | 截断小数点后 d 位（不四舍五入） | SELECT TRUNCATE(3.14159,2); → 3.14 |
| ABS(x) | 取绝对值 | SELECT ABS(-5); → 5 |
| POW(x, y) / POWER(x, y) | 幂运算 | SELECT POW(2,3); → 8 |

---

## 三、日期与时间函数（Date & Time Functions）

日期函数可获取当前时间、日期差值、提取年月日等。

| 函数 | 说明 | 示例 | 结果示例 |
|------|------|------|----------|
| CURDATE() | 当前日期 | SELECT CURDATE(); | 2025-10-06 |
| CURTIME() | 当前时间 | SELECT CURTIME(); | 12:35:20 |
| NOW() | 当前日期 + 时间 | SELECT NOW(); | 2025-10-06 12:35:20 |
| YEAR(date) | 提取年份 | SELECT YEAR('2025-10-02'); | 2025 |
| MONTH(date) | 提取月份 | SELECT MONTH('2025-10-02'); | 10 |
| DAY(date) | 提取日期（日） | SELECT DAY('2025-10-02'); | 2 |
| DATE_ADD(date, INTERVAL n unit) | 日期加减 | SELECT DATE_ADD(NOW(), INTERVAL 70 DAY); | 当前日期后70天 |
| DATEDIFF(date1, date2) | 日期差（date1 - date2） | SELECT DATEDIFF(NOW(),'2016-09-01'); | 天数差值 |

📘 **实战练习：计算员工入职天数**
```sql
SELECT ID, Name, Age, Gender, Entrydate, Workplace,
       DATEDIFF(NOW(), Entrydate) AS EntryDays
FROM employee
ORDER BY EntryDays DESC;
```
按入职天数从多到少排序。

📘 **补充函数：**

| 函数 | 说明 | 示例 |
|------|------|------|
| TIMESTAMPDIFF(unit, datetime1, datetime2) | 计算时间差（指定单位） | SELECT TIMESTAMPDIFF(YEAR, Entrydate, NOW()); → 入职年数 |
| DATE_FORMAT(date, format) | 按格式显示日期 | SELECT DATE_FORMAT(NOW(), '%Y年%m月%d日 %H:%i'); |

---

## 四、流程控制函数（Flow Control Functions）

### 1️⃣ IF 条件函数
```sql
SELECT IF(1>2, 'true', 'false');
-- 结果：false
```
语法：`IF(条件, 值1, 值2)`  
当条件成立返回值1，否则返回值2。

### 2️⃣ IFNULL 判断空值
```sql
SELECT IFNULL(NULL,'value');   -- 返回 value
SELECT IFNULL('', '1');        -- 空字符串不是 NULL，返回 ''
```
语法：`IFNULL(expr1, expr2)`  
当 expr1 为 NULL 时返回 expr2，否则返回 expr1。

💡 区别：NULL 与空字符串 '' 不一样！

### 3️⃣ CASE WHEN THEN ELSE END 条件表达式

📘 **示例1：判断工作城市类型**
```sql
SELECT Name,
       Workplace,
       CASE
           WHEN Workplace IN ('Beijing', 'Shanghai', 'Guangzhou', 'Shenzhen')
               THEN '一线城市'
           ELSE '二线城市'
       END AS CityType
FROM employee;
```

📘 **示例2：判断员工入职年限**
```sql
SELECT ID, Name, Age, Gender, Entrydate, Workplace,
       CASE
           WHEN TIMESTAMPDIFF(YEAR, Entrydate, NOW()) >= 5 THEN '老员工'
           WHEN TIMESTAMPDIFF(YEAR, Entrydate, NOW()) >= 2 THEN '普通员工'
           ELSE '新员工'
       END AS '入职年限'
FROM employee;
```

💡 技巧：
- CASE 可以嵌套或用于 ORDER BY 实现自定义排序；
- 对日期比较建议使用 TIMESTAMPDIFF()，比 DATEDIFF()/365 更准确。

---

## ✅ 小结与拓展

| 函数类型   | 常用函数 | 备注 |
|------------|----------|------|
| 字符串函数 | CONCAT, LOWER, UPPER, LPAD, RPAD, TRIM, SUBSTRING, REPLACE, LEFT, RIGHT | 处理文本拼接、截取、填充等 |
| 数值函数   | CEIL, FLOOR, MOD, RAND, ROUND, ABS, TRUNCATE, POW | 常用于数值计算与取整 |
| 日期函数   | NOW, CURDATE, DATE_ADD, DATEDIFF, TIMESTAMPDIFF, DATE_FORMAT | 时间运算和格式化 |
| 流程控制   | IF, IFNULL, CASE WHEN | 逻辑判断、数据分级 |

⚠️ **易错与补充说明**
- LPAD / RPAD 会将数字转为字符串类型；
- 日期差不要直接 /365，使用 TIMESTAMPDIFF(YEAR, date1, date2) 更精确；
- IFNULL() 仅判断 NULL，不会判断空字符串；
- CASE 语句末尾必须有 END；
- 日期函数结果受系统时区影响，可用 `@@global.time_zone` 查看。

---