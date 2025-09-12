/*
题目 ：实现一个函数，接收一个整数切片的指针，将切片中的每个元素乘以2。
考察点 ：指针运算、切片操作。
*/
package main

import "fmt"

func main() {
	slice := []int{1, 2, 3}
	multiply(slice)
	fmt.Println(slice)
}

func multiply(slice []int) {
	for i := 0; i < len(slice); i++ {
		slice[i] = 2 * slice[i]
	}
}
