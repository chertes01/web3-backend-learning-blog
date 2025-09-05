package main

import "fmt"

func main() {
	// ==========================
	// 整数类型表示
	// ==========================
	// 十六进制表示
	var hex1 uint8 = 0xF // 大写 F
	var hex2 uint8 = 0xf // 小写 f

	// 八进制表示
	var oct1 uint8 = 017  // 传统八进制
	var oct2 uint8 = 0o17 // Go 1.13+ 的八进制表示法
	var oct3 uint8 = 0o17 // 可以重复使用

	// 二进制表示
	var bin1 uint8 = 0b1111 // Go 1.13+ 二进制
	var bin2 uint8 = 0b1111

	// 十进制表示
	var dec uint8 = 15

	fmt.Println("整数表示：", hex1, hex2, oct1, oct2, oct3, bin1, bin2, dec)

	// ==========================
	// 浮点数类型
	// ==========================
	var float1 float32 = 10 // 显式指定 float32
	float2 := 10.0          // 默认 float64
	fmt.Println("float1 == float2? ", float1 == float32(float2))

	// ==========================
	// 复数类型
	// ==========================
	var c1 complex64 = 1 + 2i // 显式 complex64
	c2 := 1 + 2i              // 默认 complex128
	c3 := complex(1, 2)       // 使用 complex() 函数创建复数

	fmt.Println("c1 == complex64(c2)? ", c1 == complex64(c2))
	fmt.Println("c3 = ", c3)

	// 获取复数的实部和虚部
	x := real(c3)
	y := imag(c3)
	fmt.Println("复数 c3 的实部和虚部：", x, y)

	// ==========================
	// 字符串与字节、字符数组转换
	// ==========================
	var str string = "hello"
	var bytes []byte = []byte(str) // 字符串转字节切片
	fmt.Println("字符串转字节: ", bytes)

	var str1 string = string([]byte{97, 98, 99, 100}) // 字节切片转字符串
	fmt.Println("字节转字符串: ", str1)

	var str2 string = "abc，你好，世界！"
	var runes []rune = []rune(str2) // 字符串转 rune 切片（处理 Unicode）
	fmt.Println("字符串转 Rune: ", runes)

	var str3 string = string(runes) // rune 切片转字符串
	fmt.Println("Rune 转字符串: ", str3)

	// ==========================
	// 原生字符串与转义字符串
	// ==========================
	var s1 string = "Hello\nworld!\n" // 转义字符串
	var s2 string = `Hello
	world!
	` // 原生字符串
	fmt.Println("s1 == s2? ", s1 == s2)

	// ==========================
	// 字符串长度与切片
	// ==========================
	var str4 string = "Hello, 世界"
	var bytes2 []byte = []byte(str4) // 字节切片
	var runes2 []rune = []rune(str4) // rune 切片

	fmt.Println("字符串长度 len(str4): ", len(str4)) // 按字节长度
	fmt.Println("字节切片长度 len(bytes2): ", len(bytes2))
	fmt.Println("Rune 切片长度 len(runes2): ", len(runes2)) // 按字符长度

	// 字符串、字节切片、rune 切片的切片操作
	fmt.Println("字符串切片 str4[0:7]: ", str4[0:7]) // 按字节索引，可能会截断中文
	fmt.Println("字节切片切片 string(bytes2[0:7]): ", string(bytes2[0:7]))
	fmt.Println("Rune 切片切片 string(runes2[0:3]): ", string(runes2[0:3])) // 正确处理中文
}
