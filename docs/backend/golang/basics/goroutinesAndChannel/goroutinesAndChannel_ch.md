# Go Concurrency Study Notes — Goroutine, Channel, and Account Management Example

本笔记总结了 Go 并发编程中的 goroutine、channel、同步机制及账户管理模型，适合作为 GitHub README 参考。

---

## 1. Goroutine Execution

- goroutine 是并发执行的函数调用，通过 `go func()` 或 `go someFunc()` 启动，后台异步运行，不阻塞调用者。
- Go 的调度器负责在多个系统线程上调度 goroutine，执行顺序不可预测。
- 若需保证代码执行完毕再继续，必须使用同步手段。

**常见同步手段：**
- **channel 信号**（最常用）：通过 channel 发送/接收数据实现同步。
- **sync.WaitGroup**：统计 goroutine 数量，等待全部完成。
- **sync.Mutex**：用于共享资源保护。

---

## 2. Channel 的作用

- channel 是 goroutine 之间的通信机制，`<-` 用于发送或接收数据。
- **无缓冲 channel**：发送方/接收方会阻塞，直到另一方操作，天然实现同步点。

---

## 3. done channel 的意义

- 在存款/取款场景，`done` 或 `Result` channel 仅作为信号，数据内容无意义，只表示“操作完成”。
- 查询余额时，返回 channel 才有实际数值。

**示例：**
```go
done := make(chan struct{})
reqChan <- accountRequest{..., Result: done}
<-done // 等待操作完成
```

---

## 4. Goroutine Leak Risk

- 若请求发出后无人接收 channel，goroutine 会永远阻塞，造成泄漏。
- 必须确保 channel 被正确接收。

---

## 5. done channel vs WaitGroup

- **done channel**：用于通知单个操作完成，精细同步。
- **WaitGroup**：用于等待一组 goroutine 完成，适合批量任务。

**对比：**
- done 更细粒度，单操作完成即通知。
- WaitGroup 更宏观，等待全部任务完成。

---

## 6. Happens-Before 原理

- Happens-Before 描述并发执行顺序和可见性。
- 在 Go 中，channel 的发送-接收、mutex 的 Lock/Unlock 都隐式形成 Happens-Before。
- 保证：channel 接收方一定在发送方操作完成后继续。

---

## 7. 查询余额的关键点

**示例：**
```go
func getBalance(name string) float64 {
    bal := make(chan float64)
    reqChan <- accountRequest{name: name, Type: "query", balance: bal}
    return <-bal
}
```
- 查询请求通过 channel 发送，后台 goroutine 处理后写入结果，主 goroutine 阻塞直到收到余额，保证同步。

---

## 8. Goroutine 与 Channel 的配合

- 主 goroutine 发送请求（存款/取款/查询）。
- 后台 accountManager goroutine 顺序处理请求，保证线程安全。
- channel 作为通信桥梁，reqChan 传递请求，done/Result channel 发信号，bal channel 返回查询结果。

**设计思想：**
- channel 发送信号应在所有逻辑执行完后，确保同步。

---

## 9. 正确的账户操作模型

- 并发安全：单线程 goroutine 管理 Account map，避免并发写入。
- 多账户并发：可同时创建不同账号，互不干扰。
- 同步机制：
    - 存取款：只需信号表示完成。
    - 查询余额：需返回实际数值。

**本质：Actor 模型**  
单个 goroutine 独立管理资源，外部通过消息（channel）请求，保证线程安全。

---

## 10. 关键知识点总结

- goroutine 执行不可预测，需用 channel/WaitGroup 保证顺序或完成。
- channel 阻塞，天然同步点。
- done channel 仅传递“完成”信号，值无实际意义。
- 存取款只需信号，查询需返回数据。
- 推荐用 goroutine 独占管理共享数据，用 channel 消息交互，保证线程安