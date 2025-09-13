package main

import "fmt"

func analyzeInput(x interface{}) {
	switch str := x.(type) {
	case string:
		switch str {
		case "hello":
			fmt.Println("got hello")
		case "world":
			fmt.Println("got greeting")
		case "hi":
			fmt.Println("got greeting")
		default:
			fmt.Println("Please provide a greeting")
		}

		length := len(str)
		switch {
		case length > 5:
			fmt.Println("long string")
		case length%2 == 1:
			fmt.Println("odd number")
		case length%2 == 0:
			fmt.Println("even number")
		default:
		}
		fmt.Println("string length is", length)

		fmt.Println("string", str)

	case int:
		switch {
		case str == 0:
			fmt.Println("zero")
		case str == 1:
			fmt.Println("one")
		case str > 10:
			fmt.Println("n is big number")
		default:
			fmt.Println("n is a number")
		}
		fmt.Println("int", str)

	case float64:
		fmt.Println("float", str)
	default:
		fmt.Println("unknown")

	}
}

func checkNumber(n int) {

	flag := false
	m := -1

	if n > 100 {
		flag = true
		fmt.Println("n is a big number")
	} else if m = n % 2; m == 0 {
		fmt.Println("n is even")
	} else {
		fmt.Println("n is odd")
		if flag {
			fmt.Println("but flag is true")
		}
	}
	fmt.Println("flag is", flag)
}
func main() {
	analyzeInput("hello")
	analyzeInput("world")
	analyzeInput("hi")
	analyzeInput("test")

	analyzeInput(0)
	analyzeInput(1)
	analyzeInput(11)
	analyzeInput(5)

	analyzeInput("hello world")
	analyzeInput("hello!")
	analyzeInput("hello!!")

	analyzeInput("test")
	analyzeInput(1)
	analyzeInput(3.14)
	analyzeInput(true)
	analyzeInput([]int{1, 2, 3})

	checkNumber(50)
	checkNumber(101)
	checkNumber(7)
}
