package main

import (
	"gin-demo/res"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/index", func(c *gin.Context) {
		// 返回一个标准的 JSON 响应
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{},
		})
		// 使用 res 包中的响应函数
		res.OkMsg(c, "Success")
	})
	// 登录路由，使用 res 包中的响应函数
	r.GET("/login", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Login successful",
			"data": gin.H{},
		})
		// 使用 res 包中的响应函数
		res.OkMsg(c, "Login Success")
	})
	// 注册路由，使用 res 包中的响应函数
	r.POST("/users", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Login successful",
			"data": gin.H{},
		})
		// 使用 res 包中的响应函数
		res.OkData(c, map[string]any{
			"user": "Alice",
			"age":  30})
	})
	// 错误路由，使用 res 包中的响应函数
	r.GET("/error", func(c *gin.Context) {
		// 返回一个错误的 JSON 响应
		c.JSON(400, gin.H{
			"code": 1001,
			"msg":  "参数错误",
			"data": gin.H{},
		})
		// 使用 res 包中的响应函数
		res.FailCode(c, 1001, "参数错误")
	})

	// 简单的 ping 路由，使用 res 包中的响应函数
	r.GET("/ping", func(c *gin.Context) {
		res.OkMsg(c, "pong")
	})

	// 启动服务器，监听默认端口 8080
	r.Run(":8080")

}
