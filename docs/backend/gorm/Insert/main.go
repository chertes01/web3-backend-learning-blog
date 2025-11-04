package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
	"math/rand"
	"time"
)

func main() {
	config.InitDB()
	// 创建新用户(单次插入)
	result := config.DB.Create(&models.UserModel{
		Age:      30,
		Username: "johndoe",
		Email:    "3626352@qq.com",
		Phone:    "12345678901",
	})
	if result.Error != nil {
		panic(result.Error)
	}
	println("插入成功，用户ID:", result.RowsAffected)

	//回填插入
	user := models.UserModel{
		Age:      20,
		Username: "jack",
		Email:    "12345665432@qq.com",
		Phone:    "12345678904",
	}
	result = config.DB.Create(&user)
	if result.Error != nil {
		panic(result.Error)
	}
	println("插入成功，用户ID:", user.ID)

	// 批量插入
	users := []models.UserModel{
		{
			Age:      rand.Intn(40) + 20,                                   // 年龄在20到59之间
			Username: fmt.Sprintf("random_user_%d", time.Now().UnixNano()), // 使用时间戳确保唯一性
			Email:    fmt.Sprintf("user%d@example.com", time.Now().UnixNano()),
			Phone:    fmt.Sprintf("1%010d", rand.Intn(1e10)), // 1开头的11位手机号
		},
		{
			Age:      rand.Intn(30) + 25,                                      // 年龄在25到54之间
			Username: fmt.Sprintf("another_user_%d", time.Now().UnixNano()+1), // 再次使用时间戳确保唯一性
			Email:    fmt.Sprintf("another%d@example.com", time.Now().UnixNano()+1),
			Phone:    fmt.Sprintf("1%010d", rand.Intn(1e10)),
		},
	}

	result = config.DB.Create(&users)
	if result.Error != nil {
		panic(result.Error)
	}
	println("批量插入成功，插入行数:", result.RowsAffected)
	for _, u := range users {
		println("批量插入用户ID:", u.ID)
	}
}
