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

/*
功能要求：
        通证名称 tokenName
        通证简称 tokenSymbol
        通证发行数量 totalsupply
        owner 地址，记录合约创建者
        balance mapping: 记录每个地址持有的通证数量

核心函数：
        mint: 铸造通证
        transfer: 转账通证
        balanceOf: 查询某地址余额
*/

contract FundToken {

    // 通证名称，例如 "MyToken"
    string public tokenName;
    // 通证简称，例如 "MTK"
    string public tokenSymbol;
    // 通证总发行量
    uint256 public totalsupply;
    // 合约创建者地址
    address public owner;
    // balances 映射，记录每个地址的余额
    mapping (address => uint256) inernal balances;

    // constructor 构造函数，在合约部署时执行一次
    constructor(string memory _tokenName, string memory _tokenSymbol) {
        tokenName = _tokenName; // 初始化通证名称
        tokenSymbol = _tokenSymbol; // 初始化通证简称
        owner = msg.sender; // 部署者地址为 owner
    }

    // mint 函数：铸造通证，将 amountToMint 添加到调用者余额和总供应量
    function mint(uint256 amountToMint) public {
        balances[msg.sender] += amountToMint; // 调用者余额 + amountToMint
        totalsupply += amountToMint; // 总发行量 + amountToMint
    }

    // transfer 函数：转账功能
    function transfer(address payee, uint256 amount) public {
        // 检查调用者余额是否足够
        require(balances[msg.sender] >= amount, "You don't have enough balance to transfer");

        // 从调用者余额中扣除 amount
        balances[msg.sender] -= amount;

        // 给收款人 payee 添加 amount
        balances[payee] += amount;
    }

    // balanceOf 函数：查询地址余额
    function balanceOf(address addr) public view returns (uint256) {
        return balances[addr]; // 返回指定地址余额
    }
}
```

---

## 🔧 代码结构详解

#### 合约概览
```solidity
contract FundToken {
    string public tokenName;
    string public tokenSymbol;
    uint256 public totalsupply;
    address public owner;

    mapping (address => uint256) private balances;
}
```
✔️ **知识点**：
-   `string public tokenName` 定义通证名称，如 `USD Coin`。
-   `string public tokenSymbol` 定义通证简称，如 `USDC`。
-   `totalsupply` 表示当前总发行量。
-   `mapping` 结构存储地址余额。

#### 构造函数 `constructor`
```solidity
constructor(string memory _tokenName, string memory _tokenSymbol) {
    tokenName = _tokenName;
    tokenSymbol = _tokenSymbol;
    owner = msg.sender;
}
```
✔️ **解读**：
-   部署合约时，传入名称与简称。
-   `msg.sender` 为部署者地址，设置为 `owner`。

#### `mint` 函数
```solidity
function mint(uint256 amountToMint) public {
    balances[msg.sender] += amountToMint;
    totalsupply += amountToMint;
}
```
🔑 **功能**：铸造 Token，将调用者地址余额增加 `amountToMint`，同时增加 `totalsupply`。

❗ **注意**：本示例未限制 `mint` 权限，任何人都可 `mint`。正式 ERC20 发行一般限定 `onlyOwner` 调用或预先分配。

#### `transfer` 函数
```solidity
function transfer(address payee, uint256 amount) public {
    require(balances[msg.sender] >= amount, "You don't have enough balance to transfer");
    balances[msg.sender] -= amount;
    balances[payee] += amount;
}
```
🔑 **功能**：转账 Token，调用者余额大于等于转账金额才能执行，从 `sender` 扣除 `amount`，给 `payee` 增加 `amount`。

#### `balanceOf` 函数
```solidity
function balanceOf(address addr) public view returns(uint256) {
    return balances[addr];
}
```
🔑 **功能**：查询某地址的 Token 余额。

---

## ✅ 本课总结

- ✅ 理解了 **Coin** 是主网原生资产，而 **Token** 是基于智能合约的资产。
- ✅ 了解了 **ERC20** 作为一套标准接口，对生态系统互操作性的重要意义。
- ✅ 掌握了实现一个简易 Token 的三个核心动作：
    -   **`mint`**：无中生有地创造代币。
    -   **`transfer`**：在不同账户间转移代币。
    -   **`balanceOf`**：查询账户持有的代币数量。

---

## 🎯 练习拓展

1.  **添加权限控制**：为 `mint` 函数添加一个 `onlyOwner` 修饰符，使得只有合约的部署者才能铸造新的代币。
2.  **实现 `approve` 和 `transferFrom`**：（挑战）为合约添加 `approve` 和 `transferFrom` 函数，以实现完整的 ERC20 授权转账功能。思考为什么需要 `allowance` 机制，而不仅仅是 `transfer`？

