package main

import (
	"github.com/gin-gonic/gin"
)

type BindForm struct {
	Name string `form:"name" binding:"required"`
	Age  int    `form:"age" binding:"required"`
}

func QueryExample() {
	r := gin.Default()
	r.GET("/query", func(c *gin.Context) {
		var form BindForm
		// 绑定查询参数到结构体
		if err := c.ShouldBindQuery(&form); err != nil {
			c.JSON(400, gin.H{
				"code": 1002,
				"msg":  "参数错误",
				"data": gin.H{},
			})
			return
		}
		// 返回绑定的结果
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"name": form.Name,
				"age":  form.Age,
			},
		})
	})
	r.Run(":8080")
}

type QueryParams struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

func DynamicQueryExample() {
	r := gin.Default()
	r.GET("/search/:query", func(c *gin.Context) {
		var params QueryParams
		// 绑定查询参数到结构体
		if err := c.ShouldBindQuery(&params); err != nil {
			c.JSON(400, gin.H{
				"code": 1002,
				"msg":  "参数错误",
				"data": gin.H{},
			})
			return
		}
		query := c.Param("query")
		// 返回绑定的结果
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"query": query,
				"page":  params.Page,
				"size":  params.Size,
			},
		})
	})
	r.Run(":8080")
}

func FormExample() {
	r := gin.Default()
	r.POST("/submit", func(c *gin.Context) {
		var form BindForm
		// 绑定表单数据到结构体
		if err := c.ShouldBind(&form); err != nil {
			c.JSON(400, gin.H{
				"code": 1002,
				"msg":  "参数错误",
				"data": gin.H{},
			})
			return
		}
		// 返回绑定的结果
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"name": form.Name,
				"age":  form.Age,
			},
		})
	})
	r.Run(":8080")
}

type JSONData struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"required"`
}

func JSONExample() {
	r := gin.Default()
	r.POST("/json", func(c *gin.Context) {
		var jsonData JSONData
		// 绑定 JSON 数据到结构体
		if err := c.ShouldBindJSON(&jsonData); err != nil {
			c.JSON(400, gin.H{
				"code": 1002,
				"msg":  "参数错误",
				"data": gin.H{},
			})
			return
		}
		// 返回绑定的结果
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{
				"name": jsonData.Name,
				"age":  jsonData.Age,
			},
		})
	})
	r.Run(":8080")
}

func HeaderExample() {
	r := gin.Default()
	r.GET("/header", func(c *gin.Context) {
		// 获取请求头
		token := c.GetHeader("Authorization")
		// 验证 token
		if token != "my-secret-token" {
			c.JSON(401, gin.H{
				"code": 1003,
				"msg":  "Unauthorized",
				"data": gin.H{},
			})
			return
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "Success",
			"data": gin.H{},
		})
	})
	r.Run(":8080")
}

func main() {
	QueryExample()
}
