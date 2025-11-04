// 定义数据模型（数据库表）
package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Gender string

// 定义常量（枚举值）
const (
	Male   Gender = "male"
	Female Gender = "female"
)

// User 用户模型
type UserModel struct {
	gorm.Model
	Age             int              `gorm:"default:0" json:"age"`                    // 年龄
	Username        string           `gorm:"size:50;unique;not null" json:"username"` // 用户名，唯一且不能为空
	Email           string           `gorm:"size:255;unique;not null" json:"email"`   // 邮箱，唯一且不能为空
	Phone           string           `gorm:"size:11;unique;not null" json:"phone"`    //
	UserDetailModel *UserDetailModel `gorm:"foreignKey:UserID"`                       // 外键关联 UserDetailModel 的 UserID 字段
}

// UserDetail 用户详情模型（一对一表）
type UserDetailModel struct {
	UserID    int        `gorm:"primaryKey"`                      // 用户ID，唯一且不能为空
	Address   string     `gorm:"size:20" json:"address"`          // 地址
	UserModel *UserModel `gorm:"foreignKey:UserID;references:ID"` // 外键关联 UserModel 的 ID 字段

}

// Employer 雇主模型（一对多表）
type Employer struct {
	gorm.Model
	Name    string      `gorm:"size:50;unique;not null" json:"name"` // 雇主名称，唯一且不能为空
	Workers *[]Employee `gorm:"foreignKey:EmployerID"`               // 一对多关联 Worker 模型
}

// Employee 员工模型（一对多表）
type Employee struct {
	gorm.Model
	Name       string    `gorm:"size:50;not null" json:"name"` // 员工名称，不能为空
	EmployerID uint      `json:"employer_id"`                  // 外键，关联 Employer 模型
	Employer   *Employer `gorm:"foreignKey:EmployerID"`        // 关联的雇主
}

// 学生选课多对多表
type Trainee struct {
	gorm.Model
	Name    string   `gorm:"size:50;not null" json:"name"`
	Courses []Course `gorm:"many2many:trainee_courses;"`
}

type Course struct {
	gorm.Model
	Name     string    `gorm:"size:20;not null" json:"name"`
	Trainees []Trainee `gorm:"many2many:trainee_courses;"`
}

// 钩子函数
func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	// 在创建用户之前执行的逻辑，例如数据验证或预处理
	return nil
}

func (u *UserModel) AfterDelate(tx *gorm.DB) (err error) {
	// 在删除用户之后执行的逻辑，例如日志记录或清理相关数据
	fmt.Println("用户删除成功")
	return nil
}

// SQL中已有表
// 结构体中使用自定义类型
type Student struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Gender Gender `json:"gender"` // 使用自定义的 Gender 类型
}

// 结构体表示学生和课程的关联记录
type Student_course struct {
	StudentID int // 对应 id_stu
	CourseID  int // 对应 id_cour
}
