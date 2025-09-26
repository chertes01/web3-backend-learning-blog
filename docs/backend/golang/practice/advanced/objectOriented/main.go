/*
题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，
实现 Shape 接口。在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点 ：接口的定义与实现、面向对象编程风格。
*/

package main

import (
	"fmt"
	"math"
	"sync"
)

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	Height, Wight float64
}

type Circle struct {
	Radius float64
}

func (cir *Circle) Area() float64 {
	return math.Pi * math.Pow(cir.Radius, 2)
}

func (rec *Rectangle) Area() float64 {
	return rec.Height * rec.Wight
}

func (cir *Circle) Perimeter() float64 {
	return math.Pi * cir.Radius * 2
}

func (rec *Rectangle) Perimeter() float64 {
	return (rec.Height + rec.Wight) * 2
}

/*
题目 ：使用组合的方式创建一个 Person 结构体，包含 Name 和 Age 字段，再创建一个 Employee 结构体，
组合 Person 结构体并添加 EmployeeID 字段。为 Employee 结构体实现一个 PrintInfo() 方法，输出员工的信息。
考察点 ：组合的使用、方法接收者。
*/

var (
	mu sync.Mutex
	wg sync.WaitGroup
)

type Person struct {
	Name string
	Age  int
}

type Employee struct {
	EmployeeID int
	Person
}

func addInfo(empl *Employee, name string, age int, ID int) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()

	empl.Name = name
	empl.Age = age
	empl.EmployeeID = ID
	empl.PrintInfo()

}

func (empl *Employee) PrintInfo() {
	fmt.Printf("Employee:%v,Age:%v,ID:%08d\n", empl.Name, empl.Age, empl.EmployeeID)
}

func main() {
	rectangle1 := Rectangle{Height: 10, Wight: 8}
	circle1 := Circle{Radius: 5}

	shapes := []Shape{&rectangle1, &circle1}
	for _, shape := range shapes {
		fmt.Printf("Shape:%T,Perimeter:%v,Area:%v\n", shape, shape.Perimeter(), shape.Area())
	}

	fmt.Println()

	var emp Employee
	wg.Add(2)
	go addInfo(&emp, "Alice", 55, 8)
	go addInfo(&emp, "bob", 25, 2)
	wg.Wait()
}
