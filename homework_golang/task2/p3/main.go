/*
题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
考察点 ： go 关键字的使用、协程的并发执行。
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	wg.Add(2)

	go func() {
		defer wg.Done()
		for n := 1; n <= 10; n += 2 {
			fmt.Println("程序1:", n)
			time.Sleep(time.Second)
		}
	}()

	go func() {
		defer wg.Done()
		for n := 2; n <= 10; n += 2 {
			fmt.Println("程序2:", n)
			time.Sleep(time.Second)
		}
	}()

	wg.Wait()
}
