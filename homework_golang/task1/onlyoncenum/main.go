/*
只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。
找出那个只出现了一次的元素。可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，
例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素。
*/
package main

import (
	"fmt"
)

func main() {
	fmt.Println(OnlyOnceNum([]int{1, 2, 2, -3, 1}))
}

func OnlyOnceNum(nums []int) int {
	var res int
	mp := map[int]int{}
	for _, v := range nums {
		mp[v]++
	}
	for k, v := range mp {
		if v == 1 {
			res = k
			break
		}
	}
	return res
}
