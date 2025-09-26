package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

/*
题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次
递增操作，最后输出计数器的值。
考察点 ： sync.Mutex 的使用、并发数据安全。
*/

var (
	wg   sync.WaitGroup
	mu   sync.Mutex
	coun int
)

func increment() {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		mu.Lock()   //Lock
		coun++      //Increment share counter
		mu.Unlock() //Unlock
	}
}

/*
题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。启动10个协程，每个协程对计数器进行1000次
递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/

var (
	coun0 int64
)

func incrementAtomic() {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		atomic.AddInt64(&coun0, 1) //Atomic increment,no lock needed
	}
}

func main() {
	fmt.Println("============Task1:Self increment with \"sync.Mutex.Lock\"============")
	incrementCount := 10

	wg.Add(incrementCount)
	for i := 0; i < incrementCount; i++ {
		go increment()
	}
	wg.Wait()

	fmt.Println(coun)

	fmt.Println("============Task2:Self increment with \"atomic.AddInt64\"============")

	wg.Add(incrementCount)
	for i := 0; i < incrementCount; i++ {
		go incrementAtomic()
	}
	wg.Wait()

	fmt.Println(coun0)

}
