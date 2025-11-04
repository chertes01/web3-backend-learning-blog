package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"

	"gorm.io/gorm"
)

func migrate() {
	err := config.DB.AutoMigrate(&models.Trainee{}, &models.Course{}) // 自动创建或更新 User
	if err != nil {
		panic(err)
	}
	fmt.Println("迁移成功")

}

func createmain() {
	trainees := []models.Trainee{
		{Name: "Alice"},
		{Name: "Bob"},
		{Name: "Charlie"},
	}

	courses := []models.Course{
		{Name: "MySQL"},
		{Name: "Golang"},
		{Name: "Docker"},
		{Name: "Python"},
	}
	// 创建学员和课程
	for i := range trainees {
		config.DB.Create(&trainees[i])
	}
	for i := range courses {
		config.DB.Create(&courses[i])
	}

	// 建立多对多关联
	config.DB.Model(&trainees[0]).Association("Courses").Append(&courses[0], &courses[1])
	config.DB.Model(&trainees[1]).Association("Courses").Append(&courses[1], &courses[3])
	config.DB.Model(&trainees[2]).Association("Courses").Append(&courses[2], &courses[3])

	fmt.Println("创建学员和课程及其关联成功")

}

func insert() {
	//插入数据
	config.DB.Create(&models.Trainee{
		Name: "David",
		Courses: []models.Course{
			{Name: "Java"},
			{Name: "C++"},
		},
	})

}

func read() {
	//预加载查询
	var trainee models.Trainee
	config.DB.Preload("Courses").Take("name=?", "Bob")
	fmt.Printf("学员ID:%d,学员姓名:%s\n", trainee.ID, trainee.Name)
	//输出学员所选课程
	for _, course := range trainee.Courses {
		fmt.Printf("课程ID:%d,课程名称:%s\n", course.ID, course.Name)
	}
	fmt.Println("-----")
	//反向预加载查询
	var course models.Course
	config.DB.Preload("Trainees").Take("name=?", "Golang")
	fmt.Printf("课程ID:%d,课程名称:%s\n", course.ID, course.Name)
	//输出选修该课程的学员
	for _, trainee := range course.Trainees {
		fmt.Printf("学员ID:%d,学员姓名:%s\n", trainee.ID, trainee.Name)
	}
	fmt.Println("-----")
	//多对多关联查询
	var trainee2 models.Trainee
	config.DB.Model(&trainee2).Where("name=?", "Alice").Preload("Courses").First(&trainee2)
	fmt.Printf("学员ID:%d,学员姓名:%s\n", trainee2.ID, trainee2.Name)
	//输出学员所选课程
	for _, course := range trainee2.Courses {
		fmt.Printf("课程ID:%d,课程名称:%s\n", course.ID, course.Name)
	}
	fmt.Println("-----")
	//查询所有课程及其选修学员
	var courses []models.Course
	config.DB.Preload("Trainees").Find(&courses)
	for _, course := range courses {
		fmt.Printf("课程ID:%d,课程名称:%s\n", course.ID, course.Name)
		//输出选修该课程的学员
		for _, trainee := range course.Trainees {
			fmt.Printf("学员ID:%d,学员姓名:%s\n", trainee.ID, trainee.Name)
		}
	}
}

func changeAssociation() {
	var trainee models.Trainee
	config.DB.Preload("Courses").Take(&trainee, "name=?", "Alice")
	for _, course := range trainee.Courses {
		fmt.Printf("课程ID:%d,课程名称:%s\n", course.ID, course.Name)
	}
	// 添加关联
	config.DB.Model(&trainee).Association("Courses").Append([]models.Course{
		{Model: gorm.Model{ID: 4}},
		{Model: gorm.Model{ID: 3}},
	})
	// 替换原有关联
	config.DB.Model(&trainee).Association("Courses").Replace([]models.Course{
		{Model: gorm.Model{ID: 2}},
		{Model: gorm.Model{ID: 5}},
	})

}

func main() {
	config.InitDB2()
	migrate()
	createmain()
	insert()
	read()
	changeAssociation()
}
