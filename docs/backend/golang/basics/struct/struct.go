package main

import "fmt"

type Other struct{}

type Person struct {
	Name  string            `json:"name" gorm:"column:<name>"`
	Age   int               `json:"age" gorm:"column:<name>"`
	Call  func()            `json:"-" gorm:"column:<name>"`
	Map   map[string]string `json:"map" gorm:"column:<name>"`
	Ch    chan string       `json:"-" gorm:"column:<name>"`
	Arr   [32]uint8         `json:"arr" gorm:"column:<name>"`
	Slice []interface{}     `json:"slice" gorm:"column:<name>"`
	Ptr   *int              `json:"-"`
	O     Other             `json:"-"`
}

var noName1 = struct {
	field1 string
	field2 int
}{}

var noName2 = struct{}{}

var noNameStructCh = make(chan struct{}, 0)

type One struct {
	a string
}

type Two struct {
	One
	b string
}

type Three struct {
	One
	Two
	A string
	B string
	C string
}

func main() {

	var noName4 = struct {
		field1 string
		field2 int
	}{
		field1: "hello world",
		field2: 888,
	}

	first := One{a: "Apple"}
	second := Two{b: "Pinapple"}
	third := Three{A: "an", B: "bn", C: "cn"}

	fmt.Println(noName1, noName2, noName4, noNameStructCh)
	fmt.Print(first, second, third)
	fmt.Println(first.a, second.a, second.b, third.A, third.B, third.C)
	fmt.Println(third.One.a, third.Two.a, third.Two.b)
}
