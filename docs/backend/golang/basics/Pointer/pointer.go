package main

import (
	"fmt"
	"unsafe"
)

var p *int
var q *string

func main() {
	i := 10
	p = &i

	q = new(string)
	*q = "hello"

	q2 := &q

	println(p, q, q2)
	println(*p == i, *q == "hello", *q2 == q)

	*p = 20
	*q = "world"
	println(i, *q)

	p2 := &p
	**p2 = 50
	println(*p2, p, **p2, i)

	a := "Hello, world!"
	upA := uintptr(unsafe.Pointer(&a))
	upA += 1

	c := (*uint8)(unsafe.Pointer(upA))
	fmt.Println(*c)

}
