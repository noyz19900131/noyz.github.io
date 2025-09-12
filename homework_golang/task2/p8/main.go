/*
题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
考察点 ：通道的缓冲机制。
*/
package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	ch := make(chan int, 10)
	wg.Add(2)
	defer wg.Wait()

	go func(ch chan<- int) {
		defer wg.Done()
		for i := 1; i <= 6; i++ {
			ch <- i
			fmt.Println("生产者发送:", i)
		}
		close(ch)
	}(ch)

	go func(ch <-chan int) {
		defer wg.Done()
		for num := range ch {
			fmt.Println("消费者接收:", num)
		}
	}(ch)
}
