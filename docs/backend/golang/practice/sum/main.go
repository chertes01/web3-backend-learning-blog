/*
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出和为目标值 target 的那两个整数，并返回它们的数组下标。
你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。
你可以按任意顺序返回答案。

示例 1：

输入：nums = [2,7,11,15], target = 9
输出：[0,1]
解释：因为 nums[0] + nums[1] == 9 ，返回 [0, 1] 。
示例 2：

输入：nums = [3,2,4], target = 6
输出：[1,2]
示例 3：

输入：nums = [3,3], target = 6
输出：[0,1]
*/

package mai

import (
	"fmt"
)

func match(arra []int, tar int) (bool, []int) {
	for index, value := range arra {
		for i, v := range arra[index+1:] {
			if value+v == tar {
				res := []int{arra[index], arra[index+1+i]}
				return true, res
			}
		}
	}
	return false, nil
}

func main() {
	nums := [...]int{2, 7, 11, 15}
	target := 9
	result, value := match(nums[:], target)
	fmt.Printf("There are two numbers in nums that sum up to target: %v,%v", result, value)
}
