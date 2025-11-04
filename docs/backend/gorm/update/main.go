package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"

	"gorm.io/gorm"
)

func updateSave() {
	var user models.UserModel
	// 修改字段值
	user.ID = 2
	user.Username = "new_username"
	user.Age = 28
	user.Email = "14363467383@expl.com"
	user.Phone = "14363467383"
	// 使用 Save 方法更新记录(有主键记录则更新，无则插入，且可更新0值字段)
	config.DB.Save(&user)
	// 查询更新后的记录以验证更改
	config.DB.Find(&user, user.ID)
	fmt.Printf("更新前用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)
}

func updateUpdate() {
	var user models.UserModel
	// 修改字段值
	user.ID = 2
	user.Username = "Jason"
	// 使用 Update 方法更新记录（单字段更新）|UpdateColumns不走钩子
	config.DB.Model(&user).Update("Username", user.Username)
	// 查询更新后的记录以验证更改
	config.DB.Find(&user, user.ID)
	fmt.Printf("更新后用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)

}

func updateUpdates() {
	var user models.UserModel
	// 修改字段值
	user.ID = 2
	user.Username = "Jason"
	user.Age = 32
	user.Email = "67383143634@expl.com"
	user.Phone = "16738436343"
	// 使用 Updates 方法更新记录（多字段更新,不更新0值）
	config.DB.Model(&user).Updates(models.UserModel{Username: user.Username, Age: user.Age, Email: user.Email, Phone: user.Phone})
	// 查询更新后的记录以验证更改
	config.DB.Find(&user, user.ID)
	fmt.Printf("更新后用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)

	// 使用 Updates 方法更新记录（多字段更新,更新0值）
	config.DB.Model(&user).Updates(map[string]interface{}{"Age": 0})
	// 查询更新后的记录以验证更改
	config.DB.Find(&user, user.ID)
	fmt.Printf("更新后用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)
}

func updateExpr() {
	var user models.UserModel
	// 修改字段值
	user.ID = 2
	user.Username = "Jason"
	user.Age = 32
	// 使用 Updates 方法更新记录（多字段更新,更新0值）
	config.DB.Model(&user).Updates(map[string]interface{}{"Age": gorm.Expr("Age + ?", 18)})
	// 查询更新后的记录以验证更改
	config.DB.Find(&user, user.ID)
	fmt.Printf("更新后用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)
}

func main() {
	config.InitDB()
	//updateSave()
	//updateUpdate()
	//updateUpdates()
	updateExpr()

}
