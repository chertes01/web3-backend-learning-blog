package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
)

func where() {
	var user, user2, user3 models.UserModel
	//带条件查询
	fmt.Println("带条件查询")
	config.DB.Where("age = ?", 18).Take(&user)
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user.ID, user.Username, user.Age, user.Email, user.Phone)

	// 使用结构体作为查询条件(只查询非空字段)
	fmt.Println("使用结构体作为查询条件(只查询非空字段)")
	config.DB.Where(&models.UserModel{
		Age:      30,
		Username: "johndoe",
	}).Take(&user2)
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user2.ID, user2.Username, user2.Age, user2.Email, user2.Phone)

	// 使用 map 作为查询条件(可查询空字段)
	fmt.Println("使用 map 作为查询条件(可查询空字段)")
	config.DB.Where(map[string]interface{}{
		"age":      20,
		"username": "jack",
	}).Take(&user3)
	fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", user3.ID, user3.Username, user3.Age, user3.Email, user3.Phone)

	//嵌套where条件查询
	var stu []models.Student
	//var cour []models.Student_course
	fmt.Println("嵌套where条件查询")
	config.DB.Where("id IN (?)",
		config.DB.Model(&models.Student_course{}).
			Select("id_stu").
			Where("id_cour = ?", 2),
	).Find(&stu)

	for _, u := range stu {
		fmt.Printf("学生ID:%d,学生姓名:%s,学生性别:%s\n", u.ID, u.Name, u.Gender)
	}

	//排序查询
	fmt.Println("排序查询")
	var user4, user5 []models.UserModel
	config.DB.Order("id desc").Find(&user4)
	for _, u := range user4 {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}

	//或条件查询
	fmt.Println("或条件查询")
	config.DB.Or("age = ?", 18).Or("username = ?", "jack").Find(&user5)
	for _, u := range user5 {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}

	//Not条件查询
	fmt.Println("Not条件查询")
	var user6 []models.UserModel
	config.DB.Not("age > ?", 25).Not("username = ?", "johndoe").Find(&user6)
	for _, u := range user6 {
		fmt.Printf("用户ID:%d,用户名:%s,年龄:%d,邮箱:%s,手机号:%s\n", u.ID, u.Username, u.Age, u.Email, u.Phone)
	}
}

func main() {
	config.InitDB()
	where()

}
