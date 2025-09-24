/*
给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。可以使用双指针法，
一个慢指针 i 用于记录不重复元素的位置，一个快指针 j 用于遍历数组，当 nums[i] 与 nums[j] 不相等时，
将 nums[j] 赋值给 nums[i + 1]，并将 i 后移一位。
*/
package main

import "fmt"

func removeDuplicates(arra []int) ([]int, int) {
	if len(arra) == 0 {
		return arra, 0
	}
	// Slow pointer (i) - tracks the position of the last unique element
	i := 0

	// Fast pointer (j) - traverses the array
	for j := 1; j < len(arra); j++ {
		// If current element is different from the last unique one
		if arra[i] != arra[j] {
			// Place the new unique element right after the last unique one
			arra[i+1] = arra[j]
			// Move slow pointer forward
			i++
		}
	}

	// Return the array truncated to the new length
	newArray, newLength := getNewarray(arra, i+1)
	return newArray, newLength
}

func getNewarray(arra []int, length int) ([]int, int) {
	arra = arra[:length]
	return arra, length
}

func main() {
	var nums = []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	newArray, newLength := removeDuplicates(nums)

	fmt.Printf("After remove duplicates: %v, new length: %d\n", newArray, newLength)
}
