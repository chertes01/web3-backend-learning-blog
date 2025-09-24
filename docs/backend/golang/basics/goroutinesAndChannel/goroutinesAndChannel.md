# Go Concurrency Study Notes — Goroutine, Channel, and Account Management Example

This note summarizes Go concurrent programming concepts including goroutines, channels, synchronization mechanisms, and an account management model. Suitable as a GitHub README reference.

---

## 1. Goroutine Execution

- A goroutine is a function call executed concurrently, started with `go func()` or `go someFunc()`. It runs asynchronously in the background and does not block the caller.
- Go's scheduler manages goroutines across multiple system threads; execution order is unpredictable.
- To ensure code completes before proceeding, synchronization is required.

**Common synchronization methods:**
- **Channel signals** (most common): Synchronize by sending/receiving data through channels.
- **sync.WaitGroup**: Track the number of goroutines and wait for all to finish.
- **sync.Mutex**: Protect shared resources.

---

## 2. Purpose of Channels

- Channels are the communication mechanism between goroutines; `<-` is used to send or receive data.
- **Unbuffered channels**: Sender/receiver will block until the other side operates, naturally creating a synchronization point.

---

## 3. Meaning of the done channel

- In deposit/withdrawal scenarios, the `done` or `Result` channel is just a signal; its data content is meaningless and only indicates "operation completed."
- For balance queries, the returned channel contains the actual value.

**Example:**
```go
done := make(chan struct{})
reqChan <- accountRequest{..., Result: done}
<-done // Wait for operation to complete
```

---

## 4. Goroutine Leak Risk

- If a request is sent but no one receives from the channel, the goroutine will block forever, causing a leak.
- Always ensure the channel is properly received.

---

## 5. done channel vs WaitGroup

- **done channel**: Notifies completion of a single operation, fine-grained synchronization.
- **WaitGroup**: Waits for a group of goroutines to finish, suitable for batch tasks.

**Comparison:**
- done is more granular, notifies as soon as a single operation completes.
- WaitGroup is more macro, waits for all tasks to finish.

---

## 6. Happens-Before Principle

- Happens-Before describes the order and visibility of concurrent execution.
- In Go, channel send/receive and mutex Lock/Unlock implicitly form Happens-Before relationships.
- Guarantee: The channel receiver always continues after the sender's operation is complete.

---

## 7. Key Points for Balance Query

**Example:**
```go
func getBalance(name string) float64 {
    bal := make(chan float64)
    reqChan <- accountRequest{name: name, Type: "query", balance: bal}
    return <-bal
}
```
- The query request is sent via channel; the background goroutine processes and writes the result. The main goroutine blocks until the balance is received, ensuring synchronization.

---

## 8. Cooperation of Goroutine and Channel

- The main goroutine sends requests (deposit/withdraw/query).
- The background accountManager goroutine processes requests sequentially, ensuring thread safety.
- Channels act as communication bridges: `reqChan` passes requests, `done`/`Result` channels send signals, `bal` channel returns query results.

**Design concept:**
- Channel signals should be sent after all logic is complete to ensure synchronization.

---

## 9. Correct Account Operation Model

- Concurrency safety: A single goroutine manages the Account map, avoiding concurrent writes.
- Multi-account concurrency: Multiple accounts can be created simultaneously without interference.
- Synchronization mechanism:
    - Deposit/withdrawal: Only need a signal to indicate completion.
    - Balance query: Must return the actual value.

**Essentially: Actor Model**  
A single goroutine manages resources independently; external requests are made via messages (channels), ensuring thread safety.

---

## 10. Key Knowledge Summary

- Goroutine execution is unpredictable; use channels/WaitGroup to ensure order or completion.
- Channel blocking is a natural synchronization point.
- done channel only transmits a "done" signal; its value is meaningless.
- Deposit/withdrawal only needs a signal; queries need to return data.
- It is recommended to use a dedicated goroutine to manage shared data and interact via channel messages to ensure thread safety.