package loopGame2

import (
	"fmt"
	"math/rand"
)

func LoopGame2() {
roundloop:
	for round := 1; round <= 5; round++ {
		if round == 3 {
			continue
		}

		monstor := rand.Intn(6) + 5

		for player := 1; player <= monstor; player++ {
			hp := rand.Intn(101)
			switch {
			case hp <= 50:
				fmt.Printf("战斗成功\n")
			case hp > 50 && hp < 100:
				fmt.Printf("战斗失败\n")
				continue
			case hp == 100:
				fmt.Printf("BOSS出现，游戏结束\n")
				continue roundloop
			}
		}
	}
}
