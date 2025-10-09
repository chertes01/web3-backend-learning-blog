# 🧱 MySQL 数据类型学习笔记（建表字段定义详解）

在建表时，需要为每个字段选择合适的数据类型（Data Type），这关系到数据库的存储效率、查询速度和数据精度。

---

## 一、整数类型（INT 系列）

| 类型      | 字节   | 范围（有符号）         | 适用场景           | 备注                        |
|-----------|--------|-----------------------|--------------------|-----------------------------|
| TINYINT   | 1字节  | -128 ~ 127            | 性别、状态码、小范围计数 | UNSIGNED 可达 0~255         |
| SMALLINT  | 2字节  | -32768 ~ 32767        | 年份、年龄、评分等 | 节省空间                    |
| MEDIUMINT | 3字节  | -8388608 ~ 8388607    | 不太常用           |                             |
| INT       | 4字节  | -21亿 ~ 21亿          | 最常用ID、数量等   | 默认值                      |
| BIGINT    | 8字节  | 超大整数              | 订单号、资金类、社交ID |                             |

💡 **经验技巧**  
- 能用小类型就别用大类型（节省内存）。
- 如果数值永远不会是负的，就加 UNSIGNED。

---

## 二、浮点与定点类型

| 类型         | 含义         | 特点                | 适用场景         |
|--------------|--------------|---------------------|------------------|
| FLOAT        | 单精度浮点数 | 精度有限，约7位有效数字 | 测试、粗略测量   |
| DOUBLE       | 双精度浮点数 | 精度更高，约15位    | 科学计算         |
| DECIMAL(M,D) | 定点数       | 精确小数，金融首选  | 价格、利率       |

💡 **示例**  
```sql
price DECIMAL(10,2)  -- 最大99999999.99
```

---

## 三、字符串类型（CHAR vs VARCHAR）

| 类型      | 含义         | 特点         | 适用场景         |
|-----------|--------------|--------------|------------------|
| CHAR(n)   | 固定长度字符串 | 长度固定，速度快 | 性别、国家码等定长字段 |
| VARCHAR(n)| 可变长度字符串 | 更灵活，节省空间 | 姓名、地址、备注等    |
| TEXT      | 长文本        | 不可加默认值     | 文章内容、评论        |
| ENUM      | 枚举类型      | 限定取值        | 性别、状态（male/female） |

💡 **示例**  
```sql
gender CHAR(1) CHECK (gender IN ('M','F'));
name VARCHAR(30) NOT NULL;
```

---

## 四、日期与时间类型

| 类型      | 格式                   | 说明           | 示例                      |
|-----------|------------------------|----------------|---------------------------|
| DATE      | YYYY-MM-DD             | 日期           | '2025-10-09'              |
| TIME      | HH:MM:SS               | 时间           | '18:25:30'                |
| DATETIME  | YYYY-MM-DD HH:MM:SS    | 日期+时间      | '2025-10-09 18:25:30'     |
| TIMESTAMP | 同上                   | 随时区自动转换 | 一般用于记录修改时间      |
| YEAR      | YYYY                   | 年份           | '2025'                    |

💡 **建议**  
- 自动记录插入/修改时间 → 用 TIMESTAMP DEFAULT CURRENT_TIMESTAMP
- 仅存日期 → 用 DATE

---

## 五、布尔类型（BOOLEAN）

MySQL 没有真正的布尔类型，BOOLEAN 实际是 TINYINT(1) 的别名：

```sql
is_active BOOLEAN DEFAULT 1;  -- 实际等价于 TINYINT(1)
```

---

## 六、二进制类型（了解）

| 类型        | 用途                   |
|-------------|------------------------|
| BLOB        | 存二进制数据（图片、文件） |
| VARBINARY(n)| 可变长度二进制串        |
| BINARY(n)   | 固定长度二进制串        |

---

## 🧩 实战对比示例

```sql
CREATE TABLE userInformation (
    ID INT UNSIGNED PRIMARY KEY AUTO_INCREMENT COMMENT '用户ID',
    name VARCHAR(20) NOT NULL COMMENT '姓名',
    gender ENUM('M','F') DEFAULT 'M' COMMENT '性别',
    age TINYINT UNSIGNED COMMENT '年龄',
    idCard CHAR(18) UNIQUE COMMENT '身份证号',
    salary DECIMAL(10,2) COMMENT '薪资',
    register_time DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户信息表';
```

---

## ✨ 总结口诀

| 分类       | 选用建议                        |
|------------|---------------------------------|
| 整数类型   | 够用就小，用 UNSIGNED 节省空间  |
| 浮点定点   | 金额用 DECIMAL，普通计算用 DOUBLE |
| 字符串     | 定长 CHAR，变长 VARCHAR，超长 TEXT |
| 日期时间   | 一般 DATETIME，自动更新用 TIMESTAMP |
| 逻辑值     | BOOLEAN 实为 TINYINT(1)         |
| 编码存储   | 建议统一 utf8mb4                |

---