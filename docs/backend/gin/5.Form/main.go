package main

import (
	"github.com/gin-gonic/gin"
)

func FormExample() {
	r := gin.Default()
	r.POST("/submit", func(c *gin.Context) {
		// 获取表单字段（相当于name ，_:= c.GetPostForm("name")）
		name := c.PostForm("name")
		// 获取可选字段
		email, ok := c.GetPostForm("email")
		if !ok {
			email = "No email provided"
		}
		message := c.DefaultPostForm("message", "No message provided") // 提供默认值

		c.JSON(200, gin.H{
			"name":    name,
			"email":   email,
			"message": message,
		})
	})

	// 启动服务器
	r.Run(":8080")
}

func FileExample() {
	r := gin.Default()
	r.POST("/upload", func(c *gin.Context) {
		// 获取上传的文件
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(400, gin.H{"error": "No file is received"})
			return
		}

		// 保存文件到指定路径
		err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
		if err != nil {
			c.JSON(500, gin.H{"error": "Unable to save the file"})
			return
		}

		c.JSON(200, gin.H{
			"message": "File uploaded successfully",
			"file":    file.Filename,
		})
	})
	// 启动服务器
	r.Run(":8080")
}

func MultipartFormExample() {
	r := gin.Default()
	r.POST("/upload-multiple", func(c *gin.Context) {
		// 获取多个上传的文件
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(400, gin.H{"error": "No files are received"})
			return
		}
		files := form.File["files"]

		var uploadedFiles []string
		for _, file := range files {
			// 保存每个文件到指定路径
			err = c.SaveUploadedFile(file, "./uploads/"+file.Filename)
			if err != nil {
				c.JSON(500, gin.H{"error": "Unable to save the file: " + file.Filename})
				return
			}
			uploadedFiles = append(uploadedFiles, file.Filename)
		}

		c.JSON(200, gin.H{
			"message": "Files uploaded successfully",
			"files":   uploadedFiles,
		})
	})
	// 启动服务器
	r.Run(":8080")
}

func main() {
	FormExample()
	FileExample()
}
