package main

import (
	"fmt"
	"strconv"
)

var (
	palindromes = "上海自来水来自海上"
	num         = 6123216
)

func convertToRunes(s interface{}) []rune {
	switch v := s.(type) {
	case string:
		return []rune(v)
	case int:
		return []rune(strconv.Itoa(v))
	case uint:
		return []rune(strconv.Itoa(int(v)))
	default:
		return []rune{}
	}
}

func isPalindrome(val interface{}) bool {

	runes := convertToRunes(val)

	fmt.Println(runes)

	reversed := make([]rune, len(runes))

	if len(runes) != 0 {

		for i := len(runes) - 1; i >= 0; i-- {
			reversed[len(runes)-1-i] = runes[i]
		}

		if string(reversed) == string(runes) {
			return true
		} else {
			return false
		}

	} else {
		return false
	}
}

func main() {
	fmt.Printf("%q 是否是回文: %v\n", palindromes, isPalindrome(palindromes))

	fmt.Printf("%d 是否是回文: %v\n", num, isPalindrome(num))
}
