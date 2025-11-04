package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
)

func query() {
	var userlist []models.UserModel
	//带条件查询（Find可查多条）
	config.DB.Find(&userlist, "Age >?", 30)
	for _, u := range userlist {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}

	//根据主键查询
	var user, user1, user2 models.UserModel
	//Take查一条
	config.DB.Debug().Take(&user, 1) // SELECT * FROM `user_models` WHERE `user_models`.`id` = 1 AND `user_models`.`deleted_at` IS NULL LIMIT 1
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)

	//First查一条
	config.DB.Debug().First(&user1) //SELECT * FROM `user_models` WHERE `user_models`.`id` = 3 AND `user_models`.`deleted_at` IS NULL ORDER BY `user_models`.`id` LIMIT 1
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)

	//Last查一条
	config.DB.Debug().Last(&user2) //SELECT * FROM `user_models` WHERE `user_models`.`deleted_at` IS NULL ORDER BY `user_models`.`id` DESC LIMIT 1
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user2.ID, user2.Username, user2.Age, user2.Email, user2.Phone)
}

func main() {
	config.InitDB()
	query()

}
