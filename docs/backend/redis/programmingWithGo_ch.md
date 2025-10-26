# Redis 编程（Go语言版）

## 一、概述

Redis 编程指的是通过编程语言操作 Redis 数据库。Redis 官方为多种语言提供了客户端支持，详情请见 [Redis 官方客户端列表](https://redis.io/resources/clients/)。

在 Go 语言生态中，常用的客户端主要有以下几种：

| 客户端         | 特点                                                     |
| :------------- | :------------------------------------------------------- |
| **go-redis**   | 官方推荐、功能全面、支持哨兵、集群、Pipeline、事务     |
| **redigo**     | 轻量、简单，适合基础场景                                 |
| **goredis/v9** | `go-redis` 最新版本，支持 Context API、泛型等现代 Go 特性 |

在 Go 生态中，`go-redis` 是使用最广、功能最全的 Redis 客户端，因此以下笔记将以 **`go-redis/v9`** 为主进行讲解。

---

## 二、Go-Redis 使用入门

### 1. 初始化项目与安装依赖

首先，新建一个 Go 模块项目：

```bash
mkdir go-redis-demo
cd go-redis-demo
go mod init go-redis-demo
```

接着，安装 `go-redis`：

```bash
go get github.com/redis/go-redis/v9
```

### 2. Redis 基本配置与连接

创建文件 `redis_client.go`：

```go
package main

import (
    "context"
    "fmt"
    "github.com/redis/go-redis/v9"
    "log"
)

var (
    ctx = context.Background()
    Rdb *redis.Client
)

// InitRedis 初始化 Redis 客户端
func InitRedis() {
    Rdb = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis 地址
        Password: "",               // 密码（无则留空）
        DB:       0,                // 默认数据库
        PoolSize: 10,               // 连接池大小
    })

    // 测试连接
    _, err := Rdb.Ping(ctx).Result()
    if err != nil {
        log.Fatalf("连接 Redis 失败: %v", err)
    }
    fmt.Println("✅ 已成功连接 Redis！")
}
```

### 3. 基本 CRUD 操作示例

创建文件 `main.go`：

```go
package main

import (
    "context"
    "fmt"
)

func main() {
    InitRedis()
    defer Rdb.Close()
    ctx := context.Background()

    // 写入键值对 name:dafei
    err := Rdb.Set(ctx, "name", "dafei", 0).Err()
    if err != nil {
        panic(err)
    }

    // 读取键值
    val, err := Rdb.Get(ctx, "name").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("name =", val)
}
```

运行后，输出结果：

```
✅ 已成功连接 Redis！
name = dafei
```

---

## 三、Go-Redis 操作进阶

`go-redis` 的 API 设计与 Redis 原生命令一一对应，调用形式非常直观：

- `Rdb.Get(ctx, key)`
- `Rdb.Set(ctx, key, value, expiration)`
- `Rdb.Del(ctx, key)`
- `Rdb.Exists(ctx, key)`
- `Rdb.Keys(ctx, pattern)`

以下是常见数据结构的操作对照表：

| Redis 类型 | Go-Redis 方法示例                                           | 对应命令        |
| :--------- | :---------------------------------------------------------- | :-------------- |
| **String** | `Rdb.Set(ctx, "k", "v", 0)` / `Rdb.Get(ctx, "k")`            | `SET` / `GET`   |
| **Hash**   | `Rdb.HSet(ctx, "user:1", "name", "von")` / `Rdb.HGetAll(ctx, "user:1")` | `HSET` / `HGETALL` |
| **List**   | `Rdb.LPush(ctx, "mylist", "a", "b")` / `Rdb.LRange(ctx, "mylist", 0, -1)` | `LPUSH` / `LRANGE` |
| **Set**    | `Rdb.SAdd(ctx, "myset", "x", "y")` / `Rdb.SMembers(ctx, "myset")` | `SADD` / `SMEMBERS` |
| **ZSet**   | `Rdb.ZAdd(ctx, "myzset", redis.Z{Score: 1, Member: "tom"})` | `ZADD`          |
| **Keys**   | `Rdb.Keys(ctx, "*")`                                        | `KEYS *`        |
| **事务**   | `Rdb.TxPipeline()`                                          | `MULTI` / `EXEC` |

---

## 四、连接池与配置优化

`go-redis` 内部默认实现了连接池，可以通过 `redis.Options` 进行详细配置：

```go
Rdb = redis.NewClient(&redis.Options{
    Addr:         "localhost:6379",
    Password:     "",
    DB:           0,
    PoolSize:     100,   // 最大连接数
    MinIdleConns: 10,    // 最小空闲连接数
    ReadTimeout:  -1,    // 读取超时（-1 表示不设置超时）
    WriteTimeout: -1,    // 写入超时
    IdleTimeout:  300,   // 空闲连接的超时时间（秒）
})
```

> **注意**：如果连接不稳定或长时间空闲，可通过 `Rdb.Ping()` 定期检测连接状态。

---

## 五、错误处理与资源释放

在生产代码中，应始终对 Redis 操作返回的 `err` 进行检查。特别是查询一个不存在的键时，`go-redis` 会返回 `redis.Nil` 错误。

```go
val, err := Rdb.Get(ctx, "nonexistent_key").Result()
if err == redis.Nil {
    fmt.Println("key 不存在")
} else if err != nil {
    fmt.Println("发生其他错误:", err)
} else {
    fmt.Println("value =", val)
}
```

程序退出前，应确保调用 `Close()` 方法来释放 Redis 连接资源：

```go
defer Rdb.Close()
```

---

## 六、实战示例：缓存应用

使用 Redis 作为缓存层是其最常见的应用场景之一。以下函数演示了如何从缓存中获取用户名，如果缓存未命中，则从数据库中读取并回填到缓存。

```go
import "time"

func GetUserName(uid string) string {
    key := "user:name:" + uid
    // 1. 从 Redis 读取缓存
    val, err := Rdb.Get(ctx, key).Result()

    // 2. 判断错误类型
    if err == redis.Nil {
        fmt.Println("缓存未命中，从数据库读取...")
        // 假设从数据库查到用户名为 "von"
        name := "von" 
        
        // 3. 将数据写入缓存，并设置1小时过期
        Rdb.Set(ctx, key, name, 1*time.Hour)
        return name
    } else if err != nil {
        log.Println("Redis 错误:", err)
        return "" // 或者返回默认值、处理错误
    }

    // 4. 缓存命中，直接返回
    return val
}
```

---

## 七、与 Java 客户端对照总结

| 功能         | Java (Jedis/Lettuce) | Go (go-redis)             |
| :----------- | :------------------- | :------------------------ |
| **基础 CRUD**  | `jedis.set("k", "v")`| `Rdb.Set(ctx, "k", "v", 0)` |
| **连接池配置** | `JedisPoolConfig`    | `redis.Options`           |
| **Spring 集成**| `StringRedisTemplate`| 原生 API，无需额外封装    |
| **事务**       | `multi` / `exec`     | `TxPipeline()`            |
| **分布式锁**   | Redisson             | `SetNX()` + Lua 脚本实现  |
| **发布/订阅**  | `jedis.subscribe()`  | `Rdb.Subscribe()`         |

---

## 八、生产级配置示例

在实际项目中，通常会将配置信息、Redis 客户端初始化和业务逻辑分离开来，以提高代码的可维护性。

### 1️⃣ 配置文件：`config.yaml`

```yaml
redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
  poolSize: 50
  minIdleConns: 10
  readTimeout: 5   # 秒
  writeTimeout: 5  # 秒
  idleTimeout: 300 # 秒
```

### 2️⃣ 配置加载：`config/config.go`

此模块负责读取 `config.yaml` 并支持通过环境变量覆盖配置。

```go
package config

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "gopkg.in/yaml.v2"
)

type RedisConfig struct {
    Host         string `yaml:"host"`
    Port         int    `yaml:"port"`
    Password     string `yaml:"password"`
    DB           int    `yaml:"db"`
    PoolSize     int    `yaml:"poolSize"`
    MinIdleConns int    `yaml:"minIdleConns"`
    ReadTimeout  int    `yaml:"readTimeout"`
    WriteTimeout int    `yaml:"writeTimeout"`
    IdleTimeout  int    `yaml:"idleTimeout"`
}

type Config struct {
    Redis RedisConfig `yaml:"redis"`
}

var Conf Config

func LoadConfig(path string) {
    data, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatalf("读取配置文件失败: %v", err)
    }
    err = yaml.Unmarshal(data, &Conf)
    if err != nil {
        log.Fatalf("解析配置文件失败: %v", err)
    }

    // 环境变量覆盖，方便多环境部署
    if host := os.Getenv("REDIS_HOST"); host != "" {
        Conf.Redis.Host = host
    }
    if port := os.Getenv("REDIS_PORT"); port != "" {
        fmt.Sscanf(port, "%d", &Conf.Redis.Port)
    }
    if pwd := os.Getenv("REDIS_PASSWORD"); pwd != "" {
        Conf.Redis.Password = pwd
    }
}
```

### 3️⃣ Redis 客户端初始化：`redisclient/redis_client.go`

创建一个全局可用的 Redis 客户端实例。

```go
package redisclient

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/redis/go-redis/v9"
    "your_project/config" // 替换为你的项目路径
)

var (
    Rdb *redis.Client
    Ctx = context.Background()
)

func InitRedis() {
    c := config.Conf.Redis
    Rdb = redis.NewClient(&redis.Options{
        Addr:         fmt.Sprintf("%s:%d", c.Host, c.Port),
        Password:     c.Password,
        DB:           c.DB,
        PoolSize:     c.PoolSize,
        MinIdleConns: c.MinIdleConns,
        ReadTimeout:  time.Duration(c.ReadTimeout) * time.Second,
        WriteTimeout: time.Duration(c.WriteTimeout) * time.Second,
        IdleTimeout:  time.Duration(c.IdleTimeout) * time.Second,
    })

    // 测试连接
    if _, err := Rdb.Ping(Ctx).Result(); err != nil {
        log.Fatalf("连接 Redis 失败: %v", err)
    }
    log.Println("✅ Redis 初始化成功")
}
```

### 4️⃣ 使用示例：`main.go`

```go
package main

import (
    "fmt"
    "your_project/config"       // 替换为你的项目路径
    "your_project/redisclient"  // 替换为你的项目路径
)

func main() {
    // 1. 加载配置
    config.LoadConfig("config.yaml")

    // 2. 初始化 Redis
    redisclient.InitRedis()
    defer redisclient.Rdb.Close()

    // 3. 基本 CRUD 操作
    ctx := redisclient.Ctx
    err := redisclient.Rdb.Set(ctx, "name", "von", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := redisclient.Rdb.Get(ctx, "name").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("name =", val)

    // 4. 示例：删除
    redisclient.Rdb.Del(ctx, "name")
}
```

---

## 九、总结

- **客户端选择**：`go-redis` 是 Go 生态中最强大、最主流的 Redis 客户端。
- **API 设计**：其 API 与 Redis 原生命令高度一致，学习曲线平缓。
- **生产实践**：
    - **配置分离**：使用 YAML 或其他格式管理 Redis 参数，并通过环境变量覆盖，以适应不同环境（开发、测试、生产）。
    - **连接池**：合理配置 `PoolSize`、`MinIdleConns` 和 `IdleTimeout` 等参数，以提升性能和稳定性。
    - **缓存策略**：结合业务场景，设计合理的缓存过期策略和错误处理机制。
    - **统一客户端**：初始化一个全局的 Redis 客户端实例，方便在项目各处调用。
    - **健康检查**：在初始化时自动检测 Redis 连接，确保服务启动时依赖就绪。
