package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string    `gorm:"type:varchar(100);size:100;not null" json:"Title"`
	Content  string    `gorm:"type:text;size:100000;not null" json:"Content"`
	UserID   uint      `gorm:"not null" json:"user_id"`
	User     User      `gorm:"foreignKey:UserID;references:ID"`
	Comments []Comment `gorm:"foreignKey:PostID"`
}
