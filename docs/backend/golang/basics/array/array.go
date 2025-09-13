package main

import (
	"fmt"
)

var (
	pArr = [3]*Person{
		&Person{"Alice", 30},
		&Person{"Bob", 25},
		&Person{"Charlie", 35},
	}
	arr2D = [3][3]int{
		{0, 10, 0},
		{0, 0, 90},
	}

	a = [...]int{1, 2, 3, 4, 5}
)

// Array of structs
type Person struct {
	Name string
	Age  int
}

func printPersons(param [3]*Person) {
	for v := range param {
		fmt.Printf("Name:%s , Age:%d\n", param[v].Name, param[v].Age)
	}
}

func modifyArray(param [5]int) {
	param[0] = 100
}
func modifyArrayPtr(param *[5]int) {
	param[0] = 100
}

func main() {

	fmt.Println(arr2D)
	fmt.Println(pArr)
	printPersons(pArr)
	modifyArray(a)
	fmt.Println(a) // a[0] is still 1
	modifyArrayPtr(&a)
	fmt.Println(a) // a[0] is now 100

}
