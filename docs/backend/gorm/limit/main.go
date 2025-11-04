package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"

	"gorm.io/gorm"
	grom "gorm.io/gorm"
)

func limit() {
	var user []models.UserModel
	limit := 2
	//分页查询
	for page := 1; page <= 3; page++ {
		config.DB.Limit(limit).Offset((page - 1) * limit).Find(&user)
		for _, u := range user {
			fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
		}
	}
}

// Age30 返回一个根据年龄大于30过滤的命名范围函数
func Age30(tx *grom.DB) *gorm.DB {
	return tx.Where("age > ?", 30)
}

// AgeIn 返回一个根据年龄列表过滤的命名范围函数
func AgeIn(AgeList []int) func(tx *grom.DB) *gorm.DB {
	return func(tx *grom.DB) *gorm.DB {
		return tx.Where("age IN ?", AgeList)
	}
}

func main() {
	config.InitDB()
	limit()

	// 使用 Scopes 方法应用命名范围查询
	fmt.Println("使用 Scopes 方法应用命名范围查询")
	var usersAge []models.UserModel
	config.DB.Scopes(Age30).Find(&usersAge)
	for _, u := range usersAge {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}

	// 使用 Scopes 方法应用多个命名范围查询
	fmt.Println("使用 Scopes 方法应用多个命名范围查询")
	var usersAgeIn []models.UserModel
	config.DB.Scopes(AgeIn([]int{18, 30, 20})).Find(&usersAgeIn)
	for _, u := range usersAgeIn {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}
}
