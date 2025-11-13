package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;size:500;not null"  json:"Content"`
	UserID  uint   `gorm:"not null" json:"user_id"`
	User    User   `gorm:"foreignKey:UserID;references:ID"`
	PostID  uint   `gorm:"not null" json:"post_id"`
	Post    Post   `gorm:"foreignKey:PostID;references:ID"`
}
