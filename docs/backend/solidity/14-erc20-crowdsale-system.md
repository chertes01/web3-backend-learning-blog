# 📘 Solidity 学习笔记 | 第十三课：基于 ERC20 的募资铸币系统设计

## 🔹 课程目标

本项目构建了一个去中心化募资铸币系统，其业务逻辑为：

- ✅ 理解 `FundTokenERC20` 与 `FundMe` 两合约的**交互设计**。
- ✅ 掌握 **ERC20 标准合约调用**与**继承**。
- ✅ 透彻理解 `require` 权限判断、**跨合约调用**及 `mint`/`burn` 机制。

---

## 💻 完整代码

为了实现募资铸币系统，我们需要两个核心合约：`FundMe.sol`（负责募资和资金管理）和 `FundTokenERC20.sol`（负责代币的铸造和销毁）。

**`FundMe.sol`**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

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

    // New: Flag to indicate if fundraising was successful and funds withdrawn
    bool public getFundSuccess = false;

    // New: Address of the associated ERC20 token contract
    address public erc20Addr;

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

    // New: Function to set the ERC20 token contract address
    function setErc20Addr(address _erc20Addr) public onlyOwner {
        require(_erc20Addr != address(0), "ERC20 address cannot be zero");
        erc20Addr = _erc20Addr;
    }

    // New: Function to update funder's amount, callable only by the ERC20 contract
    function setFunderToAmount(address _funder, uint256 amountToUpdate) external {
        require(msg.sender == erc20Addr, "Only the associated ERC20 contract can call this function");
        s_addressToAmountFunded[_funder] = amountToUpdate;
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

        // New: Set flag to true after successful withdrawal
        getFundSuccess = true;
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

**`FundTokenERC20.sol`**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import {FundMe} from "./FundMe.sol"; // Assuming FundMe is in a separate file

contract FundTokenERC20 is ERC20 {
    FundMe public fundMe; // Instance of the FundMe contract
    address public immutable i_owner; // Owner of this token contract

    constructor(address fundMeaddr) ERC20("Kincoin", "KC") {
        require(fundMeaddr != address(0), "FundMe address cannot be zero");
        fundMe = FundMe(fundMeaddr);
        i_owner = msg.sender; // Set owner of the token contract
    }

    function mint(uint256 amountToMint) public {
        // Check if user has enough accumulated funds in FundMe
        require(fundMe.s_addressToAmountFunded(msg.sender) >= amountToMint, "You can't mint that much tokens");
        // Check if the fundraising was successful and funds were withdrawn by owner
        require(fundMe.getFundSuccess(), "No fund has been raised yet!");

        // Mint tokens to the caller
        _mint(msg.sender, amountToMint);

        // Update user's accumulated funds in FundMe (deducting minted amount)
        // This is a cross-contract call to FundMe
        fundMe.setFunderToAmount(msg.sender, fundMe.s_addressToAmountFunded(msg.sender) - amountToMint);
    }

    function claim(uint256 amountToClaim) public {
        // Check if caller has enough tokens to burn
        require(balanceOf(msg.sender) >= amountToClaim, "You don't have enough tokens");
        // Burn tokens from the caller
        _burn(msg.sender, amountToClaim);
    }
}
```

---

## 🎯 核心业务场景

本项目构建了一个去中心化募资铸币系统，其业务逻辑为：

1.  用户向 `FundMe` 合约转入 ETH 进行募资。
2.  `FundMe` 合约的 `owner` 在锁仓期结束后，且募资成功（达到目标或决定提现）后，提走 ETH。
3.  募资成功提取后，用户可调用 `FundTokenERC20` 合约的 `mint` 方法，根据其在 `FundMe` 中的捐款额兑换 ERC20 通证。
4.  用户也可通过 `claim` 方法销毁（`burn`）部分通证。

---

## 🔍 新知识点与代码详解

#### ✨ 1. `FundTokenERC20` 合约

##### 💡 继承 `ERC20`
```solidity
import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";
contract FundTokenERC20 is ERC20 { ... }
```
✔️ **解读**：引入 OpenZeppelin 的 `ERC20` 标准库，通过 `is ERC20` 实现继承，从而自动获得所有 ERC20 的基础功能（如 `totalSupply`、`balanceOf`、`transfer` 等）。

##### 💡 构造函数
```solidity
constructor(address fundMeaddr) ERC20("Kincoin", "KC") {
    require(fundMeaddr != address(0), "FundMe address cannot be zero");
    fundMe = FundMe(fundMeaddr);
    i_owner = msg.sender; // Set owner of the token contract
}
```
✔️ **解读**：
-   `ERC20("Kincoin", "KC")`：这是调用父合约 `ERC20` 的构造函数，用于初始化 Token 的名称 (`Kincoin`) 和简称 (`KC`)。
-   `fundMeaddr`：传入 `FundMe` 合约的地址，并将其转换为 `FundMe` 类型实例 `fundMe`，以便进行跨合约调用。
-   `i_owner = msg.sender;`：设置 `FundTokenERC20` 合约的部署者为该 Token 合约的 `owner`。

##### 💡 `mint` 函数
```solidity
function mint(uint256 amountToMint) public {
    require(fundMe.s_addressToAmountFunded(msg.sender) >= amountToMint, "You can't mint that much tokens");
    require(fundMe.getFundSuccess(), "No fund has been raised yet!");

    _mint(msg.sender, amountToMint);

    fundMe.setFunderToAmount(msg.sender, fundMe.s_addressToAmountFunded(msg.sender) - amountToMint);
}
```
✅ **功能**：
-   **权限检查**：
    -   `require(fundMe.s_addressToAmountFunded(msg.sender) >= amountToMint, ...)`：检查用户在 `FundMe` 合约中的累计募资金额是否大于或等于其希望铸造的 Token 数量。这确保了用户只能根据其贡献来铸币。
    -   `require(fundMe.getFundSuccess(), ...)`：检查 `FundMe` 合约的募资是否已成功提取（即 `getFundSuccess` 为 `true`）。这通常意味着只有在众筹目标达成且资金被 `owner` 提走后，用户才能兑换 Token。
-   **铸造 Token**：`_mint(msg.sender, amountToMint)` 调用 OpenZeppelin `ERC20` 库提供的内部函数，为 `msg.sender` 铸造指定数量的 Token，并增加总供应量。
-   **更新 `FundMe` 记录**：`fundMe.setFunderToAmount(msg.sender, fundMe.s_addressToAmountFunded(msg.sender) - amountToMint)` 进行跨合约调用，更新用户在 `FundMe` 中的募资余额，扣除已用于铸币的部分。

##### 💡 `claim` 函数
```solidity
function claim(uint256 amountToClaim) public {
    require(balanceOf(msg.sender) >= amountToClaim, "You don't have enough tokens");
    _burn(msg.sender, amountToClaim);
}
```
✅ **功能**：
-   **余额检查**：`require(balanceOf(msg.sender) >= amountToClaim, ...)` 确保用户拥有足够的 Token 来销毁。
-   **销毁 Token**：`_burn(msg.sender, amountToClaim)` 调用 OpenZeppelin `ERC20` 库提供的内部函数，销毁 `msg.sender` 指定数量的 Token，并减少总供应量。
-   **常见用例**：`claim` 函数通常用于实现抵押赎回、销毁投票权或链上销毁机制等。

#### ✨ 2. `FundMe` 合约新增逻辑

##### 💡 `getFundSuccess` 状态变量
```solidity
bool public getFundSuccess = false;
```
✔️ **解读**：这是一个公共状态变量，用于标记 `FundMe` 合约的募资是否已成功提取。`FundTokenERC20` 合约会查询此变量来验证用户铸币的权限。

##### 💡 `setFunderToAmount` 函数
```solidity
function setFunderToAmount(address _funder, uint256 amountToUpdate) external {
    require(msg.sender == erc20Addr, "Only the associated ERC20 contract can call this function");
    s_addressToAmountFunded[_funder] = amountToUpdate;
}
```
✅ **功能**：
-   **权限控制**：`require(msg.sender == erc20Addr, ...)` 确保此函数只能由与 `FundMe` 关联的 `FundTokenERC20` 合约调用。
-   **更新余额**：用于更新指定 `_funder` 在 `FundMe` 中的募资余额 `s_addressToAmountFunded`。这在用户铸造 Token 后，从其募资余额中扣除相应金额时使用。

##### 💡 `erc20Addr` 变量与 `setErc20Addr` 函数
```solidity
address public erc20Addr;
function setErc20Addr(address _erc20Addr) public onlyOwner {
    require(_erc20Addr != address(0), "ERC20 address cannot be zero");
    erc20Addr = _erc20Addr;
}
```
✔️ **解读**：
-   `erc20Addr`：存储与 `FundMe` 关联的 `FundTokenERC20` 合约的地址。
-   `setErc20Addr`：一个 `onlyOwner` 函数，允许 `FundMe` 的 `owner` 设置或更新 `erc20Addr`。这是实现跨合约权限控制的关键一步，确保 `FundMe` 知道哪个 `ERC20` 合约有权调用其 `setFunderToAmount` 函数。

#### 🎯 3. 关键知识点拓展

##### 🔑 3.1 ERC20 标准简述
ERC20 通证通常用于：交易流通、DApp 积分、链上票据凭证。

| 函数 | 功能 |
| :--- | :--- |
| `totalSupply()` | 返回总发行量 |
| `balanceOf(address)` | 查询地址余额 |
| `transfer(address,uint256)` | 转账 |
| `_mint(address,uint256)` | 内部铸币（OpenZeppelin 提供） |
| `_burn(address,uint256)` | 内部销毁（OpenZeppelin 提供） |

##### 🔑 3.2 `require` 与跨合约调用
```solidity
require(fundMe.s_addressToAmountFunded(msg.sender) >= amountToMint, "...");
```
✔️ **解读**：跨合约调用与自身函数调用在语法上相似，但需要先在构造函数中实例化目标合约（`fundMe = FundMe(fundMeaddr);`），然后通过实例变量 (`fundMe.`) 调用其公共函数。

##### 🔑 3.3 铸币 `mint` 与销毁 `burn`
✔️ **`_mint`**：增加指定地址的 Token 余额，并增加 `totalSupply`。
✔️ **`_burn`**：减少指定地址的 Token 余额，并减少 `totalSupply`。

**常见 use-case**：
-   **质押挖矿产币**：用户质押资产，合约 `_mint` 新 Token 作为奖励。
-   **赎回销毁机制**：用户将 Token 销毁，以赎回其他资产（如 USDT 赎回 USD）。
-   **治理代币**：通过销毁代币来投票或执行特定操作。

##### 🔑 3.4 构造函数中的父合约初始化
```solidity
constructor(address fundMeaddr) ERC20("Kincoin", "KC") { ... }
```
✔️ **解读**：`ERC20("Kincoin", "KC")` 是在子合约 `FundTokenERC20` 的构造函数中，调用其父合约 `ERC20` 的构造函数。这种方式与其他面向对象语言（如 Java 的 `super()`）初始化父类的方式一致。

---

## 🔧 调用图示（逻辑流程）

1.  **部署 `FundMe` 合约**：`FundMe` 部署，`i_owner` 被设置。
2.  **部署 `FundTokenERC20` 合约**：`FundTokenERC20` 部署，构造函数接收 `FundMe` 地址，并初始化 `fundMe` 实例。`i_owner` 被设置。
3.  **关联合约**：`FundMe` 的 `owner` 调用 `FundMe.setErc20Addr(FundTokenERC20_address)`，将两个合约关联起来。
4.  **用户捐款**：用户调用 `FundMe.fund()` 并发送 ETH。`FundMe` 记录用户捐款额 `s_addressToAmountFunded[user]`。
5.  **Owner 提现**：`FundMe` 的 `owner` 在锁仓期结束后调用 `FundMe.withdraw()`。`FundMe` 将 ETH 提走，并设置 `getFundSuccess = true`。
6.  **用户铸币**：用户调用 `FundTokenERC20.mint(amount)`。
    -   `FundTokenERC20` 检查：
        -   用户在 `FundMe` 中的捐款额是否足够 (`fundMe.s_addressToAmountFunded(msg.sender) >= amount`)。
        -   `FundMe` 是否已成功提现 (`fundMe.getFundSuccess()`)。
    -   如果检查通过：
        -   `FundTokenERC20` 调用 `_mint(msg.sender, amount)` 为用户铸造 Token。
        -   `FundTokenERC20` 跨合约调用 `fundMe.setFunderToAmount(msg.sender, ...)`，从用户在 `FundMe` 的捐款记录中扣除已用于铸币的金额。
7.  **用户销毁 Token**：用户调用 `FundTokenERC20.claim(amount)`。
    -   `FundTokenERC20` 检查用户是否有足够的 Token。
    -   如果检查通过，`FundTokenERC20` 调用 `_burn(msg.sender, amount)` 销毁 Token。

---

## ✅ 本课总结

- ✅ 实现了 `FundMe` 募资合约与 `FundTokenERC20` 铸币合约的**跨合约交互**。
- ✅ 理解了 **ERC20 标准**函数调用与继承结构。
- ✅ 掌握了 `mint` / `burn` 核心逻辑及其在代币经济模型中的作用。
- ✅ 学习了如何通过 `require` 和 `msg.sender` 进行**跨合约权限判断**。

---

## 🎯 练习拓展

1.  **铸币比例**：修改 `FundTokenERC20.mint` 函数，使其铸造的 Token 数量与用户在 `FundMe` 中的捐款额之间存在一个固定的兑换比例（例如 1 ETH 捐款 = 100 KC Token）。
2.  **Owner 铸币**：为 `FundTokenERC20` 添加一个 `ownerMint(address recipient, uint256 amount)` 函数，只允许 `FundTokenERC20` 的 `owner` 为任意地址铸造 Token。
3.  **募资目标与 Token 兑付**：思考如何设计募资金额与 Token 铸造比例，才能确保兑付安全，例如，如果募资未达目标，用户是否可以退款而不是铸币？

