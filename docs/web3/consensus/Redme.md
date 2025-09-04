# 分布式事务一致性：2PC、3PC 与拜占庭将军问题

## 一、2PC（Two-Phase Commit，两阶段提交）

### 1. 基本流程

```mermaid
sequenceDiagram
    participant C as 协调者
    participant A as 参与者A
    participant B as 参与者B

    C->>A: CanCommit?
    C->>B: CanCommit?
    A-->>C: Yes/No
    B-->>C: Yes/No

    alt 所有回复Yes
        C->>A: DoCommit
        C->>B: DoCommit
        A-->>C: Ack
        B-->>C: Ack
    else 有No或超时
        C->>A: Abort
        C->>B: Abort
    end
    
2. 特点
✅ 保证 原子性（要么全部提交，要么全部回滚）

❌ 缺陷：无法解决网络分区导致的阻塞问题，协调者存在单点故障

3. 脑裂问题举例
场景：协调者 C，参与者 A、B

A、B 都返回 Vote-commit

C 在发出 Commit 前网络分区

A 收到 Commit → 提交

B 未收到 → 阻塞等待

结果：系统出现 不一致

二、3PC（Three-Phase Commit，三阶段提交）
1. 改进点
在 2PC 基础上增加 PreCommit 阶段，并引入 超时机制，减少阻塞。

mermaid
复制代码
sequenceDiagram
    participant C as 协调者
    participant A as 参与者A
    participant B as 参与者B

    C->>A: CanCommit?
    C->>B: CanCommit?
    A-->>C: Yes/No
    B-->>C: Yes/No

    alt 所有Yes
        C->>A: PreCommit
        C->>B: PreCommit
        A-->>C: Yes + 写日志
        B-->>C: Yes + 写日志

        alt 协调者正常
            C->>A: DoCommit
            C->>B: DoCommit
            A-->>C: Ack
            B-->>C: Ack
        else 协调者崩溃/超时
            A->>A: 超时 → 自主提交
            B->>B: 超时 → 自主提交
        end
    else 有No或超时
        C->>A: Abort
        C->>B: Abort
    end
2. 特点
✅ 引入超时机制，参与者可在协调者失联时自主决策

✅ 相比 2PC，更少阻塞

❌ 在网络分区/拜占庭错误下，仍可能不一致

❌ 增加一轮通信，性能开销更大

3. 失效场景举例
协调者 C，参与者 A、B、X：

所有人回复 Yes

C 发送 PreCommit，A、B 收到，X 因网络延迟没收到

C 崩溃

A、B：超时后直接提交

X：未收到 PreCommit，选择回滚
→ 系统不一致

三、拜占庭将军问题
1. 问题定义
在分布式系统中，某些节点可能是恶意的（叛徒），会发送矛盾信息。如何保证忠诚节点仍能达成一致，是该问题的核心。

2. 示例（9 位将军，1 位叛徒）
4 位忠诚将军投票 进攻

4 位忠诚将军投票 撤退

1 位叛徒 → 向进攻派说自己投进攻，向撤退派说自己投撤退

结果：

进攻派看到 5 票进攻 → 选择进攻

撤退派看到 5 票撤退 → 选择撤退

系统一致性被破坏

3. 与 2PC / 3PC 的关系
2PC：保证原子性，但阻塞 & 单点风险大

3PC：减少阻塞，但无法解决拜占庭问题

拜占庭场景：需要更强一致性算法（Paxos、Raft、PBFT）

四、总结
2PC：保证原子性，但可能阻塞

3PC：通过预提交 + 超时机制减少阻塞，但仍无法应对拜占庭问题

拜占庭问题：揭示了 2PC/3PC 的局限性 → 需要更强的共识协议