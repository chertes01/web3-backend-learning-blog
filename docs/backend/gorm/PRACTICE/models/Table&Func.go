package models

import (
	"gorm.io/gorm"
)

//使用Gorm定义 User 、Post 和 Comment 模型，其中 User 与 Post 是一对多关系（一个用户可以发布多篇文章）， Post 与 Comment 也是一对多关系（一篇文章可以有多个评论）

type User struct {
	gorm.Model
	Name      string `gorm:"size:20;not null" json:"name"`
	PostCount uint   `gorm:"default:0" json:"post_count"` // 用户发布的文章数量
	Posts     []Post `gorm:"foreignKey:UserID"`           // 一对多关联 Post 模型
}

type Post struct {
	gorm.Model
	Title         string `gorm:"size:100;not null" json:"title"`
	Content       string `gorm:"size:500;not null" json:"content"`
	UserID        uint   `json:"user_id"`           // 外键，关联 User 模型
	User          User   `gorm:"foreignKey:UserID"` // 关联的用户
	CommentStatus bool   `gorm:"default:false"`     // 用于表示文章评论状态
	Comments      []Comment
}

type Comment struct {
	gorm.Model
	Content string
	PostID  uint
	Post    Post `gorm:"foreignKey:PostID"` // 关联的文章
}

func (p *Post) AfterCreate(tx *gorm.DB) (err error) {
	// 创建文章后，更新用户的文章数量
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count + ?", 1)).Error
}

func (p *Post) AfterDelete(tx *gorm.DB) (err error) {
	// 删除文章后，更新用户的文章数量
	return tx.Model(&User{}).
		Where("id = ?", p.UserID).
		UpdateColumn("post_count", gorm.Expr("post_count - ?", 1)).Error
}

func (c *Comment) AfterCreate(tx *gorm.DB) (err error) {
	// 创建评论后，更新文章的评论状态为已评论
	return tx.Model(&Post{}).
		Where("id = ?", c.PostID).
		Update("comment_status", true).Error
}

func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	// 删除评论后，检查文章是否还有评论，如果没有则更新评论状态为未评论
	var count int64
	tx.Model(&Comment{}).Where("post_id = ?", c.PostID).Count(&count)
	if count == 0 {
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("comment_status", false).Error
	}
	return nil
}
