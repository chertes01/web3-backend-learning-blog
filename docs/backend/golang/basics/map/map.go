package main

import (
	"fmt"
	"sync"
)

var (
	student = map[string]map[string]float64{
		"Tom": {
			"math":    90,
			"english": 88,
		},
		"Jack": {
			"math":    78,
			"english": 68,
		},
		"Sam": {
			"math":    85,
			"english": 80,
		},
	}
	mu sync.Mutex
	wg sync.WaitGroup
)

func addScore(name string, sub string, score float64, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	defer mu.Unlock()

	if student[name] == nil {
		student[name] = make(map[string]float64)
	}
	student[name][sub] = score

}

func addStudent(stu map[string]map[string]float64, name string) {
	if _, ok := stu[name]; !ok {
		stu[name] = map[string]float64{}
	}
}

func getScore(stu map[string]map[string]float64, name, sub string) float64 {
	if _, ok := stu[name]; ok {

		return stu[name][sub]

	} else {
		fmt.Println("no this student")
		return -1
	}
}

func deleteStudent(stu map[string]map[string]float64, name string) {
	mu.Lock()
	defer mu.Unlock()
	delete(stu, name)
}

func deletesubject(stu map[string]map[string]float64, name string, sub string) {
	mu.Lock()
	defer mu.Unlock()
	delete(stu[name], sub)
}

func printAll(stu map[string]map[string]float64) {
	mu.Lock()
	defer mu.Unlock()
	for name, subject := range stu {
		fmt.Println("Name:", name)
		for sub, score := range subject {
			fmt.Printf(" Subject:%s,Score:%f\n", sub, score)
		}
	}
}

func main() {
	addStudent(student, "Alice")
	addStudent(student, "Bob")

	wg.Add(4)

	go addScore("Alice", "math", 100, &wg)

	go addScore("Alice", "english", 95, &wg)

	go addScore("Bob", "math", 98, &wg)

	go addScore("Bob", "english", 97, &wg)

	wg.Wait()

	fmt.Printf("Name:Alice,subject:math,score:%f\n", getScore(student, "Alice", "math"))
	printAll(student)
}
