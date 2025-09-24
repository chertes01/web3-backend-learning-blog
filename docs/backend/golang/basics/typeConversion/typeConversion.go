package main

import "fmt"

func main() {
	str := "hello, 123, 你好"
	var bytes []byte = []byte(str)
	var runes []rune = []rune(str)
	fmt.Printf("bytes 的值为: %v \n", bytes)
	fmt.Printf("runes 的值为: %v \n", runes)

	str2 := string(bytes)
	str3 := string(runes)
	fmt.Printf("str2 的值为: %v \n", str2)
	fmt.Printf("str3 的值为: %v \n", str3)
}
