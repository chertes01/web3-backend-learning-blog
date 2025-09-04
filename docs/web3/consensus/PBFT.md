# PBFT and Byzantine Fault Tolerance Study Notes

PBFT (Practical Byzantine Fault Tolerance) is a consensus algorithm proposed for consortium or private blockchain scenarios, designed to achieve high throughput and strong consistency even when nodes may exhibit Byzantine behavior. Its theoretical foundation comes from the Byzantine Generals Problem.

---

## 1. Review of the Byzantine Generals Problem

- **Problem Description**: How can loyal nodes reach a consensus decision in the presence of traitors (malicious nodes)?
- **Key Conditions**:
  1. **Reliable Communication**: Messages must be successfully delivered to their destination.
  2. **Authentication**: The receiver must identify the sender to ensure messages cannot be forged.
  3. **Message Loss Detection**: Ability to detect lost messages and retransmit.
- **Oral Message Protocol**:
  - The system needs at least `3m+1` generals to tolerate `m` traitors.
  - Example:
    - **n=3, m=1**: A traitor may cause confusion or prevent a decision.
    - **n=4, m=1**: Loyal generals can usually identify the traitor through mutual verification and reach consensus.
- **Written Message Protocol**:
  - Introduces **message signatures** to prevent forgery.
  - The receiver can verify the signature to ensure message reliability.
  - This is equivalent to modern networks where messages are signed with a node's private key and verified with the public key.

> **Summary**: Byzantine Fault Tolerance (BFT) ensures correct system decisions by evaluating received messages, even in the presence of malicious nodes.

---

## 2. PBFT Basic Concepts

- **State Machine Byzantine Protocol**: All nodes maintain the same system state to ensure consistent actions.
- **Node Types**: 1 primary node (Primary), multiple replica nodes (Replicas).
- **System Requirement**: For `f` Byzantine nodes, the total number of nodes must be at least `3f+1`.
- **Design Goals**:
  - Reduce the complexity of Byzantine protocols (from exponential to polynomial)
  - Support high throughput and low latency for consortium and private chains

---

## 3. PBFT Consensus Protocol Process (5 Phases)

1. **Request**
   - The client sends a request:
   ```text
   m = [op, ts, c-id, c-sig]
   ```
   * `op`: operation  
   * `ts`: timestamp  
   * `c-id`: client ID  
   * `c-sig`: client signature

2. **Pre-Prepare (Sequence Assignment)**
   * The primary assigns a sequence number `sn` and broadcasts:
   ```text
   [PP, vn, sn, D(m), p-sig, m]
   ```
   * `PP`: Pre-Prepare message  
   * `vn`: view number  
   * `D(m)`: message digest  
   * `p-sig`: primary's signature

3. **Prepare (Interaction)**
   * After receiving the PP message, replicas broadcast to other replicas:  
     `[P, vn, sn, D(m), b-id, b-sig]`

4. **Commit (Sequence Confirmation)**
   * After receiving `2f+1` Prepare messages, broadcast Commit:  
     `[C, vn, sn, D(m), b-id, b-sig]`

5. **Reply / Execute (Response and Execution)**
   * After receiving `2f+1` Commit messages, execute the operation and