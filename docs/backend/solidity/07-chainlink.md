# 📘 Solidity 学习笔记 | 第六课：Chainlink 预言机与价格换算

## 🔹 课程目标

学会如何使用 **Chainlink 预言机**获取链下数据（如 ETH/USD 实时价格），并将 `FundMe` 合约升级为按美元价值进行验证。

- ✅ **Chainlink 预言机**：理解其作用并获取 ETH/USD 实时价格。
- ✅ **`AggregatorV3Interface`**：学会使用标准接口与预言机交互。
- ✅ **`constructor` 构造函数**：在合约部署时初始化预言机地址。
- ✅ **价格换算**：掌握如何将 ETH 数额转换为对应的 USD 价值。
- ✅ **单位一致性**：理解在进行数值比较时，保持单位一致的重要性。

---

## 💻 完整代码

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 引入 Chainlink 预言机标准接口 AggregatorV3Interface
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract FundMe {

    // 记录每个地址的投资金额 (单位: wei)
    mapping(address => uint256) public FundsToAmount;

    // 设置最低捐款额为 1 USD。注意这里的单位处理，详见下文解析。
    uint256 public constant MINIMUM_USD = 1 * 10 ** 18;

    // 声明 Chainlink 预言机接口变量
    AggregatorV3Interface internal dataFeed;

    // 构造函数，在合约部署时自动执行一次
    constructor() {
        // 初始化 dataFeed，这里使用的是 Chainlink Sepolia 测试网的 ETH/USD 价格预言机地址
        dataFeed = AggregatorV3Interface(0x694AA1769357215DE4FAC081bf1f309aDC325306);
    }

    // fund 函数，允许用户发送 ETH 进行众筹
    function fund() external payable {
        // 验证用户转入的 ETH 按美元计价后，是否 >= 最低捐款额
        require(getConversionRate(msg.value) >= MINIMUM_USD, "Did not send enough USD!");

        // 记录用户捐款金额到 mapping 中
        FundsToAmount[msg.sender] += msg.value;
    }

    // 获取 ETH/USD 最新价格
    function getPrice() public view returns (uint256) {
        // 调用 latestRoundData()，返回多个值，我们仅需要第二个值 answer (价格)
        (, int256 answer, , , ) = dataFeed.latestRoundData();
        // Chainlink 返回的价格带有 8 位小数，例如 $2000.12345678 会返回 200012345678
        // 我们需要将其转换为无符号整数
        return uint256(answer);
    }

    // 将输入的 ETH 数量 (单位: wei) 转换为对应的 USD 价值
    function getConversionRate(uint256 ethAmount) public view returns (uint256) {
        uint256 ethPrice = getPrice();
        // 核心换算逻辑，详见下文解析
        uint256 ethAmountInUsd = (ethPrice * ethAmount) / 10**8;
        return ethAmountInUsd;
    }
}
```

---

## 🔍 新知识点解析

#### 1️⃣ `import {AggregatorV3Interface}`
> **接口导入**：我们从 Chainlink 提供的官方合约库中导入 `AggregatorV3Interface`。这是一个标准的接口定义，它告诉我们的合约如何与 Chainlink 的价格数据源进行通信，例如它规定了必须有一个名为 `latestRoundData()` 的函数。

#### 2️⃣ `constructor()`
> **构造函数**：`constructor` 是一个特殊的函数，它只在合约**首次部署到区块链上时**执行一次。我们用它来完成初始化设置，比如在这里，我们将 `dataFeed` 变量指向了 Sepolia 测试网上 ETH/USD 价格预言机的**具体地址**。

#### 3️⃣ `latestRoundData()`
> **获取数据**：这是 `AggregatorV3Interface` 接口中最重要的函数。调用它会从预言机返回一组最新的数据，包括价格、时间戳等。我们只关心价格（`answer`），所以用 `(, int256 answer, , , )` 的语法来忽略其他返回值。

#### 4️⃣ 价格单位与换算 (核心)
> 这是本课最关键的知识点。为了在 `require` 中正确比较，我们必须确保两边的单位完全一致。

**让我们来分析 `getConversionRate` 函数中的单位：**

`uint256 ethAmountInUsd = (ethPrice * ethAmount) / 10**8;`

| 变量 | 原始单位 | 解释 |
| :--- | :--- | :--- |
| `ethAmount` (`msg.value`) | `wei` | `1 ETH = 10^18 wei` |
| `ethPrice` (`getPrice()`) | `USD / ETH * 10^8` | Chainlink 返回的价格，带 8 位小数 |
| `ethPrice * ethAmount` | `wei * USD / ETH * 10^8` | 两者相乘 |
| `ethAmountInUsd` | `wei * USD / ETH` | 除以 `10^8` 修正小数位 |

**结果是什么？**
由于 `1 ETH` 等于 `10^18 wei`，所以 `wei / ETH` 等于 `1 / 10^18`。
因此，`ethAmountInUsd` 的最终单位是 `USD * 10^18`。

**为什么 `MINIMUM_USD` 要等于 `1 * 10**18`？**
因为 `getConversionRate` 返回值的单位是 `USD * 10^18`，所以我们的最低要求 `MINIMUM_USD` 也必须是相同的单位，这样 `require(getConversionRate(...) >= MINIMUM_USD)` 的比较才有意义。
`1 * 10**18` 在这里代表的正是 **"1 美元"** 在我们统一单位制下的数值。

---

## 🔧 调用图示（逻辑流程）

1.  **部署合约**
    - `constructor` 被调用，`dataFeed` 变量被设置为 Chainlink 预言机的地址。

2.  **用户调用 `fund()` 并发送 0.0005 ETH**
    - `msg.value` 等于 `5 * 10^14 wei`。
    - 合约调用 `getConversionRate(5 * 10^14)`。

3.  **合约内部执行换算**
    - 合约向 `dataFeed` 地址请求价格，假设返回 `2000 * 10^8` (即 $2000/ETH)。
    - 计算 `(2000 * 10^8 * 5 * 10^14) / 10^8`，结果约等于 `1 * 10^18`。

4.  **执行 `require` 校验**
    - 检查 `1 * 10^18 >= MINIMUM_USD (1 * 10^18)`，结果为 `true`。
    - 交易继续，用户的捐款被记录。

---

## ✅ 本课总结

- ✅ **Chainlink 预言机**是连接智能合约与真实世界数据的桥梁。
- ✅ 使用 `AggregatorV3Interface` 是与 Chainlink 价格源交互的标准方式。
- ✅ `constructor` 是初始化合约状态（如设置关键地址）的理想位置。
- ✅ 在 Solidity 中进行数学运算时，**处理好单位和小数位**至关重要，否则会导致逻辑错误和安全漏洞。

---

## 🎯 练习拓展

1.  **添加提款功能**：为合约添加 `owner` 和 `withdraw` 功能，只允许合约的部署者提走所有资金。
2.  **记录所有捐赠者**：添加一个 `address[] public funders;` 数组，每次有新的地址捐款时，将其添加到这个数组中（注意处理重复捐赠者）。