package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	// 加载模板文件
	r.LoadHTMLGlob("../templates/*")
	// 定义一个路由，渲染 HTML 模板
	r.GET("/index", func(c *gin.Context) {
		// 渲染模板并传递数据
		c.HTML(200, "index.html", map[string]any{
			"title":   "Gin HTML 渲染示例",
			"message": "Hello, Gin with HTML!",
		})
	})
	r.Run(":8080")
}
