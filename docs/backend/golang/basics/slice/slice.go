package main

import (
	"fmt"
)

var (
	array  = [...]string{"zero", "one", "two", "three", "four", "five", "six"}
	slice1 = []int{1, 2, 3, 4}
	slice2 = make([]int, 4, 5)
	slice3 = array[3:7]
)

func lengthOfslice() {
	fmt.Println("array:", array)
	fmt.Printf("slice1:%v ,length:%d ,capacity:%d\n", slice1, len(slice1), cap(slice1))
	fmt.Printf("slice2:%v ,length:%d ,capacity:%d\n", slice2, len(slice2), cap(slice2))
	fmt.Printf("slice3:%v ,length:%d ,capacity:%d\n", slice3, len(slice3), cap(slice3))
}

func accessModification() {
	s1 := []int{5, 4, 3, 2, 1}
	fmt.Println("s1:", s1[:3])
	s1 = append([]int{10, 9, 8}, s1[3:]...)
	fmt.Println("s1 after modification:", s1)
}

func rangeSlice(s []int) {
	for v := range s {
		fmt.Println(v, s[v])
	}
}

func nilSlice() {
	var nilSlice []int
	s3 := []int{}
	fmt.Printf("length of nilSlice:%d,length of s3:%d\n", len(nilSlice), len(s3))
	fmt.Printf("capacity of nilSlice:%d , capacity of s3:%d\n", cap(nilSlice), cap(s3))
	fmt.Println(nilSlice == nil) // true
	fmt.Println(s3 == nil)       // false
}

func useAppend() {
	s := []string{}
	s = append(s, "hello")
	fmt.Println("append 1 s:", s)
	s = append(s, "world", "wellcome")
	fmt.Println("append 2 s:", s)
	s = append(s[:3], []string{"to", "go", "language"}...)
	fmt.Println("append 3 s:", s)
}

func insertSlice() {
	s5 := array[0:7]
	s5 = append(s5[0:2], append([]string{"3"}, s5[2:7]...)...)
	fmt.Println(s5)
}

func deleteSlice() {
	s6 := array[0:7]
	s6 = append(s6[0:3], s6[4:7]...)
	fmt.Println(s6)
}

func main() {
	lengthOfslice()
	accessModification()
	rangeSlice(slice1)
	nilSlice()
	useAppend()
	insertSlice()
	deleteSlice()
}
