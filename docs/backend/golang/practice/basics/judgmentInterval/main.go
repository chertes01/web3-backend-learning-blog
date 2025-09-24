/*
以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。请你合并所有
重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。可以先对区间数组按照区间的
起始位置进行排序，然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，将当前区间与切片中最后
一个区间进行比较，如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中。
*/
package main

import (
	"fmt"
	"sort"
)

func judge(arry2D [][]int) [][]int {
	// 按起始位置排序
	sort.Slice(arry2D, func(i, j int) bool {
		return arry2D[i][0] < arry2D[j][0]
	})

	var res [][]int
	// 先放第一个区间
	res = append(res, arry2D[0])

	for i := 1; i < len(arry2D); i++ {
		last := res[len(res)-1]
		cur := arry2D[i]

		if cur[0] <= last[1] { // 有重叠
			if cur[1] > last[1] {
				res[len(res)-1][1] = cur[1] // 更新右边界
			}
		} else { // 无重叠
			res = append(res, cur)
		}
	}
	return res
}

func main() {
	intervals := [][]int{{1, 3}, {3, 5}, {4, 8}, {9, 13}}
	fmt.Println(judge(intervals)) // [[1 8] [9 13]]
}
