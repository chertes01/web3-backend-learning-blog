package config

import (
	modles "blog/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

var (
	username = "testuser"
	password = "123456test"
	host     = "127.0.0.1"
	port     = 3306
	dbname   = "ginTestDB"
)

func InitDB() {
	//连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束(禁用实体外键)
	})

	if err != nil {
		log.Fatalf("❌ 连接数据库失败: %v", err) //打印错误信息 并立即终止程序（os.Exit(1))
	}

	log.Println("✅ 数据库连接成功！")
	DB = db

}

func CreatDB() {
	//连接MySQL服务器
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true, // 禁用自动创建外键约束(禁用实体外键)
	})

	if err != nil {
		log.Fatalf("❌ MySQL连接失败: %v", err) //打印错误信息 并立即终止程序（os.Exit(1))
	}
	log.Println("✅ MySQL连接成功！")
	DB = db
	//创建数据库
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", dbname)
	err = DB.Exec(createDBSQL).Error
	if err != nil {
		log.Fatalf("❌ 创建数据库失败: %v", err) //打印错误信息 并立即终止程序（os.Exit(1))
	}
	log.Println("✅ 数据库创建成功！")
}

func MigrateDB() {
	err := DB.AutoMigrate(
		&modles.User{},
		&modles.Post{},
		&modles.Comment{},
	)
	if err != nil {
		log.Fatalf("❌ 数据库迁移失败: %v", err) //打印错误信息 并立即终止程序（os.Exit(1))
	}
	log.Println("✅ 数据库迁移成功！")
}
