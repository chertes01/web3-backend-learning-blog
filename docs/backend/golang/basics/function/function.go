package main

import "fmt"

type Counter struct {
	value int
}

func (a *Counter) Add(n int) int {
	a.value += n
	return a.value
}

var square func(int) int = func(x int) int {
	x *= x
	return x
}

func makeMultiplier(factor int) func(y int) int {

	return func(y int) int {
		return factor * y
	}
}

func main() {
	Counter := Counter{0}
	fn := Counter.Add

	time10 := makeMultiplier(10)

	fmt.Println(fn(3))
	fmt.Println(fn(5))
	fmt.Println(square(5))
	fmt.Println(time10(5))

}
