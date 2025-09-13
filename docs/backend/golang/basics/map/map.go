package main

import "fmt"

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
)

func addScore(stu map[string]map[string]float64, name string, sub string, score float64) {
	if stu[name] == nil {
		stu[name][sub] = score
	} else if stu[name][sub] == 0 {
		stu[name][sub] = score
	}
}

func addStudent(stu map[string]map[string]float64, name string) {
	if _, v := stu[name]; !v {
		stu[name] = map[string]float64{}
	}
}

func getScore(stu map[string]map[string]float64, name, sub string) float64 {
	if stu[name] != nil {
		if stu[name][sub] != 0 {
			return stu[name][sub]
		} else {
			fmt.Println("no this subject")
			return -1
		}
	} else {
		fmt.Println("no this student")
		return -1
	}
}

func deleteStudent(stu map[string]map[string]float64, name string) {
	delete(stu, name)
}

func deletesubject(stu map[string]map[string]float64, name string, sub string) {
	delete(stu[name], sub)
}

func printAll(stu map[string]map[string]float64) {
	for name, subject := range stu {
		fmt.Println("Name:", name)
		for sub, score := range subject {
			fmt.Printf(" Subject:%s,Score:%f\n", sub, score)
		}
	}
}

func main() {
	addStudent(student, "Alice")
	addScore(student, "Alice", "math", 100)
	addScore(student, "Alice", "english", 57)
	addStudent(student, "Bob")
	addScore(student, "Bob", "english", 99)
	addScore(student, "Bob", "math", 89)

	fmt.Printf("Name:Alice,subject:math,score:%f\n", getScore(student, "Alice", "math"))
	printAll(student)
}
