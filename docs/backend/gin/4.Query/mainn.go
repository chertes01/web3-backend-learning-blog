package main

import (
	"github.com/gin-gonic/gin"
)

func QueryExample() {
	r := gin.Default()

	r.GET("/search", func(c *gin.Context) {
		// 获取查询参数
		query := c.Query("q")               // 等同于 c.Request.URL.Query().Get("q")
		page := c.DefaultQuery("page", "1") // 提供默认值

		c.JSON(200, gin.H{
			"query": query,
			"page":  page,
		})
	})

	// 启动服务器
	r.Run(":8080")
}

func DynamicQueryExample() {
	r := gin.Default()

	r.GET("/search/:query", func(c *gin.Context) {
		// 获取动态参数
		query := c.Param("query")

		c.JSON(200, gin.H{
			"query": query,
		})
	})

	// 启动服务器
	r.Run(":8080")
}

func main() {
	QueryExample()
	DynamicQueryExample()
}
