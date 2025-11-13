package controllers

import (
	"blog/config"
	"blog/middle"
	"blog/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	// 定义输入结构体
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,max=20,min=8"`
	}
	// 绑定 JSON 到结构体
	if err := c.ShouldBindJSON(&input); err != nil {
		// [日志] 记录参数绑定失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":    c.ClientIP(),
			"error": err.Error(),
		}).Warn("注册失败：参数格式错误")
		// 返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		// [日志] 记录密码加密失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":    c.ClientIP(),
			"error": err.Error(),
		}).Error("注册失败：密码加密失败")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}
	// 创建用户实例
	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}
	// 保存用户到数据库
	if err := config.DB.Create(&user).Error; err != nil {
		// [日志] 记录用户创建失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":    c.ClientIP(),
			"email": input.Email,
			"error": err.Error(),
		}).Error("注册失败：用户创建失败")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户创建失败"})
		return
	}
	// 返回成功响应
	User_ID := user.ID
	c.JSON(http.StatusOK, gin.H{
		"user_id": User_ID,
		"message": "用户注册成功",
	})
}

func Login(c *gin.Context) {
	// 定义输入结构体
	var input struct {
		ID       uint   `json:"id"  binding:"omitempty"`
		Name     string `json:"name"`
		Email    string `json:"email" binding:"omitempty,email"`
		Password string `json:"password" binding:"required,max=20,min=8"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		// [日志] 记录参数绑定失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":    c.ClientIP(),
			"error": err.Error(),
		}).Warn("登录失败：参数格式错误")
		// 返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	//基础查询
	query := config.DB.Model(&user)

	//根据用户ID查询
	if input.ID != 0 {
		query = query.Where("id = ?", input.ID)
		//根据邮箱查询
	} else if input.Email != "" {
		query = query.Where("email = ?", input.Email)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请输入用户名或邮箱"})
		return
	}

	//执行查询
	if err := query.First(&user).Error; err != nil {
		// [日志] 记录用户不存在的信息
		config.Log.WithFields(logrus.Fields{
			"ip":    c.ClientIP(),
			"error": err.Error(),
		}).Warn("登录失败：用户不存在")
		// 返回错误响应
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户不存在"})
		return
	}

	// 比较密码
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		// [日志] 记录密码错误的信息
		config.Log.WithFields(logrus.Fields{
			"ip": c.ClientIP(),
		}).Warn("登录失败：密码错误")
		// 返回错误响应
		c.JSON(http.StatusUnauthorized, gin.H{"error": "密码错误"})
		return
	}
	// 生成JWT令牌
	token, err := middle.GenerateToken(user.ID, c)
	if err != nil {
		// [日志] 记录令牌生成失败的信息
		config.Log.WithFields(logrus.Fields{
			"ip":      c.ClientIP(),
			"user_id": user.ID,
			"error":   err.Error(),
		}).Error("登录失败：令牌生成失败")
		// 返回错误响应
		c.JSON(http.StatusInternalServerError, gin.H{"error": "令牌生成失败"})
		return
	}
	// 返回成功响应和令牌
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
	})
}
