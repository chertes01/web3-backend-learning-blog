/*
有效的括号

考察：字符串处理、栈的使用

题目：给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效
*/

package main

import (
	"fmt"
)

// bracket 映射表：右括号 -> 左括号
var bracket = map[rune]rune{
	'}': '{',
	')': '(',
	']': '[',
}

// isValidDebug 判断字符串是否是有效括号
func isValidDebug(s string) bool {
	stack := []rune{} // 模拟栈

	fmt.Println("输入字符串:", s)
	fmt.Println("开始逐步匹配括号：")

	for _, value := range s {
		switch value {
		// 遇到左括号：入栈
		case '{', '[', '(':
			stack = append(stack, value)

		// 遇到右括号：检查栈顶是否匹配
		case '}', ']', ')':
			if len(stack) > 0 && stack[len(stack)-1] == bracket[value] {

				stack = stack[:len(stack)-1] // 出栈

			} else {
				fmt.Printf("遇到右括号 %c -> 栈顶不匹配或栈空 ❌\n", value)
				return false
			}
		default:
			// 如果遇到非括号字符，这里直接忽略，也可以报错
			fmt.Printf("遇到非括号字符 %c -> 忽略\n", value)
		}
	}

	// 遍历完后，栈为空说明所有左括号都匹配
	if len(stack) == 0 {
		fmt.Println("最终: 栈为空括号匹配成功")
		return true
	} else {
		fmt.Printf("最终: 栈未清空剩余栈: %q\n", string(stack))
		return false
	}
}

func main() {
	str := "[(hello)]" // 测试字符串

	fmt.Println("======================================")
	result := isValidDebug(str)
	fmt.Println("======================================")
	fmt.Printf("最终结果: %v\n", result)
}
