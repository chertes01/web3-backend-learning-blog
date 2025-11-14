# 📘 Solidity 学习笔记 | 第三课：使用 Mapping 优化查询

## 🔹 课程目标

在上一课的基础上，引入 `mapping` 来优化合约，学习：

- ✅ `mapping` 的定义与使用场景
- ✅ 比较 **数组遍历查询** vs. **mapping 直接查询** 的效率差异
- ✅ 掌握 `address(0x0)` 在 `mapping` 中的妙用
- ✅ 理解代码优化对 Gas 消耗的影响

---

## 💻 优化后的完整代码

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract HelloWorld {

    string strVar = "hello world";

    struct Info {
        string phrase;
        uint256 id;
        address addr;
    }

    // 虽然 mapping 效率更高，但我们仍保留数组，以便未来需要遍历所有数据
    Info[] infos;

    // ✅ 新增 mapping 定义，用于通过 id 快速查找 Info
    mapping (uint256 => Info) public infoMapping;

    // ✅ 查询函数已优化
    // 直接使用 mapping 查询，替代了 for 循环
    function sayHelloWorld(uint256 _id) public view returns(string memory) {
        // 通过检查 addr 是否为零地址来判断 _id 是否存在
        if (infoMapping[_id].addr == address(0x0)) {
            // 如果不存在，返回默认 strVar
            return addinfo(strVar);
        } else {
            // 如果存在，直接返回 mapping 中存储的 phrase
            return addinfo(infoMapping[_id].phrase);
        }
    }

    // ✅ 设置函数已优化
    // 现在同时向 mapping 和数组中写入数据
    function setHelloWorld(string memory newstrVar, uint _id) public {
        // 检查 id 是否已被使用，防止数据被覆盖
        require(infoMapping[_id].addr == address(0x0), "ID already exists.");

        Info memory info = Info(newstrVar, _id, msg.sender);
        
        // 同时写入 mapping 和 array
        infoMapping[_id] = info;
        infos.push(info);
    }

    // 内部函数，无变化
    function addinfo(string memory HelloWorldstr) internal pure returns(string memory) {
        return string.concat(HelloWorldstr, ",from Von's Smart Contract");
    }

}
```

---

## 🔍 新知识点解析

#### 1️⃣ `mapping` (映射)
> **定义**：`mapping` 是 Solidity 中的哈希表（或字典）结构，用于存储键值对。

```solidity
mapping (uint256 => Info) public infoMapping;
```
这行代码定义了一个名为 `infoMapping` 的 `public` 映射。
- **Key (键)**：`uint256` 类型，我们用它来存储 `id`。
- **Value (值)**：`Info` 结构体，与 `id` 对应的数据。

#### 2️⃣ `address(0x0)` (零地址)
> **定义**：`address(0x0)` 是一个特殊的、全为零的以太坊地址，通常表示一个“空”或“不存在”的地址。

在 `mapping` 中，如果一个 `key` 从未被赋值，那么它对应的 `value` 会是该类型的**默认零值**。
- 对于 `uint`，零值是 `0`。
- 对于 `bool`，零值是 `false`。
- 对于 `address`，零值是 `address(0x0)`。
- 对于 `struct`，其所有成员都是零值。

我们利用这个特性，通过检查 `infoMapping[_id].addr == address(0x0)` 来高效地判断一个 `id` 是否已经被使用。

---

## 🚀 数组查询 vs. Mapping 查询 (性能与 Gas 对比)

| 特性 | 数组遍历查询 (上一课) | Mapping 直接查询 (本课) |
| :--- | :--- | :--- |
| **查询方式** | `for` 循环遍历整个数组 | `infoMapping[_id]` 直接访问 |
| **时间复杂度** | O(N)，查询时间随数组增大而线性增加 | O(1)，查询时间恒定，与数据量无关 |
| **Gas 消耗** | **高**。每次循环都是一次操作，数组越大，Gas 越多。 | **极低**。无论存储多少数据，查询 Gas 成本几乎不变。 |
| **适用场景** | 需要遍历所有数据时。 | 需要通过唯一标识符快速查找单个数据时。 |

**结论**：在需要通过 ID 快速查找的场景下，`mapping` 是比数组遍历**高效得多**的解决方案。

---

## ✅ 本课总结

- ✅ 掌握了 `mapping` 的定义和作为高效键值存储的用法。
- ✅ 理解了 `mapping` 在查询效率和 Gas 消耗上相对于数组遍历的巨大优势。
- ✅ 学会了使用 `address(0x0)` 来检查 `mapping` 中某个键是否存在。
- ✅ 在 `setHelloWorld` 中加入了 `require` 语句，增强了代码的健壮性，防止数据被意外覆盖。

---

## 🎯 练习拓展

1.  为合约添加一个 `updatePhrase(uint256 _id, string memory _newPhrase)` 函数，允许信息的创建者 (`msg.sender == infoMapping[_id].addr`) 修改 `phrase`。
2.  思考一下：为什么我们在使用了 `mapping` 之后，仍然保留了 `infos` 数组？它在什么场景下会有用？