/*
题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。

考察点 ： sync.Mutex 的使用、并发数据安全。
*/
package main

import (
	"fmt"
	"sync"
)

var (
	count int
	mutex sync.Mutex
	wg    sync.WaitGroup
)

func main() {
	tasks := [10]func(){}
	taskInit(tasks[:])
	taskExecute(tasks[:])
}

func taskInit(tasks []func()) {
	for k, _ := range tasks {
		tasks[k] = func() {
			defer wg.Done()
			mutex.Lock()
			for i := 0; i < 1000; i++ {
				count++
				fmt.Printf("goroutine%d: %d\n", k+1, count)
			}
			mutex.Unlock()
		}
	}
}

func taskExecute(tasks []func()) {
	wg.Add(10)
	defer wg.Wait()

	for _, task := range tasks {
		go task()
	}
}
