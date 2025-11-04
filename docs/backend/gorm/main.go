package main

import (
	"gorm/config"
	//"gorm/models"
)

type user struct {
	Name string
	Age  int
}

func main() {
	// 初始化数据库
	config.InitDB()
	//查询
	var userlist []user             // 定义一个 user 类型的切片，用于存储查询结果
	config.DB.Find(&userlist)       // 使用 DB.Find 方法查询所有 user 记录，并将结果存储到 userlist 中
	for _, user := range userlist { // 遍历 userlist 切片
		println("Name:", user.Name, "Age:", user.Age)
	}
}
