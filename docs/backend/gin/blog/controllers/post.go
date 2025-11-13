package controllers

import (
	"blog/config"
	"blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func CreatePost(c *gin.Context) {
	var post models.Post
	// 获取用户ID
	userID := c.GetUint("user_id")
	// 绑定 JSON 数据到 post 结构体(传入文章)
	if err := c.ShouldBindJSON(&post); err != nil {
		// [日志] 记录参数绑定失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
			"error":   err.Error(),
		}).Warn("创建文章失败：参数格式错误")
		// 返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.UserID = userID
	// 保存文章到数据库
	if err := config.DB.Create(&post).Error; err != nil {
		// [日志] 记录文章创建失败的信息
		config.Log.WithFields(logrus.Fields{
			"post_id": post.ID,
			"user_id": userID,
			"error":   err.Error(),
		}).Error("创建文章失败：数据库错误")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章创建成功", "post": post})
}

func GetAllPosts(c *gin.Context) {
	var posts []models.Post
	if err := config.DB.Preload("User").Preload("Comments").Find(&posts).Error; err != nil {
		// [日志] 记录获取文章列表失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": c.GetUint("user_id"),
			"error":   err.Error(),
		}).Error("获取文章列表失败：数据库错误")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取文章列表失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func GetPostByID(c *gin.Context) {
	var post models.Post
	postID := c.Param("post_id")
	if err := config.DB.Preload("User").Preload("Comments").First(&post, postID).Error; err != nil {
		// [日志] 记录获取文章失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": c.GetUint("user_id"),
			"error":   err.Error(),
		}).Error("获取文章失败：文章未找到")
		// 返回错误响应
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"post": post})
}

func GetPostsByUser(c *gin.Context) {
	var posts []models.Post
	userID := c.Param("user_id")
	if err := config.DB.Preload("User").Preload("Comments").Where("user_id = ?", userID).Find(&posts).Error; err != nil {
		// [日志] 记录获取文章失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": c.GetUint("user_id"),
			"error":   err.Error(),
		}).Error("获取文章失败：文章未找到")
		// 返回错误响应
		c.JSON(http.StatusNotFound, gin.H{"error": "用户未找到"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func UpdatePost(c *gin.Context) {
	var post models.Post
	// 当前用户 ID（从 JWT 提取）
	userID := c.GetUint("user_id")
	// 从 URL 获取文章 ID
	postID := c.Param("post_id")
	if err := config.DB.First(&post, postID).Error; err != nil {
		// [日志] 记录获取文章失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
			"error":   err.Error(),
		}).Error("获取文章失败：文章未找到")
		// 返回错误响应
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	// 判断当前用户是否为文章的作者
	if post.UserID != userID {
		// [日志] 记录获取文章失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
		}).Error("获取文章失败：没有权限更新此文章")
		// 返回错误响应
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限更新此文章"})
		return
	}
	// 绑定更新数据
	var input struct {
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	// 绑定 JSON 数据到输入结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		// [日志] 记录参数绑定失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
			"error":   err.Error(),
		}).Warn("更新文章失败：参数格式错误")
		// 返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 更新文章字段
	post.Title = input.Title
	post.Content = input.Content
	if err := config.DB.Save(&post).Error; err != nil {
		// [日志] 记录参数绑定失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
			"post_id": postID,
			"error":   err.Error(),
		}).Warn("更新文章失败：数据库错误")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章更新成功", "post": post})
}

func DeletePost(c *gin.Context) {
	var post models.Post
	// 从 URL 获取文章 ID
	userID := c.GetUint("user_id")
	// 当前用户 ID（从 JWT 提取）
	postID := c.Param("post_id")
	if err := config.DB.First(&post, postID).Error; err != nil {
		// [日志] 记录获取文章失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
			"post_id": postID,
			"error":   err.Error(),
		}).Error("获取文章失败：文章未找到")
		// 返回错误响应
		c.JSON(http.StatusNotFound, gin.H{"error": "文章未找到"})
		return
	}
	// 判断当前用户是否为文章的作者
	if post.UserID != userID {
		// [日志] 记录获取文章失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
			"post_id": postID,
		}).Error("获取文章失败：没有权限删除此文章")
		// 返回错误响应
		c.JSON(http.StatusForbidden, gin.H{"error": "没有权限删除此文章"})
		return
	}
	// 删除文章
	if err := config.DB.Delete(&post).Error; err != nil {
		// [日志] 记录参数绑定失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": userID,
			"post_id": postID,
			"error":   err.Error(),
		}).Warn("删除文章失败：数据库错误")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除文章失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "文章删除成功"})
}
