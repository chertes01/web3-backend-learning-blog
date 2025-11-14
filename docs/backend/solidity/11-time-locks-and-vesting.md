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
// 许可证声明，指明此合约的开源协议为 MIT，可以自由复制、使用、修改
pragma solidity ^0.8.20;


// 导入 Chainlink 预言机接口，用于获取 ETH/USD 价格
import {AggregatorV3Interface} from "@chainlink/contracts/src/v0.8/shared/interfaces/AggregatorV3Interface.sol";


contract FundMe {


// 定义映射，记录每个用户最后一次 fund 的金额
mapping(address => uint256) public FundsToAmount;

// 定义映射，记录每个用户累计的 fund 总金额
mapping(address => uint256) public User_Amount;

// 定义最小转账金额，2 * 10^18，单位为 USD * 10^18（使用 Chainlink 数据时常以 8 位精度返回，需要自行换算）
uint256 MINIMUM_VALUE = 2 * 10 ** 18; // 2 USD

// 定义目标金额常量，10 * 10^18，达到此值后 owner 可提取资金
uint256 constant getAmount = 10 * 10 ** 18;

// 合约拥有者地址
address public owner;

//时间戳
uint256 deploymentTimestamp;
//锁仓时长
uint256 lockTime;

// Chainlink 预言机接口对象
AggregatorV3Interface internal dataFeed;

// 构造函数，部署合约时执行
constructor(uint256 _lockTime) {
    // 指定 Chainlink Sepolia testnet ETH/USD 价格预言机合约地址
    dataFeed = AggregatorV3Interface(0x694AA1769357215DE4FAC081bf1f309aDC325306);
    // 将部署者设置为 owner
    owner = msg.sender;
    // 部署时记录时间
    deploymentTimestamp=block.timestamp;
    //锁仓时长
    lockTime=_lockTime;
}

// 更换 owner
function transferOwnership(address new_owner) public onlyOwner {
    owner = new_owner;
}

// fund 函数：用户向合约转账 ETH
function fund() external payable {
    // 检查转入 ETH 金额（换算成 USD）是否大于最小值
    require(convertETHToUSD(msg.value) >= MINIMUM_VALUE, "Send more ETH");

    //判断是否在锁定期外
    require(block.timestamp<lockTime+deploymentTimestamp,"Windows is closed");

    // 记录用户本次转账金额
    FundsToAmount[msg.sender] = msg.value;

    // 累加用户累计转账金额
    User_Amount[msg.sender] += msg.value;
}

// 从 Chainlink 获取最新 ETH/USD 价格（带8位小数，单位：USD*10^8）
function getChainlinkDataFeedLatestAnswer() public view returns (int) {
    (
        /* uint80 roundId */,
        int256 answer, // 预言机返回的 ETH/USD 价格
        /* uint256 startedAt */,
        /* uint256 updatedAt */,
        /* uint80 answeredInRound */
    ) = dataFeed.latestRoundData();
    return answer;
}

// 将 ETH 金额转换为 USD 金额
function convertETHToUSD(uint256 ETH_Amount) internal view returns (uint256) {
    uint256 ETH_Price = uint256(getChainlinkDataFeedLatestAnswer()); // 当前 ETH/USD 价格，单位 USD*10^8
    // 计算公式：ETH 数量 * ETH 单价 / 10^8（Chainlink 返回带 8 位小数）
    return ETH_Amount * ETH_Price / (10 ** 8);
}

// getFund 函数：当合约总资金 >= getAmount 时，owner 可提取全部资金
function getFund() external windowClosed onlyOwner {
    // 检查合约当前余额（换算成 USD）是否 >= getAmount
    require(convertETHToUSD(address(this).balance) >= getAmount, "Not enough Funds");

    // transfer 方法：转账 ETH，若失败则 revert
    payable(msg.sender).transfer(address(this).balance);

    // send 方法：转账 ETH，返回 bool 表示成功或失败（此处注释掉）
    /*
    bool success = payable(msg.sender).send(address(this).balance);
    require(success, "Fail to Transfer Fund");
    */

    // call 方法：更低层调用，返回 (bool success, bytes memory data)
    // 格式示例: (bool success, bytes memory result) = addr.call{value:value}("");
    bool success;
    (success, ) = payable(msg.sender).call{value:address(this).balance}("");
    require(success, "transfer failed");

    // 清空调用者的 FundsToAmount 记录（这里设计上有争议，因为 getFund 本意是 owner 提取全部资金）
    FundsToAmount[msg.sender] = 0;
}

// refund 函数：当合约总资金 < getAmount 时，用户可申请退款
function refund() external {
    // 检查合约余额（换算成 USD）是否小于目标金额
    require(convertETHToUSD(address(this).balance) < getAmount, "Enough Funds");

    require(block.timestamp>=lockTime+deploymentTimestamp,"Windows is not closed");//需判断合约余额是否达到目标，因此不能直接调用onlyOwner

    // 检查当前用户是否有累计资金记录
    require(User_Amount[msg.sender] != 0, "No Funds to refund");

    // 记录用户退款金额
    uint256 refundAmount = User_Amount[msg.sender];

    // 先将用户余额置零，防止重入攻击
    User_Amount[msg.sender] = 0;
    FundsToAmount[msg.sender] = 0;

    // 调用 call 方法退款给用户
    bool success;
    (success, ) = payable(msg.sender).call{value:refundAmount}("");
    require(success, "transfer failed");
}

//判断是否在锁定期内,
modifier windowClosed(){
    require(block.timestamp>=lockTime+deploymentTimestamp,"Windows is not closed");
    _;//放在require后可节省gas
}

// 确保只有 owner 可以调用
modifier onlyOwner(){
    require(msg.sender==owner, "This function can only called by the owner");
    _; //放在require后可节省gas
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