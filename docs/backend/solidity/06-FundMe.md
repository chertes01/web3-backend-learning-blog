# 📘 Solidity 学习笔记 | 第五课：FundMe 合约与收款基础

## 🔹 课程目标

通过一个简单的众筹合约 `FundMe`，掌握 Solidity 中与资金处理相关的核心概念：

- ✅ **收款函数**：如何让合约接收 ETH。
- ✅ **`payable` 修饰符**：允许函数处理以太币转账。
- ✅ **`require` 校验**：如何设置条件并验证转账金额。
- ✅ **`mapping` 数据记录**：如何记录每个投资者的存款信息。
- ✅ **`msg.value`**：如何获取随函数调用一同发送的 ETH 数量。

---

## 💻 完整代码

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract FundMe{

    // 定义一个 mapping，记录每个地址的投资金额
    mapping(address => uint256) public FundsToAmount;

    // 最小投资额，单位为 wei，这里设置为 1 ETH
    uint256 public MINIMUM_VALUE = 1 * 10 ** 18;

    // fund() 函数，external 表示只能外部调用，payable 允许接收 ETH
    function fund() external payable {
        // require 语句：如果投资金额小于最小要求，则 revert 交易
        require(msg.value >= MINIMUM_VALUE, "Did not send enough ETH!");

        // 记录投资者的转账金额 (此处为覆盖式，可改为累加)
        FundsToAmount[msg.sender] = msg.value;
    }
}
```

---

## 🔍 新知识点解析

#### 1️⃣ `payable` 修饰符
> **定义**：`payable` 是一个函数修饰符，它允许该函数在被调用时接收以太币（ETH）。如果一个函数没有 `payable` 修饰符，任何向其发送 ETH 的尝试都会失败。

```solidity
function fund() external payable { ... }
```

#### 2️⃣ `msg.value`
> **定义**：与 `msg.sender` 类似，`msg.value` 是一个全局变量，它代表了**随本次函数调用一同发送的 ETH 数量**，单位是 `wei`。

**单位换算**：`1 ether == 10^18 wei`。因此，`1 * 10 ** 18` 就代表 1 ETH。

#### 3️⃣ `require` 语句
> **用途**：`require` 用于设置一个必须为 `true` 的条件。如果条件为 `false`，交易将被**回滚 (revert)**，所有状态变更都会被撤销，剩余的 Gas 会退还给调用者。

```solidity
require(msg.value >= MINIMUM_VALUE, "Did not send enough ETH!");
```
- **条件**：`msg.value >= MINIMUM_VALUE`
- **失败提示信息**：`"Did not send enough ETH!"`

#### 4️⃣ `mapping(address => uint256)`
> **用途**：我们再次使用 `mapping`，这次是用来记录每个投资者（`address`）投资了多少钱（`uint256`）。这是一种非常高效的数据记录方式。

---

## 🔧 调用图示（逻辑流程）

1.  **部署合约**
    - `MINIMUM_VALUE` 被初始化为 `1 * 10 ** 18`。
    - `FundsToAmount` mapping 为空。

2.  **用户A 调用 `fund()` 并发送 0.5 ETH**
    - `msg.value` 等于 `0.5 * 10 ** 18`。
    - `require` 语句检查 `0.5 * 10 ** 18 >= 1 * 10 ** 18`，结果为 `false`。
    - **交易失败并回滚**，合约状态未改变，用户A的 0.5 ETH 未被扣除（仅消耗少量 Gas）。

3.  **用户B 调用 `fund()` 并发送 2 ETH**
    - `msg.value` 等于 `2 * 10 ** 18`。
    - `require` 语句检查 `2 * 10 ** 18 >= 1 * 10 ** 18`，结果为 `true`。
    - `FundsToAmount[用户B的地址]` 被设置为 `2 * 10 ** 18`。
    - **交易成功**，合约收到了 2 ETH。

---

## ✅ 本课总结

- ✅ 掌握了如何使用 `payable` 关键字定义一个可以接收 ETH 的函数。
- ✅ 理解了 `msg.value` 是获取随调用发送的 ETH 金额的关键。
- ✅ 学会了使用 `require` 来强制执行规则（如最低存款额），保护合约逻辑。
- ✅ 再次实践了 `mapping`，用它来高效地记录与地址相关的数据。

---

## 🎯 练习拓展

1.  **累计存款**：修改 `fund` 函数，将 `FundsToAmount[msg.sender] = msg.value;` 改为 `FundsToAmount[msg.sender] += msg.value;`，实现存款的累加而不是覆盖。
2.  **查看合约余额**：添加一个名为 `getBalance()` 的 `public view` 函数，返回合约当前的总 ETH 余额。（提示：使用 `address(this).balance`）。
3.  **设置管理员**：添加一个 `address public owner;` 状态变量，并在 `constructor` 中将其设置为 `msg.sender`，为后续的提款功能做准备。
4.  **提款功能**：（挑战）添加一个 `withdraw()` 函数，只允许 `owner` 调用，将合约中的所有 ETH 转给 `owner`。