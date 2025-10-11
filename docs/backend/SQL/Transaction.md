# 🛡️ SQL Transaction Control Study Notes

---

## Table of Contents

- 🤔 What is a Transaction? An Analogy
- 💎 The Four ACID Properties of Transactions
- 🚦 Transaction Isolation Levels (Supplementary Knowledge)
- 🕹️ Methods of Transaction Control
- 🏦 Practical Exercise: Simulating Bank Transfers
- ✨ Key Points Summary

---

## 1. 🤔 What is a Transaction? An Analogy

**Definition**: A transaction is a set of database operations that are treated as an indivisible unit. They either all succeed or all fail.

**ATM Withdrawal Example**:

1. Insert card and enter password
2. Enter withdrawal amount
3. ATM dispenses cash
4. Bank account is debited

These four steps must be treated as a whole. If the ATM dispenses cash but the account isn't debited (the bank loses), or the account is debited but the ATM doesn't dispense cash (you lose), both are disasters. The transaction mechanism ensures that such operations are "all-or-nothing".

---

## 2. 💎 The Four ACID Properties of Transactions

ACID is the cornerstone of transaction reliability and a frequent topic in database interviews.

- **A - Atomicity**  
  Meaning: All operations in a transaction are completed or all are canceled.  
  Analogy: In a transfer, both debit and credit must happen together.

- **C - Consistency**  
  Meaning: Before and after a transaction, the database transitions from one valid state to another.  
  Analogy: No matter how many transfers, the total amount in all bank accounts remains unchanged.

- **I - Isolation**  
  Meaning: When multiple transactions are executed concurrently, one transaction's execution should not be affected by others.  
  Analogy: When you withdraw money from an ATM, others cannot see the intermediate changes in your account balance.

- **D - Durability**  
  Meaning: Once a transaction is committed, its results are permanent.  
  Analogy: After the ATM shows "Transaction Successful", even if the bank system loses power, the transaction will not be lost.

---

## 3. 🚦 Transaction Isolation Levels (Supplementary Knowledge)

Isolation is not absolute. To balance performance and data consistency, SQL defines four isolation levels:

| Isolation Level             | Dirty Read | Non-repeatable Read | Phantom Read |
|----------------------------|------------|---------------------|--------------|
| Read Uncommitted           | Possible   | Possible            | Possible     |
| Read Committed             | Avoided    | Possible            | Possible     |
| Repeatable Read            | Avoided    | Avoided             | Possible     |
| Serializable               | Avoided    | Avoided             | Avoided      |

💡 MySQL InnoDB default isolation level: **Repeatable Read**, which solves most concurrency problems.

---

## 4. 🕹️ Methods of Transaction Control

In MySQL, there are two main ways to control transactions:

### Method 1: Temporarily Disable Auto-commit

By default, each SQL statement is a transaction and is automatically committed. You can temporarily disable this:

```sql
SET @@autocommit = 0; -- Disable auto-commit

-- ... Execute your SQL operations ...

COMMIT;    -- Manually commit
-- Or ROLLBACK; -- Manually rollback

SET @@autocommit = 1; -- Restore default setting
```

### Method 2: Explicitly Start a Transaction (Recommended)

More commonly used, clearer, smaller scope, best practice:

```sql
START TRANSACTION;
-- ...SQL operations...
COMMIT;    -- Commit
-- Or ROLLBACK; -- Rollback
```

---

## 5. 🏦 Practical Exercise: Simulating Bank Transfers

**Scenario**: Zhang San (money: 2000) transfers 1000 yuan to Li Si (money: 2000).

### Preparation

```sql
CREATE TABLE account (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(10),
    money INT
);
INSERT INTO account(name, money) VALUES ('Zhang San', 2000), ('Li Si', 2000);
```

### Transfer Process

```sql
-- 1. Start transaction
START TRANSACTION;

-- 2. Deduct 1000 from Zhang San's account
UPDATE account SET money = money - 1000 WHERE name = 'Zhang San';
-- At this point, only in the current session, Zhang San's balance decreases; others cannot see it.

-- (Suppose an unexpected event occurs here, such as a program crash or network interruption)

-- 3. Add 1000 to Li Si's account
UPDATE account SET money = money + 1000 WHERE name = 'Li Si';

-- 4. Check the result and decide the final outcome
COMMIT;    -- If all goes well, commit, and everyone can see the new balance
-- Or
ROLLBACK;  -- If there is an error, rollback, and the database returns to the state before the transaction started
```

---

## 6. ✨ Key Points Summary

- Transactions are the core mechanism for ensuring data consistency.
- ACID are the four golden rules that transactions must follow.
- Use `START TRANSACTION` to clearly manage transaction units; safest and most reliable.
- On success, use `COMMIT` to commit.
- On error, use `ROLLBACK` to rollback and protect data safety.

---