package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
)

func migrate() {
	err := config.DB.AutoMigrate(&models.UserModel{}) // 自动创建或更新 User 表结构
	if err != nil {
		panic(err)
	}
	fmt.Println("迁移成功")
}

func main() {
	// 初始化数据库
	config.InitDB()
	// 自动迁移模式
	migrate()

}
