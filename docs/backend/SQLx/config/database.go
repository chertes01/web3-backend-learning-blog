package config

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	username = "testuser"
	password = "123456test"
	host     = "127.0.0.1"
	port     = 3306
	dbname   = "SQLxDB"
)
var DB *sqlx.DB

func CreatDB() error {
	//连接MySQL服务器
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ MySQL连接失败: %v", err) //打印错误信息 并立即终止程序（os.Exit(1))
		return err
	}
	log.Println("✅ MySQL连接成功！")
	//创建数据库
	createDBSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;", dbname)
	_, err = DB.Exec(createDBSQL)
	if err != nil {
		log.Fatalf("❌ 创建数据库失败: %v", err) //打印错误信息 并立即终止程序（os.Exit(1))
		return err
	}
	log.Println("✅ 数据库创建成功！")
	return nil
}

func InitDB() error {
	//连接数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	var err error
	DB, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ 连接数据库失败: %v", err) //打印错误信息 并立即终止程序（os.Exit(1))
		return err
	}

	log.Println("✅ 数据库连接成功！")
	return nil
}
