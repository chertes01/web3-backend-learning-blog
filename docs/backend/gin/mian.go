package main

import (
	"github.com/gin-gonic/gin"
)

type Response struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data"`
}

func main() {
	// 初始化：创建一个默认的 Gin 路由器
	gin.SetMode("release") // 设置为发布模式
	r := gin.Default()

	// 挂载一个简单的 GET 路由
	r.GET("/index", func(c *gin.Context) {
		c.JSON(200, Response{
			Message: "success",
			Code:    200,
			Data:    "Hello, Gin!",
		})
	})

	// 启动服务器，监听默认端口 8080
	r.Run(":8080") // 监听并在
}
