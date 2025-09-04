# Consensus Systems Study Notes

This note summarizes the concepts of distributed system consensus, Byzantine Fault Tolerance (BFT), communication models, FLP and CAP theorems, and consensus algorithm classifications. It provides a theoretical foundation for understanding blockchain consensus mechanisms.

---

## 1. Basic Definition of Consensus Systems

In distributed systems, the goal of a consensus algorithm is to allow multiple nodes to reach agreement on system state or operations, even in the presence of faults or malicious nodes.

### 1.1 Byzantine Faults and Failures

1. **Byzantine Fault**  
   - Node behavior is variable and unpredictable, and different observers may see different outcomes.

2. **Byzantine Failure**  
   - When a Byzantine fault causes the system to lose service, e.g., data corruption, malware interference, or arbitrary malicious behavior by nodes.

3. **Crash Fault**  
   - Node process stops running but causes no other side effects.

4. **Crash-Recovery Fault**  
   - Node process stops but can resume normal operation without affecting the rest of the system.

### 1.2 Core Properties of Consensus Algorithms

- **Consistency**: All nodes agree on the same decision.  
- **Validity**: Accepted decisions must be proposed by at least one node.  
- **Termination / Liveness**: All nodes eventually complete the decision process.

> **Safety**: Satisfies consistency and validity  
> **Liveness**: Satisfies termination

---

## 2. Communication Models

Communication models define the timing and constraints of message delivery:

1. **Synchronous Model**  
   - There is a known upper bound on message delay; exceeding it is considered a fault.  
   - Ideal for algorithm design, but hard to guarantee in real networks.

2. **Asynchronous Model**  
   - No bound on message delay; messages may arrive at any time.  
   - Closer to real-world networks, but FLP impossibility applies.

3. **Partially Synchronous Model**  
   - System behaves asynchronously most of the time but occasionally synchronously.  
   - Balances theoretical analysis and practical application.

---

## 3. FLP Theorem

- **Statement**: In an asynchronous network, if even a single node can fail, there is no perfect consensus algorithm that guarantees termination (final agreement).  
- **Significance**:  
  - Highlights fundamental limitations of consensus in asynchronous networks.  
  - Modern algorithms often adopt partial synchrony to ensure safety during critical periods, even if liveness is temporarily affected.

---

## 4. CAP Theorem

- **Consistency (C)**: Every read returns the most recent write.  
- **Availability (A)**: Every request receives a response within a bounded time.  
- **Partition Tolerance (P)**: The system continues to operate despite network partitions.  

> **Core idea**: A distributed system cannot fully achieve all three simultaneously. Trade-offs are necessary.  
> Example:  
> - Dynamo prioritizes availability over strict consistency.  
> - GFS prioritizes consistency over availability.

---

## 5. Consensus Algorithm Classification

### 5.1 By Fault Tolerance

- **Byzantine Fault Tolerant (BFT)**  
  - Can tolerate malicious nodes  
  - Examples: PBFT, PoW, PoS, DPoS  

- **Non-Byzantine Fault Tolerant**  
  - Assumes honest nodes; handles crashes or network faults  
  - Examples: Paxos, Raft  

### 5.2 By Determinism

- **Deterministic Consensus**  
  - Decisions are final and cannot be reverted  
  - Examples: PBFT, Paxos, Raft  

- **Probabilistic Consensus**  
  - Decisions may be reverted, probability decreases over time  
  - Examples: PoW, some PoS variants  

### 5.3 By Leader Selection Strategy

- **Election-based**  
  - Nodes vote to elect a block producer  
  - Examples: PBFT, Raft  

- **Proof-based**  
  - Nodes win block rights via computational power or stake  
  - Examples: PoW, PoS  

---

## 6. Summary

1. Consensus in distributed systems must balance **safety** and **liveness**.  
2. Communication model and fault assumptions directly affect algorithm design.  
3. FLP and CAP theorems provide theoretical limits and design constraints.  
4. Algorithm selection depends on scenario:  
   - Public chains → BFT + probabilistic (PoW/PoS)  
   - Consortium/private chains → BFT + deterministic (PBFT/Raft)  
5. Consensus design is fundamentally a trade-off between **safety, liveness, and efficiency**.
