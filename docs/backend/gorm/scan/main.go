package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
)

func scan() {
	var user []models.UserModel
	// Scan 方法将查询结果扫描到指定的结构体或切片中
	config.DB.Model(&models.UserModel{}).
		Select("id", "username").
		Where("Age>?", 30).
		Scan(&user)

	for _, u := range user {
		fmt.Printf("ID:%d, 用户名:%s\n", u.ID, u.Username)
	}
}

func Pluck() {
	//Pluck 方法用于查询单个字段并将结果存储在切片中
	//定义一个新的结构体来存储 Pluck 的结果
	type UserScan struct {
		ID       uint   // 主键ID
		Username string // 用户名，唯一且不能为空
	}
	var userscan []UserScan

	config.DB.Debug().Model(&models.UserModel{}).
		Where("Age>?", 30).
		Pluck("id,username", &userscan)

	for _, u := range userscan {
		fmt.Printf("ID:%d, 用户名:%s\n", u.ID, u.Username)
	}
}

func cunt() {
	type course struct {
		Cour  int `gorm:"column:id_cour"` // 对应 id_cour
		Count int
	}
	var cour []course
	config.DB.Debug().Model(&models.Student_course{}).
		Group("id_cour").
		Select("id_cour", "count(id_stu) as Count").
		Scan(&cour)
	for _, u := range cour {
		fmt.Printf("ID:%d, 选课人数:%d\n", u.Cour, u.Count)
	}
}

func distinct() {
	var student []int
	config.DB.Debug().Model(&models.Student_course{}).
		Distinct("id_stu").
		Pluck("id_stu", &student)

	for _, u := range student {
		fmt.Printf("ID:%d\n", u)
	}

}

func main() {
	config.InitDB()
	scan()
	Pluck()
	cunt()
	distinct()
}
