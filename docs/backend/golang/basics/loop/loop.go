package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	loopGame1 "github.com/learn/loop/loop_game1"
	loopGame2 "github.com/learn/loop/loop_game2"
)

var (
	scores = []int{}
)

func getScore() {
	for i := 0; i < 10; i++ {
		scores = append(scores, rand.Intn(101))
	}
}

func averge() {
	total := 0
	for i := 0; i < len(scores); i++ {
		total += scores[i]
	}
	avg := float64(total) / float64(len(scores))
	fmt.Println("Average:", avg)
}

func highScore() {
	for i := 0; i < len(scores); i++ {
		max := scores[i]
		if max > 90 {
			fmt.Println("High Score: ", max)
			break
		}
		fmt.Println("Score:", scores[i])
	}
}

func scoreManagement() {
	m := make(map[string]int)
	for i, v := range scores {
		m[fmt.Sprintf("Student%d", i+1)] = v
	}

	keys := make([]string, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	fmt.Println(keys)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	getScore()
	averge()
	highScore()
	scoreManagement()

	loopGame1.ScoreGame()
	loopGame2.LoopGame2()
}
