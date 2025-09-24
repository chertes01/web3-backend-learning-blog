/*
编写一个函数来查找字符串数组中的最长公共前缀。

如果不存在公共前缀，返回空字符串 ""。



示例 1：

输入：strs = ["flower","flow","flight"]
输出："fl"
示例 2：

输入：strs = ["dog","racecar","car"]
输出：""
解释：输入不存在公共前缀。
*/

package main

import (
	"fmt"
)

func jugement(s []string) string {
	prefix := []rune{}

	for index, value := range s[0] {

		for i := 0; len(s) > i; i++ {

			if index >= len(s[i]) || value != rune(s[i][index]) {
				return string(prefix)
			}
		}
		prefix = append(prefix, value)
	}
	return string(prefix)
}

func publicPrefix(s []string) {

	pubPrefix := jugement(s)
	if len(pubPrefix) == 0 {
		fmt.Println("That string haven't public prefix")

	} else {
		fmt.Printf("In that string ,the longest public prefix is:%s\n", pubPrefix)
	}
}

func main() {

	strs := []string{"flow", "flower", "flowing"}
	strs2 := []string{"dog", "racecar", "car"}

	publicPrefix(strs)
	publicPrefix(strs2)
}
