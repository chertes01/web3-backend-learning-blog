# 📘 Solidity 学习笔记 | 第二课：结构体、数组与 msg.sender

## 🔹 课程目标

通过一个可记录信息的 HelloWorld 合约，学习：

- ✅ `struct`（结构体）的定义与应用
- ✅ 使用 `array` (数组) 存储多个结构体
- ✅ 掌握 `msg.sender` 记录调用者地址
- ✅ 理解 `for` 循环查询与 `push` 新增元素

---

## 💻 完整代码

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 这是一个扩展版的 HelloWorld 合约
// 它可以存储不同用户设置的 "hello" 信息
contract HelloWorld {


// 状态变量 strVar，初始值为 "hello world"
string strVar = "hello world";

// ✅ 定义结构体 Info，包含 phrase(短语)、id(编号)、addr(调用者地址)
struct Info {
    string phrase;
    uint256 id;
    address addr;
}

// ✅ 定义结构体数组 infos，用于保存所有 Info
Info[] infos;

// ✅ 查询函数 sayHelloWorld
// 根据传入的 id 在 infos 中查询 phrase，否则返回默认 strVar
function sayHelloWorld(uint256 _id) public view returns(string memory) {
    // for 循环遍历 infos 数组
    for (uint256 i = 0; i < infos.length; i++) {
        if (infos[i].id == _id) {
            // 如果找到匹配 id，调用 addinfo 拼接后返回
            return addinfo(infos[i].phrase);
        }
    }
    // 如果未找到，返回 strVar 拼接结果
    return addinfo(strVar);
}

// ✅ 设置函数 setHelloWorld
// 将外部传入的 phrase 和 id 封装为 Info 结构体，并添加到 infos 数组
function setHelloWorld(string memory newstrVar, uint _id) public {
    // 创建 Info 结构体，msg.sender 为调用者地址
    Info memory info = Info(newstrVar, _id, msg.sender);
    // 将 info 添加到 infos 数组
    infos.push(info);
}

// ✅ 内部函数 addinfo
// 拼接传入字符串和固定后缀
function addinfo(string memory HelloWorldstr) internal pure returns(string memory) {
    return string.concat(HelloWorldstr, ",from Von's Smart Contract");
}

}
```

---

## 🔍 新知识点解析

#### 1️⃣ `struct`（结构体）
> **用途**：自定义复合数据类型，将不同变量封装为一个整体，类似于其他语言中的 `struct` 或 `object`。

```solidity
struct Info {
    string phrase;
    uint256 id;
    address addr;
}
```
- `phrase`: 字符串短语。
- `id`: 无符号整数，用作唯一标识。
- `addr`: 以太坊地址（20字节），用于记录调用者。

#### 2️⃣ 数组 (Array)
> **定义**：用于存储一组相同类型元素的集合。

| 类型 | 声明方式 | 特点 |
| :--- | :--- | :--- |
| **动态数组** | `Info[] infos;` | 长度不固定，初始为 0，可使用 `push` 动态增长。 |
| **定长数组** | `Info[5] infos;` | 长度固定为 5，在声明时确定，无法增加或减少元素。|

本合约中使用的是**动态数组**，更具灵活性。

#### 3️⃣ `msg.sender`
> **定义**：Solidity 内置的全局变量，它代表**调用当前函数的账户地址**。

---

## 🔧 调用图示（逻辑流程）

1.  **部署合约**
    - `strVar` 的初始值被设为 `"hello world"`。
    - `infos` 数组被初始化为空 `[]`。

2.  **调用 `setHelloWorld("你好", 1)`**
    - 创建一个 `Info` 内存变量：`Info(phrase="你好", id=1, addr=调用者地址)`。
    - 通过 `infos.push()` 将该结构体存入 `infos` 状态变量数组中。

3.  **调用 `sayHelloWorld(1)`**
    - `for` 循环遍历 `infos` 数组，找到 `id` 为 `1` 的匹配项。
    - 返回结果：`"你好,from Von's Smart Contract"`。

---

## ✅ 本课总结

- ✅ 学会了如何定义和实例化**结构体 `struct`**。
- ✅ 掌握了如何使用**动态数组 `array`** 来存储和管理一组复杂数据。
- ✅ 理解了 `msg.sender` 的含义，并用它来追踪**函数调用者**。
- ✅ 实践了通过 `for` 循环**遍历数组**进行数据查询。

---

## 🎯 练习拓展

1.  修改 `sayHelloWorld` 函数，如果找到对应 `id`，让它同时返回 `phrase` 和 `addr`。
2.  在 `setHelloWorld` 函数中添加 `require` 语句，确保新添加的 `_id` 不与数组中已有的 `id` 重复。
3.  （挑战）为合约添加一个 `deleteInfo(uint _id)` 函数，用于删除指定 `id` 的元素。