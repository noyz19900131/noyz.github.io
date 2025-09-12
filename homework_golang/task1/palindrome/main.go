/*
题目：判断一个整数是否是回文数 。回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
*/
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// for i := 0; i < 1000000000; i++ {
	// 	if isPalindrome(i) {
	// 		fmt.Println(i)
	// 	}
	// }
	fmt.Println(isPalindrome(123321))
	fmt.Println(isPalindrome(12321))
	fmt.Println(isPalindrome(189437341))
}

func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	arr := strconv.Itoa(x)
	for i := 0; i < len(arr)/2; i++ {
		if arr[i] != arr[len(arr)-1-i] {
			return false
		}
	}
	return true
}
