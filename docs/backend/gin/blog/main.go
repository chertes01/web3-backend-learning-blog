package main

import (
	routes "blog/Routes"
	"blog/config"
)

func main() {
	config.InitLog()
	// 创建数据库
	config.CreatDB()
	// 初始化数据库
	config.InitDB()
	//迁移数据库结构
	config.MigrateDB()
	// 设置路由
	r := routes.SetupRouter()
	// 运行服务器
	config.Log.Info("服务器启动在 :8080 端口")
	r.Run(":8080")
}
