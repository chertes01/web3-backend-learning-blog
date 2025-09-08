package main

import "fmt"

// 方式1
const a int = 1

// 方式2
const b = "test"

// 方式3
const c, d = 2, "hello"

// 方式4
const e, f bool = true, false

// 方式5
const (
	h    byte = 3
	i         = "value"
	j, k      = "v", 4
	l, m      = 5, false
)

const (
	n = 6
)

type Gender string

const (
	Male   Gender = "Male"
	Female Gender = "Female"
)

func (g *Gender) IsMale() bool {
	return *g == Male
}

func (g *Gender) IsFemale() bool {
	return *g == Female
}

type Month int

const (
	January Month = 1 + iota
	February
	March
	April
	May
	June
	July
	August
	September
	October
	November
	December
)

func main() {
	var gender = Male
	// fmt.Println(a, b, c, d, e, f)
	// fmt.Println(h, i, j, k, l, m)
	// fmt.Println(n)
	fmt.Println(gender.IsMale())
	fmt.Println(September, October, November, December)
}
