package main

import (
	"github.com/gin-gonic/gin"
)

func directlyRequestAPI() {
	r := gin.Default()

	r.GET("/file", func(c *gin.Context) {
		c.Header("Content-Type", "application/octet-stream")            // 表示是文件流，唤起浏览器下载，一般设置了这个，就要设置文件名
		c.Header("Content-Disposition", "attachment; filename=文件下载.go") // 用来指定下载下来的文件名
		c.File("文件下载.go")                                               // 指定要下载的文件路径
	})
	// 启动服务器
	r.Run(":8080") //文件不存在 会报错 404
}

func frontendRequestsBackendAPI() {

}

func main() {
	directlyRequestAPI()
}
