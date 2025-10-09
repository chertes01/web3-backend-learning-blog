package main

import (
	"fmt"
	"time"
)

// 向 channel 中添加数据，并在发送完毕后关闭 channel
func addData(ch chan int) {
	size := cap(ch) // 获取 channel 的容量
	for i := 0; i < size; i++ {
		ch <- i                     // 发送数据到 channel
		time.Sleep(1 * time.Second) // 每次发送后休眠 1 秒
	}
	close(ch) // 发送完毕后关闭 channel
}

func main() {
	ch := make(chan int, 10) // 创建一个容量为 10 的缓冲 channel

	go addData(ch) // 启动 goroutine 向 channel 添加数据

	// 使用 range 从 channel 读取数据，直到 channel 被关闭
	for i := range ch {
		fmt.Println(i) // 打印接收到的数据
	}
}
