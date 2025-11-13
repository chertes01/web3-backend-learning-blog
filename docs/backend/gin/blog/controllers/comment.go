package controllers

import (
	"blog/config"
	"blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreateComment(c *gin.Context) {
	var comment models.Comment
	// 获取用户ID
	userID := c.GetUint("user_id")
	// 绑定 JSON 数据到 comment 结构体(传入评论)
	if err := c.ShouldBindJSON(&comment); err != nil {
		// [日志] 记录参数绑定失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"user_id":    userID,
			"error":      err.Error(),
		}).Warn("创建评论失败：参数格式错误")
		// 返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	comment.UserID = userID
	// 保存评论到数据库
	if err := config.DB.Create(&comment).Error; err != nil {
		// [日志] 记录评论创建失败的信息
		config.Log.WithFields(logrus.Fields{
			"user_id": userID,
			"error":   err.Error(),
		}).Error("创建评论失败：数据库错误")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "评论创建成功", "comment": comment})
}

func GetCommentsByPost(c *gin.Context) {
	var comments []models.Comment
	// 从 URL 获取文章 ID
	postID := c.Param("post_id")
	// 查询该文章的所有评论
	if err := config.DB.Preload("User").Preload("Post").Where("post_id = ?", postID).Find(&comments).Error; err != nil {
		// [日志] 记录查询评论失败的信息
		config.Log.WithFields(logrus.Fields{
			"post_id": postID,
			"error":   err.Error(),
		}).Error("获取评论失败：数据库错误")
		// 返回错误响应
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func DeleteComment(c *gin.Context) {
	var comment models.Comment
	// 当前用户 ID（从 JWT 提取）
	userID := c.GetUint("user_id")
	// 从 URL 获取评论 ID
	commentID := c.Param("comment_id")
	// 查找评论
	if err := config.DB.First(&comment, commentID).Error; err != nil {
		// [日志] 记录评论未找到的信息
		config.Log.WithFields(logrus.Fields{
			"comment_id": commentID,
			"user_id":    userID,
			"error":      err.Error(),
		}).Warn("删除评论失败：评论未找到")
		// 返回错误响应
		c.JSON(http.StatusNotFound, gin.H{"error": "评论未找到"})
		return
	}
	// 检查当前用户是否是评论的作者
	if comment.UserID != userID {
		// [日志] 记录无权限删除评论的信息
		config.Log.WithFields(logrus.Fields{
			"comment_id": commentID,
			"user_id":    userID,
		}).Warn("删除评论失败：无权限删除此评论")
		// 返回错误响应
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除此评论"})
		return
	}
	// 删除评论
	if err := config.DB.Delete(&comment).Error; err != nil {
		// [日志] 记录删除评论失败的信息
		config.Log.WithFields(logrus.Fields{
			"comment_id": commentID,
			"user_id":    userID,
			"error":      err.Error(),
		}).Error("删除评论失败：数据库错误")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除评论失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "评论删除成功"})
}
