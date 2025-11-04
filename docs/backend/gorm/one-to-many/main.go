package main

import (
	"fmt"
	"gorm/config"
	"gorm/models"
)

func migrate() {
	err := config.DB.AutoMigrate(&models.Employer{}, &models.Employee{}) // 自动创建或更新 User
	if err != nil {
		panic(err)
	}
	fmt.Println("迁移成功")

}

func createmain() {
	employers := []models.Employer{
		{
			Name: "TechCorp",
			Workers: &[]models.Employee{
				{Name: "Alice"},
				{Name: "Bob"},
			},
		},
		{
			Name: "BizInc",
			Workers: &[]models.Employee{
				{Name: "Charlie"},
				{Name: "David"},
			},
		},
	}

	for _, employer := range employers {
		result := config.DB.Create(&employer)
		if result.Error != nil {
			fmt.Printf("创建雇主 %s 失败: %v\n", employer.Name, result.Error)
		} else {
			fmt.Printf("创建雇主 %s 成功，包含员工:\n", employer.Name)
			for _, worker := range *employer.Workers {
				fmt.Printf(" - 员工姓名: %s\n", worker.Name)
			}
		}
	}
}

func read() {
	// 预加载雇主及其员工
	var employer models.Employer
	config.DB.Preload("Workers").Take(&employer, "id = ?", 1)
	// 输出雇主及其员工信息
	fmt.Printf("雇主ID:%d,雇主名称:%s,员工数:%d\n", employer.ID, employer.Name, len(*employer.Workers))
	for _, worker := range *employer.Workers {
		fmt.Printf(" - 员工姓名: %s\n", worker.Name)
	}

	//
	var worker models.Employee
	config.DB.Preload("Employer").Take(&worker, "name = ?", "David")
	fmt.Printf("员工ID:%d,员工姓名:%s,雇主ID:%d,雇主名称:%s\n", worker.ID, worker.Name, worker.EmployerID, worker.Employer.Name)
}

func insert() {

	newWorkers := []models.Employee{
		{Name: "Eve"},
		{Name: "Frank"},
	}

	// 插入新员工并关联到雇主ID为1的雇主
	for _, worker := range newWorkers {
		worker.EmployerID = 1
		result := config.DB.Create(&worker)
		if result.Error != nil {
			fmt.Printf("插入员工 %s 失败: %v\n", worker.Name, result.Error)
		} else {
			fmt.Printf("插入员工 %s 成功，关联雇主ID: %d\n", worker.Name, worker.EmployerID)
		}
	}
	// 插入新员工并关联到雇主ID为2的雇主
	newWorker := models.Employee{Name: "Hany", EmployerID: 2}
	result := config.DB.Create(&newWorker)
	if result.Error != nil {
		fmt.Printf("插入员工 %s 失败: %v\n", newWorker.Name, result.Error)
	} else {
		fmt.Printf("插入员工 %s 成功，关联雇主ID: %d\n", newWorker.Name, newWorker.EmployerID)
	}
}

func changeAssociation() {
	// 将员工 David 从雇主 BizInc 转移到 TechCorp(Append操作)
	var employer models.Employer
	config.DB.Preload("Workers").Take(&employer, "name = ?", "TechCorp")

	var newWorker models.Employee
	config.DB.Preload("Employer").Take(&newWorker, "name = ?", "David")
	// 输出变更前的信息
	fmt.Printf("更新前David雇主：%s\n", newWorker.Employer.Name)
	// 将 David 关联到 TechCorp
	err := config.DB.Model(&employer).Association("Workers").Append(&newWorker)
	if err != nil {
		fmt.Printf("关联员工失败: %v\n", err)
		return
	}
	// 重新加载员工及其雇主
	config.DB.Preload("Employer").Take(&newWorker, "name = ?", "David")
	fmt.Printf("变更后David的雇主:%s，雇主ID:%d\n", newWorker.Employer.Name, newWorker.EmployerID)

	//Replace操作
	var anotherWorker models.Employee
	config.DB.Preload("Employer").Take(&anotherWorker, "name = ?", "Charlie")
	fmt.Printf("更新前Charlie雇主：%s\n", anotherWorker.Employer.Name)
	//Replace会将其他关联的记录移除，只保留指定的记录
	err = config.DB.Model(&employer).Association("Workers").Replace(&anotherWorker)
	if err != nil {
		fmt.Printf("替换员工失败: %v\n", err)
		return
	}
	// 重新加载员工及其雇主
	config.DB.Preload("Employer").Take(&anotherWorker, "name = ?", "Charlie")
	fmt.Printf("变更后Charlie的雇主:%s，雇主ID:%d\n", anotherWorker.Employer.Name, anotherWorker.EmployerID)
}

func update() {
	// 更新雇主的名称
	result := config.DB.Model(&models.Employer{}).Where("name = ?", "TechCorp").Update("name", "NewTechCorp")
	if result.Error != nil {
		fmt.Printf("更新雇主失败: %v\n", result.Error)
		return
	}

	fmt.Println("雇主更新成功")
}

func delete() {
	// 删除雇主及其员工
	result := config.DB.Select("Workers").Unscoped().Delete(&models.Employer{}, 2)
	if result.Error != nil {
		fmt.Printf("删除雇主失败: %v\n", result.Error)
		return
	}

	fmt.Println("雇主及其员工删除成功")
}

func main() {
	config.InitDB2()
	migrate()
	createmain()
	read()
	insert()
	changeAssociation()
	update()
	delete()
}
