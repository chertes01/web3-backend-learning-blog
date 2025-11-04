package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
)

func originalSQL() {
	var user []models.UserModel
	// 使用原生 SQL 查询
	config.DB.Raw("SELECT * FROM user_models WHERE age > ?", 25).Scan(&user)
	for _, u := range user {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}
	//Esec 原生 SQL 更新
	fmt.Println("Esec 原生 SQL 更新")
	var user2 []models.UserModel
	config.DB.Exec("UPDATE user_models SET age = age + 1 WHERE age < ?", 30)
	config.DB.Raw("SELECT * FROM user_models WHERE age < ?", 30).Scan(&user2)
	for _, u := range user2 {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}
}

func main() {
	config.InitDB()
	originalSQL()
}
