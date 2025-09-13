package loopGame1

import (
	"fmt"
	"math/rand" // 修复4: 必须导入math/rand包
	"time"      // 修复5: 必须导入time包
)

func ScoreGame() {
outter:
	for round := 1; round <= 5; round++ {
		fmt.Printf("Round %d:\n", round)
		if round == 4 {
			break
		}
		for player := 1; player <= 3; player++ {

			Scores := rand.Intn(101)

			fmt.Printf("scored %d points\n", Scores)

			switch {
			case Scores < 60:
				fmt.Println("Failed")
			case Scores >= 60 && Scores < 80:
				fmt.Println("Passed")
			case Scores >= 80 && Scores < 100:
				fmt.Println("Excellent")
			case Scores == 100:
				fmt.Println("Perfect Score!")
				break outter
			}

			select {
			case <-time.After(2 * time.Second):
				fmt.Println("玩家操作完成")
			case <-time.After(time.Second):
				fmt.Println("玩家操作超时，继续下一分数")
			}

			if Scores > 90 {
				fmt.Println("高分出现，停止本轮")
				break
			}
		}
	}
	fmt.Println("游戏结束")
}
