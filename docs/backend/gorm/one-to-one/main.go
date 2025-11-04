package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"

	"gorm.io/gorm"
)

func creat() {
	config.CreatDB()
	config.InitDB2()
}

func migrate() {
	err := config.DB.AutoMigrate(&models.UserModel{}) // 自动创建或更新 User 表结构
	if err != nil {
		panic(err)
	}
	fmt.Println("迁移成功")
	err = config.DB.AutoMigrate(&models.UserDetailModel{}) // 自动创建或更新 User 表结构
	if err != nil {
		panic(err)
	}
	fmt.Println("迁移成功")
}

func createmain() {
	users := []models.UserModel{
		{
			Age:      25,
			Username: "john_doe",
			Email:    "john_doe@qq.com",
			Phone:    "12345678901",
			UserDetailModel: &models.UserDetailModel{
				Address: "123 Main St",
			},
		},
		{
			Age:      30,
			Username: "alice_wang",
			Email:    "alice_wang@qq.com",
			Phone:    "13388889999",
			UserDetailModel: &models.UserDetailModel{
				Address: "456 Park Ave",
			},
		},
		{
			Age:      28,
			Username: "mike_chen",
			Email:    "mike_chen@qq.com",
			Phone:    "13911112222",
		},
		{
			Age:      22,
			Username: "sarah_lin",
			Email:    "sarah_lin@qq.com",
			Phone:    "13755556666",
			UserDetailModel: &models.UserDetailModel{
				Address: "321 Hill Rd",
			},
		},
	}

	for _, user := range users {
		if err := config.DB.Create(&user).Error; err != nil {
			fmt.Printf("❌ 创建用户 %s 失败：%v\n", user.Username, err)
			continue
		}
		fmt.Printf("✅ 创建用户 %s 成功！\n", user.Username)
	}
}

func insert() {
	err := config.DB.Create(&models.UserDetailModel{
		Address:   "789 Ocean Blvd",
		UserModel: &models.UserModel{Model: gorm.Model{ID: 3}},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func read() {
	//正向查询
	var user models.UserModel
	config.DB.Preload("UserDetailModel").Find(&user, "ID=?", 3)
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s,地址:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone, user.UserDetailModel.Address)
	//反向查询
	var userDetail models.UserDetailModel
	config.DB.Preload("UserModel").Find(&userDetail, "user_id=?", 4)
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s,地址:%s\n", userDetail.UserModel.ID, userDetail.UserModel.Username, userDetail.UserModel.Age, userDetail.UserModel.Email, userDetail.UserModel.Phone, userDetail.Address)
}

func delete() {
	//级联删除
	config.DB.Unscoped().Select("UserDetailModel").Delete(&models.UserModel{Model: gorm.Model{ID: 4}})
	//set null 删除
	config.DB.Model(&models.UserModel{Model: gorm.Model{ID: 3}}).Association("UserDetailModel").Clear()
	config.DB.Unscoped().Delete(&models.UserModel{Model: gorm.Model{ID: 3}})
}

func main() {
	creat()
	config.InitDB2()
	migrate()
	createmain()
	insert()
	read()
	delete()
}
