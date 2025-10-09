package main

import (
	"fmt"
	"math/rand"
)

/*
题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到
通道中，另一个协程从通道中接收这些整数并打印出来。
考察点 ：通道的基本使用、协程间通信。
*/

var (
	ch = make(chan int)
)

func getNums() {

	for i := 1; i <= 10; i++ {
		//Send value to channel
		ch <- i
	}
	//Close the channel after sending
	close(ch)
}

/*
题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/

var (
	cacheCh = make(chan int, 100)
)

func getRandomNum() {
	for i := 0; i < 100; i++ {
		//Creat a random value and send the value to channel
		cacheCh <- rand.Intn(1000)
	}
	//Close the channel after sending
	close(cacheCh)
}

func main() {
	fmt.Println("=== Task 1: Unbuffered channel, print 1~10 ===")
	var nums []int

	//start producter goroutine
	go getNums()
	//Consumer in main goroutine prints values as they come
	for v := range ch {
		nums = append(nums, v)
	}
	fmt.Println(nums)

	fmt.Println("\n=== Task 2: Buffered channel, print 100 random numbers ===")
	go getRandomNum()

	counter := 0
	for v := range cacheCh {
		fmt.Printf("%4d ", v)
		counter++
		if counter%10 == 0 {
			fmt.Println()
		}
	}
}
