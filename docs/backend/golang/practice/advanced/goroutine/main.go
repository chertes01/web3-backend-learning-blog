/*
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

/*
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

var (
	wg   sync.WaitGroup
	odd  []int
	even []int
)

func printEven() {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		if i%2 == 0 {
			even = append(even, i)
		}
	}

	fmt.Println("In 0~10,even is:", even)
}

func printOdd() {
	defer wg.Done()
	for i := 1; i <= 10; i++ {
		if i%2 != 0 {
			odd = append(odd, i)
		}
	}
	fmt.Println("In 0~10,odd is:", odd)
}

/*
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/

func reciveTask(tas []func(int)) {
	for index, task := range tas {

		wg.Add(1)
		go func(i int, tas func(int)) {
			defer wg.Done()
			//Start timer
			start := time.Now()
			//Simulate I/O port transmission and reception
			tas(i)

			fmt.Printf("Excution time of task %d:%v\n", i, time.Since(start))

		}(index, task)
	}
	wg.Wait()
}

func main() {

	wg.Add(2)

	go printEven()
	go printOdd()

	wg.Wait()

	//"int" is pass in the number of excuting program
	tasksIO := []func(int){
		func(i int) {
			//Simulate waiting for 500ms~5000ms
			se := rand.Intn(4501) + 500
			time.Sleep(time.Duration(se) * time.Millisecond)
			fmt.Printf("Task %d is ending\n", i)
		},
		func(i int) {

			se := rand.Intn(4501) + 500
			time.Sleep(time.Duration(se) * time.Millisecond)
			fmt.Printf("Task %d is ending\n", i)
		},
		func(i int) {

			se := rand.Intn(4501) + 500
			time.Sleep(time.Duration(se) * time.Millisecond)
			fmt.Printf("Task %d is ending\n", i)
		},
	}

	reciveTask(tasksIO)

}
