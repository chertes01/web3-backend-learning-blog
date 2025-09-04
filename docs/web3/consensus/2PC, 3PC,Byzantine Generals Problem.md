# Distributed Transaction Consistency: 2PC, 3PC, and Byzantine Generals Problem

## 1. 2PC (Two-Phase Commit)

### 1. Basic Flow

**Phase 1: Voting**

1. Coordinator → All participants: Send `CanCommit?`  
2. Participants:
   - If ready to execute the transaction → reply `Yes`
   - If unable to execute → reply `No`
3. Coordinator:
   - If everyone replies `Yes` → proceed to Phase 2
   - If any `No` or timeout → send `Abort` to all participants

**Phase 2: Commit**

1. If all agreed in Phase 1:
   - Coordinator → All participants: Send `DoCommit`
   - Participants: Execute commit and reply `Ack`
2. Coordinator:
   - Receives all `Ack` → transaction successful
   - If timeout or error → send `Abort`
   - Participants receiving `Abort` → rollback transaction

### 2. Features

- ✅ Guarantees atomicity: either all commit or all rollback  
- ❌ Limitations: cannot solve network partition blocking, single point of failure at coordinator

### 3. Split-Brain Example

Assume Coordinator C and two participants A, B:

1. C sends `Prepare` to A and B  
2. A, B reply `Vote-commit`  
3. Before C sends `Commit`, a network partition occurs:
   - A receives Commit, commits transaction  
   - B does not receive Commit, waits  

**Result:**

- A commits  
- B blocks or rolls back  
- System becomes inconsistent

### 4. Risks

- Data inconsistency  
- System blocking  
- Possible concurrent conflicts (double-write)

---

## 2. 3PC (Three-Phase Commit)

### 1. Improvements

Based on 2PC, introduce **PreCommit phase** and add timeout mechanism to reduce blocking.

**Phase 1: CanCommit**

- Coordinator asks participants if they can commit  
- Participants reply Yes or No  
- If timeout → participant automatically aborts transaction

**Phase 2: PreCommit**

- If all participants reply Yes:
  - Coordinator sends `PreCommit`  
  - Participants write Undo/Redo logs, reply Yes, start timeout  
- If final instruction not received before timeout: participants can commit directly to avoid blocking

**Phase 3: DoCommit**

- Coordinator sends `DoCommit`  
- Participants execute commit and reply Ack  
- Coordinator confirms transaction completion

### 2. Features

- ✅ Timeout mechanism avoids blocking  
- ✅ Stronger state consistency than 2PC  
- ❌ Still cannot handle network partition or Byzantine faults  
- ❌ Adds one extra communication round, higher overhead

### 3. 3PC Failure Example

System: Coordinator C, participants A, B, X

1. Phase 1: all reply Yes  
2. Phase 2: C sends `PreCommit`; A, B receive and log; X does not receive due to network delay  
3. Before Phase 3, C crashes  

**Result:**

- A, B: commit directly after timeout according to protocol  
- X: did not receive PreCommit, rolls back  

→ System becomes inconsistent

---

## 3. Byzantine Generals Problem

### 1. Problem Definition

In distributed systems, some nodes may be malicious or faulty. The core challenge of the Byzantine Generals Problem is: **how can loyal nodes reach consensus even in the presence of traitors?**

### 2. Example Scenario (9 generals, 1 traitor)

- 4 loyal generals vote **Attack**  
- 4 loyal generals vote **Retreat**  
- 1 traitor (malicious node) → spreads conflicting information  

**Result:**

- Attack faction sees 5 votes to attack → decides to attack  
- Retreat faction sees 5 votes to retreat → decides to retreat  

→ System consistency is broken

### 3. Relation to 2PC / 3PC

- 2PC: can block, single point of failure  
- 3PC: reduces blocking, but cannot handle Byzantine faults  
- Byzantine scenarios require stronger consensus algorithms (e.g., PBFT, Raft, Paxos)

---

## 4. Summary

| Protocol | Advantages | Disadvantages |
|----------|------------|---------------|
| 2PC      | Guarantees atomicity | Blocking, single point of failure, split-brain risk |
| 3PC      | PreCommit + timeout reduces blocking | Network partition or Byzantine faults may still cause inconsistency, higher communication overhead |
| Byzantine Generals | Highlights risk of malicious nodes | Standard 2PC/3PC cannot guarantee consistency; requires robust consensus algorithms |

