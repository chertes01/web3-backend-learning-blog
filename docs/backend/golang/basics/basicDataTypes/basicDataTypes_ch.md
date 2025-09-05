# Go Basic Data Types 学习笔记

---

## 1. 整数类型（Integer）

Go 支持多种进制表示整数，包括十进制、十六进制、八进制和二进制。

```go
// 十六进制表示
var hex1 uint8 = 0xF   // 大写 F
var hex2 uint8 = 0xf   // 小写 f

// 八进制表示
var oct1 uint8 = 017    // 传统八进制
var oct2 uint8 = 0o17   // Go 1.13+ 八进制表示法
var oct3 uint8 = 0o17   // 可以重复使用

// 二进制表示
var bin1 uint8 = 0b1111 // Go 1.13+ 二进制
var bin2 uint8 = 0b1111

// 十进制表示
var dec uint8 = 15

fmt.Println("整数表示：", hex1, hex2, oct1, oct2, oct3, bin1, bin2, dec)
```

**说明：**
- Go 支持 uint8、uint16、uint32、uint64 等多种整数类型。
- 二进制、八进制和十六进制常用于底层开发或与硬件打交道时。

---

## 2. 浮点数类型（Float）

Go 支持 float32 和 float64 类型。

```go
var float1 float32 = 10      // 显式指定 float32
float2 := 10.0               // 默认 float64
fmt.Println("float1 == float2? ", float1 == float32(float2))
```

**说明：**
- 默认浮点数类型是 float64。
- 浮点数在比较时需要注意类型转换，避免类型不同导致比较结果错误。

---

## 3. 复数类型（Complex）

Go 支持复数类型 complex64 和 complex128，可通过 complex() 函数创建。

```go
var c1 complex64 = 1 + 2i      // 显式 complex64
c2 := 1 + 2i                   // 默认 complex128
c3 := complex(1, 2)            // 使用 complex() 函数

fmt.Println("c1 == complex64(c2)? ", c1 == complex64(c2))
fmt.Println("c3 = ", c3)

x := real(c3)  // 获取实部
y := imag(c3)  // 获取虚部
fmt.Println("复数 c3 的实部和虚部：", x, y)
```

**说明：**
- real() 获取复数的实部，imag() 获取虚部。
- 复数常用于信号处理、科学计算等领域。

---

## 4. 字符串与字节切片（String & []byte）

**字符串转字节切片：**
```go
var str string = "hello"
var bytes []byte = []byte(str)
fmt.Println("字符串转字节: ", bytes)
```

**字节切片转字符串：**
```go
var str1 string = string([]byte{97, 98, 99, 100})
fmt.Println("字节转字符串: ", str1)
```

---

## 5. 字符串与 Unicode 字符（Rune）

**字符串转 Rune 切片：**
```go
var str2 string = "abc，你好，世界！"
var runes []rune = []rune(str2)
fmt.Println("字符串转 Rune: ", runes)
```

**Rune 切片转字符串：**
```go
var str3 string = string(runes)
fmt.Println("Rune 转字符串: ", str3)
```

**说明：**
- rune 类型是 int32 的别名，用于表示 Unicode 字符。
- 中文或其他多字节字符需要使用 rune 切片正确处理。

---

## 6. 原生字符串与转义字符串

```go
var s1 string = "Hello\nworld!\n" // 转义字符串
var s2 string = `Hello
    world!
    ` // 原生字符串
fmt.Println("s1 == s2? ", s1 == s2)
```

**说明：**
- 转义字符串支持 \n、\t 等特殊字符。
- 原生字符串用反引号 ` 包裹，保留换行和空格。

---

## 7. 字符串长度与切片操作

```go
var str4 string = "Hello, 世界"
var bytes2 []byte = []byte(str4)
var runes2 []rune = []rune(str4)

fmt.Println("字符串长度 len(str4): ", len(str4))         // 按字节长度
fmt.Println("字节切片长度 len(bytes2): ", len(bytes2))
fmt.Println("Rune 切片长度 len(runes2): ", len(runes2)) // 按字符长度

// 切片操作
fmt.Println("字符串切片 str4[0:7]: ", str4[0:7])
fmt.Println("字节切片切片 string(bytes2[0:7]): ", string(bytes2[0:7]))
fmt.Println("Rune 切片切片 string(runes2[0:3]): ", string(runes2[0:3])) // 正确处理中文
```

**说明：**
- len(str) 返回的是字节数，而不是字符数。
- 中文字符占 3 个字节（UTF-8），所以需要用 rune 切片切片才能正确处理。

---

## 总结

- Go 数据类型多样，整数、浮点、复数、字符串各有特点。
- 字符串和 Unicode 字符处理需要区分字节和字符。
- 类型转换在比较或运算时很重要。
- 原生字符串和转义字符串可根据需求选择使用。