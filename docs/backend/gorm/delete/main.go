package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
)

func delete() {
	var user models.UserModel
	user.ID = 2
	// 使用 Delete 方法删除记录
	config.DB.Delete(&user)
	// 根据主键删除记录
	config.DB.Delete(&models.UserModel{}, 3)
	// 批量软删除记录
	config.DB.Delete(&models.UserModel{}, []int{3, 4})
	// 根据条件删除记录
	config.DB.Where("age > ?", 40).Delete(&models.UserModel{})

	// 查询软删除的记录
	var users []models.UserModel
	config.DB.Unscoped().Where("deleted_at IS NOT NULL").Find(&users)
	for _, u := range users {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}

	// 物理删除记录(硬删除)
	config.DB.Unscoped().Delete(&models.UserModel{}, []int{3, 4})

}

func main() {
	config.InitDB()
	delete()

}
