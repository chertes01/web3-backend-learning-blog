package main

import (
	"fmt"
)

// 定义全局变量，包括数组和切片
var (
	array  = [...]string{"zero", "one", "two", "three", "four", "five", "six"} // 固定长度数组
	slice1 = []int{1, 2, 3, 4}                                                 // 初始化切片
	slice2 = make([]int, 4, 5)                                                 // 使用 make 创建切片，长度为4，容量为5
	slice3 = array[3:7]                                                        // 从数组中切片，范围为 [3:7)
)

// 打印切片的长度和容量
func lengthOfslice() {
	fmt.Println("array:", array)
	fmt.Printf("slice1:%v ,length:%d ,capacity:%d\n", slice1, len(slice1), cap(slice1))
	fmt.Printf("slice2:%v ,length:%d ,capacity:%d\n", slice2, len(slice2), cap(slice2))
	fmt.Printf("slice3:%v ,length:%d ,capacity:%d\n", slice3, len(slice3), cap(slice3))
}

// 访问和修改切片
func accessModification() {
	s1 := []int{5, 4, 3, 2, 1}              // 初始化切片
	fmt.Println("s1:", s1[:3])              // 截取前3个元素
	s1 = append([]int{10, 9, 8}, s1[3:]...) // 在切片前插入新元素
	fmt.Println("s1 after modification:", s1)
}

// 遍历切片并打印索引和值
func rangeSlice(s []int) {
	for v := range s {
		fmt.Println(v, s[v])
	}
}

// 检查 nil 切片与空切片的区别
func nilSlice() {
	var nilSlice []int // 声明 nil 切片
	s3 := []int{}      // 初始化空切片
	fmt.Printf("length of nilSlice:%d,length of s3:%d\n", len(nilSlice), len(s3))
	fmt.Printf("capacity of nilSlice:%d , capacity of s3:%d\n", cap(nilSlice), cap(s3))
	fmt.Println(nilSlice == nil) // true，nil 切片
	fmt.Println(s3 == nil)       // false，空切片
}

// 使用 append 添加元素到切片
func useAppend() {
	s := []string{}        // 初始化空切片
	s = append(s, "hello") // 添加一个元素
	fmt.Println("append 1 s:", s)
	s = append(s, "world", "wellcome") // 添加多个元素
	fmt.Println("append 2 s:", s)
	s = append(s[:3], []string{"to", "go", "language"}...) // 在切片中间插入元素
	fmt.Println("append 3 s:", s)
}

// 在切片中插入元素
func insertSlice() {
	s5 := array[0:7]                                           // 从数组创建切片
	s5 = append(s5[0:2], append([]string{"3"}, s5[2:7]...)...) // 在索引 2 处插入元素
	fmt.Println(s5)
}

// 从切片中删除元素
func deleteSlice() {
	s6 := array[0:7]                 // 从数组创建切片
	s6 = append(s6[0:3], s6[4:7]...) // 删除索引 3 的元素
	fmt.Println(s6)
}

func main() {
	lengthOfslice()      // 打印切片长度和容量
	accessModification() // 访问和修改切片
	rangeSlice(slice1)   // 遍历切片
	nilSlice()           // 检查 nil 切片与空切片
	useAppend()          // 使用 append 添加元素
	insertSlice()        // 在切片中插入元素
	deleteSlice()        // 从切片中删除元素
}
