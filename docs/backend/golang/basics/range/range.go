package main

import (
	"fmt"
	"time"
)

func addData(ch chan int) {
	size := cap(ch)
	for i := 0; i < size; i++ {
		ch <- i
		time.Sleep(1 * time.Second)
	}
	close(ch)
}

func main() {
	ch := make(chan int, 10)

	go addData(ch)

	for i := range ch {
		fmt.Println(i)
	}
}
