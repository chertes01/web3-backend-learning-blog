# 📘 Solidity 学习笔记 | 第一课：创建你的第一个智能合约 —— HelloWorld

## 🔹 课程目标

 通过最经典的 Hello World 智能合约，学习：

- ✅ Solidity 基本语法结构
- ✅ 状态变量与函数声明
- ✅ `public` / `internal` / `view` / `pure` 修饰符
- ✅ 部署与调用执行流程

---

## 💻 完整代码

```solidity
// SPDX-License-Identifier: MIT
// 许可证声明，指明此合约的开源协议为 MIT，可以自由复制、使用、修改
pragma solidity ^0.8.20;
// 声明 Solidity 编译器版本，使用 0.8.20 版本及以上兼容版本


// 定义一个名为 HelloWorld 的合约（contract 是 Solidity 中的合约结构体）
contract HelloWorld {


// 声明一个状态变量 strVar，类型为 string，初始值为 "hello world"
// 状态变量会被存储在区块链上
string strVar = "hello world";

// 定义一个公共的只读函数 sayHelloWorld，返回类型是 string
// 使用 view 修饰符表示该函数不会修改区块链上的数据
function sayHelloWorld() public view returns(string memory)
{
    // 调用内部函数 addinfo，将 strVar 作为参数传入
    return addinfo(strVar);
}

// 定义一个公共函数 setHelloWorld，用于修改 strVar 的值
// 参数 newstrVar 是从外部传入的字符串，memory 表示该参数是临时数据
function setHelloWorld(string memory newstrVar) public
{
    // 将传入的新字符串赋值给状态变量 strVar，实现更新操作
    strVar = newstrVar;
}

// 定义一个内部纯函数 addinfo，只在合约内部调用
// pure 修饰符表示该函数既不读取也不修改区块链上的状态
// 接收一个字符串参数 HelloWorldstr，并返回添加说明后的新字符串
function addinfo(string memory HelloWorldstr) internal pure returns(string memory)
{
    // 使用 string.concat 拼接两个字符串
    // 返回形如 "hello world，from Von's Smart Contract" 的信息
    return string.concat(HelloWorldstr, "，from Von's Smart Contract");
}

}
```

---

## 🔍 逐行解析你的第一个智能合约

#### 1️⃣ `// SPDX-License-Identifier: MIT`
> **许可证声明**
> 指明此合约遵循 MIT 开源协议，可自由复制、使用、修改。
> 🔖 推荐每个 Solidity 文件都添加 SPDX 许可证。

#### 2️⃣ `pragma solidity ^0.8.20;`
> **编译器版本声明**
> 告诉编译器使用 `0.8.20` 或更高版本编译（不含 `0.9.0` 及以上），避免版本差异导致报错。

#### 3️⃣ `contract HelloWorld { ... }`
> **合约主体定义**
> `contract` 关键字定义一个智能合约，类似其他语言的 `class` 类。
> 部署到区块链上的是这个合约。

#### 4️⃣ `string strVar = "hello world";`
> **状态变量声明**
> 定义了 `strVar`，类型为 `string`，初始值为 `"hello world"`。
> ✔️ **状态变量** 默认存储在 `storage`（区块链永久存储区），可被本合约内所有函数访问。

#### 5️⃣ `function sayHelloWorld() public view returns(string memory)`
> **只读函数**
>
> | 修饰符 | 作用 |
> |---|---|
> | `public` | 任何外部用户或合约可调用 |
> | `view` | 读取状态变量但不修改 |
> | `pure` | 不读取也不修改状态变量 |
> | `internal`| 只能在合约内部或继承的合约中调用 |
>
> ✔️ **功能**：接收一个字符串参数，拼接固定后缀：`",from Von's Smart Contract"`

---

## 🔧 调用图示（逻辑流程）

1. **部署合约**
    - `strVar` 的初始值被设为 `"hello world"`
    - `owner` 被设为部署者的地址

2. **调用 `sayHelloWorld()`**
    - 调用 `addinfo(strVar)`
    - 返回 `"hello world，from Von's Smart Contract"`

3. **调用 `setHelloWorld("你好 GPT")`**
    - 更新 `strVar` 的值为 `"你好 GPT"`

4. **再次调用 `sayHelloWorld()`**
    - 调用 `addinfo(strVar)`
    - 返回 `"你好 GPT，from Von's Smart Contract"`

---

## ✅ 本课小结

- ✔️ 学会 Solidity 合约基本语法结构（变量、函数、合约）
- ✔️ 理解 **状态变量** 与 数据位置 `memory`
- ✔️ 掌握 `view` 与 `pure` 的区别
- ✔️ 熟悉 **合约部署、调用、返回值** 的完整流程

---

## 📚 下一课预告

在第二课中，我们将学习：

- 如何定义并使用 **结构体 `struct`**
- 如何使用 **数组与映射（`mapping`）**
- 实现简单的 **数据增删查改**

---

## 🎯 完成本课练习

1. 修改 `addinfo` 函数返回内容为 `"xxx, Welcome to Solidity World!"`
2. 在 Remix 部署并调用，体会区块链上的存储读写与返回机制。