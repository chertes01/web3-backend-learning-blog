package config

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = logrus.New()

func InitLog() {
	// 1. 配置 lumberjack 日志切割归档功能
	fileWriter := &lumberjack.Logger{
		Filename:   "logs/app.log", // 日志文件路径
		MaxSize:    10,             // 单个日志文件最大尺寸（单位：MB），超过就会切割
		MaxBackups: 5,              // 最多保留旧日志文件的个数
		MaxAge:     30,             // 保留最近多少天的日志
		Compress:   true,           // 是否压缩旧日志（变成 .gz 格式，节省磁盘空间）
	}

	// 2. 组合输出：同时输出到“控制台”和“日志文件”
	// os.Stdout: 让你在开发时看控制台
	// fileWriter: 让你的日志自动保存并切割
	writers := io.MultiWriter(os.Stdout, fileWriter)

	Log.SetOutput(writers)

	// 3. 设置日志格式
	Log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 让人类读得懂的时间格式
	})

	// 4. 设置日志级别
	Log.SetLevel(logrus.InfoLevel)

	// [可选] 如果你想知道是哪行代码打印的日志，可以开启这个，但会消耗一点点性能
	// Log.SetReportCaller(true)
}
