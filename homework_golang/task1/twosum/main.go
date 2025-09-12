/*
给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
*/
package main

import "fmt"

func main() {
	// nums := []int{2, 7, 11, 15}
	// target := 9
	// fmt.Println(twoSum(nums, target))

	nums := []int{3, 2, 4}
	target := 6
	fmt.Println(twoSum(nums, target))
}

func twoSum(nums []int, target int) []int {
	mp := map[int]int{}
	for i, num := range nums {
		if j, ok := mp[target-num]; ok {
			return []int{j, i}
		}
		mp[num] = i
	}
	return nil
}
