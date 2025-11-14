# 📘 Solidity 学习笔记 | 第十一课：ERC20 通证基础

## 🔹 课程目标

学习什么是 **ERC20** 标准，理解 **Coin** 与 **Token** 的区别，并通过编写一个简易的 Token 合约，掌握其核心功能。

- ✅ **ERC20 标准**：了解其定义和核心作用。
- ✅ **Coin vs. Token**：掌握区块链主网原生币与合约代币的根本区别。
- ✅ **核心函数**：通过代码实践 `mint`、`transfer` 和 `balanceOf` 的实现原理。
- ✅ **简易 Token 合约**：从零开始编写一个自己的 Token。

---

## 🔍 新知识点解析

#### 1️⃣ Coin 与 Token 的区别

> 这是理解加密资产分类的基础。

| 类型 | 定义 | 示例 |
| :--- | :--- | :--- |
| **Coin (币)** | 区块链**主网的原生代币**，用于支付 Gas 费和网络安全。 | BTC (比特币), ETH (以太坊) |
| **Token (通证)** | 基于某条区块链（如以太坊）的**智能合约**发行的代币。 | USDT, LINK, SHIB (都是 ERC20 Token) |

简单来说，**Coin** 是区块链的“燃料”，而 **Token** 是运行在这条区块链上的“应用程序”发行的资产。

#### 2️⃣ ERC20 标准
> **ERC20** (Ethereum Request for Comments 20) 是以太坊上最成功的、应用最广泛的**同质化代币标准**。它定义了一套所有 Token 合约都应遵循的通用接口（函数和事件）。

**核心作用**：确保了不同项目方发行的 Token 能被所有钱包、交易所和 DeFi 协议无缝集成和交互。

**主要函数接口**：
-   `totalSupply()`: 返回 Token 总供应量。
-   `balanceOf(address account)`: 返回指定账户的余额。
-   `transfer(address recipient, uint256 amount)`: 将 Token 从调用者账户转给接收者。
-   `approve(address spender, uint256 amount)`: 授权 `spender` 从调用者账户中提取不超过 `amount` 的 Token。
-   `transferFrom(address sender, address recipient, uint256 amount)`: 由 `spender` 调用，从 `sender` 账户向 `recipient` 转账。
-   `allowance(address owner, address spender)`: 查询 `spender` 仍被授权从 `owner` 账户提取的额度。

> 本课将实现其中最基础的 `balanceOf` 和 `transfer` 功能，并添加一个 `mint` 函数来创建代币。

---

## 💻 完整代码

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract FundToken {
    // State Variables
    string public tokenName;
    string public tokenSymbol;
    uint256 public totalsupply;
    address public owner;

    // balances 映射，记录每个地址的余额
    // 设为 private，因为我们提供了自定义的 public balanceOf 函数
    mapping(address => uint256) private balances;

    // Functions
    constructor(string memory _tokenName, string memory _tokenSymbol) {
        tokenName = _tokenName;
        tokenSymbol = _tokenSymbol;
        owner = msg.sender;
    }

    /**
     * @notice 铸造通证，将 amountToMint 添加到调用者余额和总供应量
     * @dev 在真实的 ERC20 合约中，此函数通常会带有权限控制
     */
    function mint(uint256 amountToMint) public {
        balances[msg.sender] += amountToMint;
        totalsupply += amountToMint;
    }

    /**
     * @notice 将 `amount` 数量的通证从调用者账户转移到 `payee` 账户
     */
    function transfer(address payee, uint256 amount) public {
        require(balances[msg.sender] >= amount, "Not enough balance to transfer");
        balances[msg.sender] -= amount;
        balances[payee] += amount;
    }

    /**
     * @notice 查询指定地址 `addr` 的余额
     */
    function balanceOf(address addr) public view returns (uint256) {
        return balances[addr];
    }
}
```

---

## 🔧 调用图示（逻辑流程）

1.  **部署合约**
    -   调用 `constructor("MyToken", "MTK")`。
    -   `tokenName` 被设为 "MyToken"，`tokenSymbol` 被设为 "MTK"。
    -   `owner` 被设为部署者的地址。
    -   `totalsupply` 初始为 `0`。

2.  **用户A 调用 `mint(1000)`**
    -   `balances[用户A]` 增加 `1000`。
    -   `totalsupply` 增加 `1000`，现在总供应量为 `1000`。

3.  **用户A 调用 `balanceOf(用户A)`**
    -   函数返回 `1000`。

4.  **用户A 调用 `transfer(用户B, 300)`**
    -   `require` 检查通过 (1000 >= 300)。
    -   `balances[用户A]` 减少 `300`，变为 `700`。
    -   `balances[用户B]` 增加 `300`，变为 `300`。

5.  **用户B 调用 `balanceOf(用户B)`**
    -   函数返回 `300`。

---

## ✅ 本课总结

- ✅ 理解了 **Coin** 是主网原生资产，而 **Token** 是基于智能合约的资产。
- ✅ 了解了 **ERC20** 作为一套标准接口，对生态系统互操作性的重要意义。
- ✅ 掌握了实现一个简易 Token 的三个核心动作：
    -   **`mint`**：无中生有地创造代币。
    -   **`transfer`**：在不同账户间转移代币。
    -   **`balanceOf`**：查询账户持有的代-   币数量。

---

## 🎯 练习拓展

1.  **添加权限控制**：为 `mint` 函数添加一个 `onlyOwner` 修饰符，使得只有合约的部署者才能铸造新的代币。
2.  **实现 `approve` 和 `transferFrom`**：（挑战）为合约添加 `approve` 和 `transferFrom` 函数，以实现完整的 ERC20 授权转账功能。思考为什么需要 `allowance` 机制，而不仅仅是 `transfer`？