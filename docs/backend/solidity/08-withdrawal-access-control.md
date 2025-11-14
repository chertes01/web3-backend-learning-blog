# 📘 Solidity 学习笔记 | 第七课：提现函数与权限控制

## 🔹 课程目标

在 `FundMe` 合约的基础上，实现提现功能和权限管理，学习：

- ✅ **`constant` 关键字**：定义不可变常量，优化 Gas。
- ✅ **`owner` 权限模式**：设置合约管理员，保护关键函数。
- ✅ **`address(this).balance`**：获取合约自身的 ETH 余额。
- ✅ **提现函数**：实现一个只有 `owner` 才能调用的提现逻辑。
- ✅ **`payable(address).transfer()`**：掌握基础的 ETH 转账方法。

---

## 💻 完整代码

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract FundMe {

    // 记录每个用户的捐款金额
    mapping(address => uint256) public FundsToAmount;

    // 最低捐款额，单位 USD * 10^18
    uint256 public constant MINIMUM_USD = 1 * 10 ** 17; // 0.1 USD

    // constant 常量，表示提现最低金额阈值，部署后不可更改，Gas 更优
    uint256 public constant WITHDRAW_THRESHOLD_USD = 1 * 10 ** 18; // 1 USD

    // 合约拥有者地址
    address public owner;

    // Chainlink 预言机接口
    AggregatorV3Interface internal dataFeed;

    // 构造函数，在合约部署时执行一次
    constructor() {
        // 初始化 ETH/USD 预言机地址 (Sepolia)
        dataFeed = AggregatorV3Interface(0x694AA1769357215DE4FAC081bf1f309aDC325306);
        // 设置合约部署者为 owner
        owner = msg.sender;
    }

    // 允许 owner 转移合约所有权
    function transferOwnership(address new_owner) public {
        // 只有当前 owner 可调用
        require(msg.sender == owner, "You are not the owner!");
        owner = new_owner;
    }

    // 捐款函数
    function fund() external payable {
        // 验证捐款金额（按 USD 计价） >= 最低值
        require(getConversionRate(msg.value) >= MINIMUM_USD, "Did not send enough USD!");
        // 记录捐款金额
        FundsToAmount[msg.sender] += msg.value;
    }

    // 提现函数
    function withdraw() external {
        // 1. 验证合约余额按 USD 计价 >= 提现要求金额
        require(getConversionRate(address(this).balance) >= WITHDRAW_THRESHOLD_USD, "Not enough funds to withdraw!");
        // 2. 验证调用者是否为 owner
        require(msg.sender == owner, "You are not the owner!");

        // 3. 将合约全部余额转给 owner
        payable(msg.sender).transfer(address(this).balance);
    }

    // --- 辅助函数 ---

    function getPrice() public view returns (uint256) {
        (, int256 answer, , , ) = dataFeed.latestRoundData();
        return uint256(answer);
    }

    function getConversionRate(uint256 ethAmount) internal view returns (uint256) {
        uint256 ethPrice = getPrice();
        // ETH_Amount (wei) * ETH_Price (USD / ETH * 10^8) / 10^8 = USD * 10^18
        return (ethAmount * ethPrice) / (10 ** 8);
    }
}
```

---

## 🔍 新知识点解析

#### 1️⃣ `constant` 关键字
> **定义**：`uint256 public constant WITHDRAW_THRESHOLD_USD = 1 * 10 ** 18;`
> `constant` 用于声明一个**常量**。该变量的值在编译时就已确定，并直接硬编码到合约的字节码中。
> **优势**：因为它不占用合约的存储槽（storage slot），每次访问它时都无需从存储中读取，从而**显著节省 Gas**。

#### 2️⃣ `owner` 权限控制模式
> **定义**：`address public owner;`
> 这是一个非常常见的智能合约设计模式。
> 1.  在状态变量中声明一个 `owner` 地址。
> 2.  在 `constructor` 中，将 `owner` 初始化为 `msg.sender`（即合约的部署者）。
> 3.  在需要保护的函数（如 `withdraw`）中，添加 `require(msg.sender == owner, ...)` 来验证调用者身份。

#### 3️⃣ `address(this).balance`
> **含义**：这是一个特殊的表达式，用于获取**当前合约地址上持有的 ETH 余额**，单位是 `wei`。
> 它是合约资金管理的核心，常用于提现函数或余额检查逻辑。

#### 4️⃣ `payable(address).transfer(amount)`
> **用途**：这是 Solidity 中用于发送 ETH 的一种基础方法。
> - `payable(msg.sender)`：将 `msg.sender` 的地址（`address` 类型）强制转换为 `address payable` 类型，使其能够接收 ETH。
> - `.transfer(address(this).balance)`：调用 `transfer` 函数，将指定数量的 ETH（这里是合约的全部余额）发送到该地址。
> **注意**：`transfer` 有 2300 Gas 的限制，足以完成简单的转账，但如果接收方是执行复杂逻辑的合约，可能会因 Gas 不足而失败。

---

## 🔧 调用图示（逻辑流程）

1.  **部署合约**
    - `constructor` 被调用，`owner` 状态变量被设置为**部署者的地址**。

2.  **用户A 调用 `fund()` 并发送 1 ETH**
    - 交易成功，合约余额变为 1 ETH。

3.  **用户B (非 owner) 调用 `withdraw()`**
    - 第一个 `require` 检查余额，通过。
    - 第二个 `require(msg.sender == owner, ...)` 检查身份，`msg.sender` 是用户B的地址，不等于 `owner` 地址，**交易失败并回滚**。

4.  **Owner 调用 `withdraw()`**
    - 第一个 `require` 检查余额，通过。
    - 第二个 `require` 检查身份，`msg.sender` 是 `owner` 地址，通过。
    - `payable(owner).transfer(1 ETH)` 执行，合约的 1 ETH 被转入 `owner` 的钱包地址。
    - 交易成功，合约余额清零。

---

## ✅ 本课总结

- ✅ 学会了使用 `constant` 定义常量以优化 Gas 并提高代码可读性。
- ✅ 掌握了 Solidity 中最基础也是最重要的 `owner` 权限控制模式。
- ✅ 理解了 `address(this).balance` 是获取合约自身余额的关键。
- ✅ 成功实现了一个包含权限校验和资金转移的 `withdraw` 提现函数。

---

## 🎯 练习拓展

1.  **创建 `modifier`**：将 `require(msg.sender == owner, ...)` 这行代码抽象成一个名为 `onlyOwner` 的 `modifier`，并将其应用到 `withdraw` 和 `transferOwnership` 函数上。
2.  **事件日志**：为 `withdraw` 函数添加一个 `event`，在每次成功提现时，记录提现的金额和时间戳。
3.  **更安全的提现**：（挑战）研究并使用 `call` 方法替代 `transfer` 来实现提现，并解释这样做的好处和需要注意的安全问题（如重入攻击）。