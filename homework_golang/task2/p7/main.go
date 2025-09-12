/*
编写一个程序，使用通道实现两个协程之间的通信。
一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。

考察点 ：通道的基本使用、协程间通信。
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

var wd sync.WaitGroup

func main() {
	c := make(chan int)
	wd.Add(2)
	defer wd.Wait()
	go sendData(c)
	go recvData(c)
}

func sendData(c chan<- int) {
	defer wd.Done()
	for n := 1; n <= 10; n++ {
		fmt.Println("send:", n)
		c <- n
		time.Sleep(time.Second)
	}
	close(c)
}

func recvData(c <-chan int) {
	defer wd.Done()
	for n := range c {
		fmt.Println("receive:", n)
		time.Sleep(time.Second)
	}
}
