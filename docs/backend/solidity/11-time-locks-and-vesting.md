# 📘 Solidity 学习笔记 | 第十课：锁仓机制与时间控制

## 🔹 课程目标

在 `FundMe` 合约中引入时间维度的控制，学习如何实现一个有时限的众筹。

- ✅ **`block.timestamp`**：理解并使用区块时间戳进行时间判断。
- ✅ **`immutable` 关键字**：通过在构造函数中初始化变量来优化 Gas。
- ✅ **`modifier` 进阶**：创建基于时间的 `modifier` 来保护函数。
- ✅ **锁仓逻辑**：实现一个有时限的捐款窗口，以及窗口关闭后的提现/退款逻辑。

---

## 💻 完整代码

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 导入 Chainlink 预言机接口，用于获取 ETH/USD 价格
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";

contract FundMe {

    // State Variables
    mapping(address => uint256) public s_addressToAmountFunded;
    address[] public s_funders;

    uint256 public constant MINIMUM_USD = 2 * 10**18; // 2 USD

    address public immutable i_owner;
    uint256 public immutable i_deploymentTimestamp;
    uint256 public immutable i_lockTime;

    AggregatorV3Interface internal s_dataFeed;

    // Modifiers
    modifier onlyOwner() {
        require(msg.sender == i_owner, "Not owner");
        _;
    }

    modifier beforeLockTimeEnds() {
        require(block.timestamp < i_deploymentTimestamp + i_lockTime, "Lock time has passed");
        _;
    }

    modifier afterLockTimeEnds() {
        require(block.timestamp >= i_deploymentTimestamp + i_lockTime, "Lock time has not passed");
        _;
    }

    // Functions
    constructor(uint256 _lockTime) {
        i_owner = msg.sender;
        i_deploymentTimestamp = block.timestamp;
        i_lockTime = _lockTime;
        s_dataFeed = AggregatorV3Interface(0x694AA1769357215DE4FAC081bf1f309aDC325306); // Sepolia ETH/USD
    }

    function fund() external payable beforeLockTimeEnds {
        require(getConversionRate(msg.value) >= MINIMUM_USD, "Did not send enough USD!");
        s_funders.push(msg.sender);
        s_addressToAmountFunded[msg.sender] += msg.value;
    }

    function withdraw() external onlyOwner afterLockTimeEnds {
        // Reset all funders' balances
        for (uint256 i = 0; i < s_funders.length; i++) {
            address funder = s_funders[i];
            s_addressToAmountFunded[funder] = 0;
        }
        s_funders = new address[](0);

        // Withdraw the funds
        (bool success, ) = payable(i_owner).call{value: address(this).balance}("");
        require(success, "Transfer failed.");
    }

    function refund() external afterLockTimeEnds {
        uint256 amountToRefund = s_addressToAmountFunded[msg.sender];
        require(amountToRefund > 0, "No funds to refund");

        // Checks-Effects-Interactions Pattern
        s_addressToAmountFunded[msg.sender] = 0;

        (bool success, ) = payable(msg.sender).call{value: amountToRefund}("");
        require(success, "Refund failed.");
    }

    // --- Helper Functions ---

    function getPrice() public view returns (uint256) {
        (, int256 answer, , , ) = s_dataFeed.latestRoundData();
        return uint256(answer);
    }

    function getConversionRate(uint256 ethAmount) internal view returns (uint256) {
        uint256 ethPrice = getPrice();
        return (ethAmount * ethPrice) / (10**8);
    }
}
```

---

## 🔍 新知识点与代码详解

#### 1️⃣ `block.timestamp` 与 `immutable` 关键字
> **`block.timestamp`** 是 Solidity 的全局变量，返回当前区块的 UNIX 时间戳（秒）。我们用它来标记合约的部署时间 `i_deploymentTimestamp`。
> **`immutable`** 变量可以在**构造函数 (`constructor`) 中被赋值一次**。它非常适合用于在部署时才确定的值，如 `i_owner` 地址、`i_lockTime` 时长等。赋值后，它的值被硬编码到字节码中，访问它时**能显著节省 Gas**。

```solidity
// 在构造函数中初始化 immutable 变量
constructor(uint256 _lockTime) {
    i_owner = msg.sender;
    i_deploymentTimestamp = block.timestamp; // 记录部署时间
    i_lockTime = _lockTime; // 设置锁仓时长
}
```

#### 2️⃣ 基于时间的 `modifier`
> 我们创建了两个新的修饰符来封装时间检查逻辑，使代码更清晰、更可重用。

```solidity
modifier beforeLockTimeEnds() {
    // 检查当前时间是否在锁仓结束时间之前
    require(block.timestamp < i_deploymentTimestamp + i_lockTime, "Lock time has passed");
    _;
}

modifier afterLockTimeEnds() {
    // 检查当前时间是否在锁仓结束时间之后或等于
    require(block.timestamp >= i_deploymentTimestamp + i_lockTime, "Lock time has not passed");
    _;
}
```

#### 3️⃣ `fund()` 函数：新增锁仓校验
> 通过应用 `beforeLockTimeEnds` 修饰符，我们确保了 `fund` 函数只能在锁仓期内被调用。

```solidity
function fund() external payable beforeLockTimeEnds {
    // 1. 检查金额（来自旧逻辑）
    require(getConversionRate(msg.value) >= MINIMUM_USD, "Did not send enough USD!");
    // 2. 检查时间（由 modifier 完成）
    // 3. 执行核心逻辑
    s_funders.push(msg.sender);
    s_addressToAmountFunded[msg.sender] += msg.value;
}
```

#### 4️⃣ `withdraw()` 提现函数：新增时间与权限控制
> 通过组合使用 `onlyOwner` 和 `afterLockTimeEnds` 修饰符，我们确保了只有 `owner` 才能在锁仓期结束后提现。

```solidity
function withdraw() external onlyOwner afterLockTimeEnds {
    // ... 提现逻辑 ...
}
```

#### 5️⃣ `refund()` 退款函数：新增时间控制
> 同样，`refund` 函数也应用了 `afterLockTimeEnds` 修饰符，意味着用户只能在众筹窗口关闭后才能申请退款。

```solidity
function refund() external afterLockTimeEnds {
    // 1. 检查用户是否有可退款余额
    uint256 amountToRefund = s_addressToAmountFunded[msg.sender];
    require(amountToRefund > 0, "No funds to refund");

    // 2. 先更新状态（清零），防止重入攻击
    s_addressToAmountFunded[msg.sender] = 0;

    // 3. 最后执行转账
    (bool success, ) = payable(msg.sender).call{value: amountToRefund}("");
    require(success, "Refund failed.");
}
```

---

## 🔧 调用图示（逻辑流程）

1.  **部署合约**
    -   调用 `constructor(86400)`，传入锁仓时长为 1 天 (86400 秒)。
    -   `i_owner`, `i_deploymentTimestamp`, `i_lockTime` 被初始化并设为不可变。

2.  **锁仓期内 (第 1-23 小时)**
    -   **用户A 调用 `fund()`**：`beforeLockTimeEnds` 检查通过，捐款成功。
    -   **Owner 调用 `withdraw()`**：`afterLockTimeEnds` 检查失败，交易 `revert`。
    -   **用户A 调用 `refund()`**：`afterLockTimeEnds` 检查失败，交易 `revert`。

3.  **锁仓期后 (第 25 小时)**
    -   **用户B 调用 `fund()`**：`beforeLockTimeEnds` 检查失败，交易 `revert`。
    -   **用户A 调用 `refund()`**：`afterLockTimeEnds` 检查通过，用户A 成功取回自己的捐款。
    -   **Owner 调用 `withdraw()`**：`afterLockTimeEnds` 检查通过，Owner 成功提走合约中剩余的所有资金。

---

## ✅ 本课总结

- ✅ 学会了使用 `block.timestamp` 来实现基于时间的逻辑控制。
- ✅ 掌握了 `immutable` 关键字，并理解了它与 `constant` 的区别和 Gas 优化优势。
- ✅ 能够编写和使用 `modifier` 来封装可重用的检查逻辑，如时间锁和权限控制。
- ✅ 成功构建了一个包含捐款窗口、锁仓期、管理员提现和用户退款等完整功能的众筹合约。

---

## 🎯 练习拓展

1.  **分阶段解锁**：修改合约，实现一个可以分多次提现的功能。例如，锁仓期结束后，Owner 每周只能提现总额的 10%。
2.  **软顶和硬顶**：为众筹添加一个“硬顶”（Hard Cap），即一个总募资目标。一旦达到这个目标，`fund` 函数就应立即停止接受捐款，即使还没到锁仓期结束时间。