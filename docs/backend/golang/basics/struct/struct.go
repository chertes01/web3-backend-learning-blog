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
	a string // 普通字段
}

type Two struct {
	One        // 匿名嵌入 One，实现继承效果
	b   string // 新增字段
}

type Three struct {
	One // 多重嵌入
	Two
	A string
	B string
	C string
}

type User struct {
	name string
}

// ===== 值接收者方法 =====
func (u User) GetName() string {
	return u.name
}

func (u User) SetNameCopy(newName string) {
	u.name = newName // 修改副本，不影响原对象
}

// ===== 指针接收者方法 =====
func (u *User) GetNamePtr() string {
	return u.name
}

func (u *User) SetNamePtr(newName string) {
	u.name = newName // 修改原对象
}

// ===== 普通函数 =====
func UpdateUserByValue(u User, newName string) {
	u.name = newName // 修改副本，不影响原对象
}

func UpdateUserByPointer(u *User, newName string) {
	u.name = newName // 修改原对象
}

func main() {

	var noName4 = struct {
		field1 string
		field2 int
	}{
		field1: "hello world",
		field2: 888,
	}

	user := User{name: "Initial"}

	fmt.Println("初始值:", user.name)

	first := One{a: "Apple"}
	second := Two{b: "Pinapple"}
	third := Three{A: "an", B: "bn", C: "cn"}

	fmt.Println(noName1, noName2, noName4, noNameStructCh)
	fmt.Print(first, second, third)
	fmt.Println(first.a, second.a, second.b, third.A, third.B, third.C)
	fmt.Println(third.One.a, third.Two.a, third.Two.b)

	// 1. 普通函数传值 vs 传指针
	UpdateUserByValue(user, "ValueFunc") // 传值：不改变 user
	fmt.Println("UpdateUserByValue 后:", user.name)

	UpdateUserByPointer(&user, "PointerFunc") // 传指针：修改原对象
	fmt.Println("UpdateUserByPointer 后:", user.name)

	// 2. 值接收者方法 vs 指针接收者方法
	user.SetNameCopy("SetNameCopy 方法") // 值接收者：修改副本
	fmt.Println("SetNameCopy 方法后:", user.name)

	user.SetNamePtr("SetNamePtr 方法") // 指针接收者：修改原对象
	fmt.Println("SetNamePtr 方法后:", user.name)

	// 3. 不论 user 是值还是指针，都能调用两种方法（Go 自动转换）
	userPtr := &user
	fmt.Println("userPtr.GetName():", userPtr.GetName())       // 自动解引用
	fmt.Println("userPtr.GetNamePtr():", userPtr.GetNamePtr()) // 直接调用
}
