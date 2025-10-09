# 🧭 MySQL DCL 学习笔记（用户与权限管理）

---

## 一、DCL 概述

DCL（Data Control Language，数据控制语言）主要用于管理数据库用户及其权限，常用命令包括：

- `CREATE USER`：创建用户
- `ALTER USER`：修改用户属性或密码
- `DROP USER`：删除用户
- `GRANT`：授予权限
- `REVOKE`：撤销权限
- `SHOW GRANTS`：查看权限

---

## 二、用户管理

### 1️⃣ 创建用户

```sql
CREATE USER '用户名'@'主机名' IDENTIFIED BY '密码';
```

**示例：**

```sql
CREATE USER 'itCast'@'localhost' IDENTIFIED BY '123456';
CREATE USER 'wind'@'%' IDENTIFIED BY '123456';
```

- `'用户名'@'主机名'` 用来区分用户。
    - `'localhost'`：仅允许本机登录
    - `'%'`：允许任意主机远程访问
    - `'192.168.1.%'`：允许指定网段访问
- 默认用户信息存放在系统数据库 `mysql.user` 表中。
- 可用 `SELECT User,Host,authentication_string FROM mysql.user;` 查看所有用户。

---

### 2️⃣ 修改用户密码

```sql
ALTER USER '用户名'@'主机名' IDENTIFIED WITH mysql_native_password BY '新密码';
```

**示例：**

```sql
ALTER USER 'wind'@'%' IDENTIFIED WITH mysql_native_password BY '1234';
```

- `WITH mysql_native_password` 指定认证插件。
- MySQL 8.0 默认使用 `caching_sha2_password`，部分客户端不兼容时可改为 `mysql_native_password`。

---

### 3️⃣ 删除用户

```sql
DROP USER '用户名'@'主机名';
```

**示例：**

```sql
DROP USER 'itCast'@'localhost';
```

- 删除后，该用户的所有权限记录会同时被清除。
- 删除前可用 `SHOW GRANTS FOR 'user'@'host';` 查看权限。

---

## 三、权限管理

### 1️⃣ 授权

```sql
GRANT 权限列表 ON 数据库名.对象名 TO '用户名'@'主机名';
```

**示例：**

```sql
GRANT ALL ON testdb.employee TO 'wind'@'%';
```

- `ALL` 表示授予所有权限（等价于 SELECT, INSERT, UPDATE, DELETE, CREATE, DROP 等）。
- 数据库名.对象名 可指定不同层级：
    - `*.*` → 全局权限
    - `testdb.*` → 某数据库权限
    - `testdb.employee` → 某表权限
    - `ON testdb.employee(column)` → 某列权限
- 授权后无需 `FLUSH PRIVILEGES`，`GRANT` 会自动刷新权限表。

---

### 2️⃣ 查看权限

```sql
SHOW GRANTS FOR '用户名'@'主机名';
```

**示例：**

```sql
SHOW GRANTS FOR 'wind'@'%';
```

**返回示例：**

```
GRANT USAGE ON *.* TO 'wind'@'%'
GRANT ALL PRIVILEGES ON `testdb`.`employee` TO 'wind'@'%'
```

---

### 3️⃣ 撤销权限

```sql
REVOKE 权限列表 ON 数据库名.对象名 FROM '用户名'@'主机名';
```

**示例：**

```sql
REVOKE ALL ON testdb.employee FROM 'wind'@'%';
```

- 撤销权限后，用户仍然存在，只是不能再访问该资源。
- 若需彻底移除用户，需使用 `DROP USER`。

---

## 四、权限生效机制

- 用户权限存储在系统数据库 `mysql` 的以下表中：
    - `user`：全局权限
    - `db`：数据库级权限
    - `tables_priv`：表级权限
    - `columns_priv`：列级权限
- `GRANT` / `REVOKE` 操作后自动刷新权限表；
- 如果直接修改这些表，必须执行：
    ```sql
    FLUSH PRIVILEGES;
    ```

---

## 五、常见权限类型速查表

| 权限名        | 作用范围      | 含义                   |
|---------------|--------------|------------------------|
| SELECT        | 表级/库级    | 读取表数据             |
| INSERT        | 表级/库级    | 插入新数据             |
| UPDATE        | 表级/库级    | 修改已有数据           |
| DELETE        | 表级/库级    | 删除数据               |
| CREATE        | 库级/全局    | 创建数据库或表         |
| DROP          | 库级/全局    | 删除数据库或表         |
| ALTER         | 表级         | 修改表结构             |
| INDEX         | 表级         | 创建或删除索引         |
| EXECUTE       | 存储过程级   | 执行存储过程           |
| GRANT OPTION  | 任意级       | 允许用户将权限再授权    |
| RELOAD        | 全局         | 执行 FLUSH 相关操作    |
| SHUTDOWN      | 全局         | 关闭 MySQL 服务        |
| SUPER         | 全局         | 管理级权限（如 KILL 查询）|

---

## 六、执行顺序建议

1. `USE mysql;` — 切换到系统数据库查看用户
2. `CREATE USER` — 创建用户
3. `ALTER USER` — 设置或修改密码
4. `GRANT` — 授权
5. `SHOW GRANTS` — 验证授权
6. `REVOKE` — 撤销权限
7. `DROP USER` — 删除用户

---

## ✅ 小结

- DCL 是数据库安全管理的核心，主要操作：用户创建 + 权限分配 + 权限撤销 + 用户删除。
- 授权、撤销时要注意权限层级（全局 / 数据库 / 表）。
- 修改权限表后需刷新或重新登录才能生效。

---