# 📘 Solidity 学习笔记 | 第四课：合约工厂模式 (Factory Pattern)

## 🔹 课程目标

在本节课，我们将学习如何通过**工厂模式 (Factory Pattern)** 批量部署和管理合约，实现更模块化、工程化的开发。

- ✅ 理解什么是**合约工厂模式**。
- ✅ 学会使用 `new` 关键字**批量部署**子合约。
- ✅ 掌握如何通过工厂合约**集中管理和调用**子合约实例。
- ✅ 学习使用 `import` 导入其他合约文件。

---

## 💻 完整代码

首先，我们假设 `HelloWorld` 合约存在于一个名为 `HelloWorld.sol` 的独立文件中。然后，我们创建工厂合约 `HelloWorldFactory.sol`。

**`HelloWorldFactory.sol`**
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

// 1. 引入同级目录下 HelloWorld.sol 文件中定义的 HelloWorld 合约
import {HelloWorld} from "./HelloWorld.sol";

// 工厂合约定义
contract HelloWorldFactory {
    // 2. 定义一个 HelloWorld 类型的动态数组，用于保存所有已部署的 HelloWorld 合约实例。
    HelloWorld[] public hws;

    // 3. 创建并部署一个新的 HelloWorld 子合约
    function CreateHelloWorld() public {
        // 使用 new 关键字在区块链上部署一个新的 HelloWorld 合约
        HelloWorld hw = new HelloWorld();
        // 将新部署的 HelloWorld 合约实例存入 hws 动态数组，实现批量管理
        hws.push(hw);
    }

    // 4. 根据索引获取 HelloWorld 合约实例的地址
    function GetHelloWorldByIndex(uint256 _index) public view returns(HelloWorld) {
        // 返回 hws 数组中指定索引的 HelloWorld 合约实例
        return hws[_index];
    }

    // 5. 通过工厂调用指定子合约的 sayHelloWorld 函数
    function CallHelloFromFactory(uint256 _index, uint256 _id) public view returns(string memory) {
        // 调用第 _index 个 HelloWorld 合约实例的 sayHelloWorld 函数，并返回结果
        return hws[_index].sayHelloWorld(_id);
    }

    // 6. 通过工厂调用指定子合约的 setHelloWorld 函数
    function SetHelloWorldFactory(uint256 index, string memory newstrVar, uint _id) public {
        // 调用第 index 个 HelloWorld 合约实例的 setHelloWorld 函数，实现向子合约写入数据
        hws[index].setHelloWorld(newstrVar, _id);
    }
}
```

---

## 🔍 新知识点解析

#### 1️⃣ `import {HelloWorld} from "./HelloWorld.sol";`
> **文件导入**：`import` 语句允许你引入其他 Solidity 文件中定义的合约、库或接口。这使得代码可以被复用和模块化，是构建复杂项目的基石。

#### 2️⃣ `new HelloWorld()`
> **合约创建**：`new` 是一个特殊的关键字，用于在区块链上**部署一个新的合约实例**。
> - `new HelloWorld()` 会创建一个新的 `HelloWorld` 合约。
> - 这个操作会返回新创建合约的**地址**，我们可以将其存储在一个 `HelloWorld` 类型的变量中。
> - 每次调用 `new` 都会在链上生成一个全新的、独立的合约，拥有自己的存储空间，并消耗 Gas。

#### 3️⃣ `hws[_index].sayHelloWorld(_id)`
> **合约间交互**：一旦你拥有了另一个合约的实例（地址），你就可以像调用普通对象的方法一样，调用它的 `public` 或 `external` 函数。
> - `hws[_index]` 从数组中获取一个 `HelloWorld` 合约实例。
> - `.sayHelloWorld(_id)` 紧接着调用该实例上的函数。

---

## 🔧 调用图示（逻辑流程）

1.  **部署 `HelloWorldFactory` 合约**
    - 一个 `HelloWorldFactory` 实例被创建在区块链上。
    - 其内部的 `hws` 数组此时为空。

2.  **调用 `CreateHelloWorld()`**
    - 工厂合约执行 `new HelloWorld()`，在链上**创建了一个全新的 `HelloWorld` 子合约**。
    - 这个新的子合约地址被 `push` 到工厂合约的 `hws` 数组中。
    - 假设这是第一次调用，`hws[0]` 现在就指向了这个新的子合约。

3.  **调用 `SetHelloWorldFactory(0, "你好工厂", 1)`**
    - 工厂合约找到 `hws[0]` 对应的子合约。
    - 它代理调用该子合约的 `setHelloWorld("你好工厂", 1)` 函数，将数据写入**子合约的存储空间**。

4.  **调用 `CallHelloFromFactory(0, 1)`**
    - 工厂合约再次找到 `hws[0]` 对应的子合约。
    - 它代理调用该子合约的 `sayHelloWorld(1)` 函数，并返回从**子合约**读取到的数据。

---

## ✅ 本课总结

- ✅ **工厂模式**是一种强大的设计模式，它允许一个合约（工厂）动态地创建和管理其他合约（产品）。
- ✅ 使用 `import` 可以组织和复用代码，实现项目模块化。
- ✅ `new` 关键字是动态部署新合约实例的核心。
- ✅ 通过持有子合约的地址/实例，一个合约可以方便地与另一个合约进行交互（函数调用）。

---

## 🎯 练习拓展

1.  在 `CreateHelloWorld` 函数中添加一个 `event`，每当新合约被创建时，就发出一个事件，记录新合约的地址和创建者。
2.  添加一个 `getHwsLength()` 函数，用于返回已创建的 `HelloWorld` 合约的总数。
3.  （挑战）修改 `CreateHelloWorld` 函数，增加权限控制，只允许工厂合约的部署者（owner）创建新的子合约。