/*
题目 ：使用原子操作（ sync/atomic 包）实现一个无锁的计数器。
启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
考察点 ：原子操作、并发数据安全。
*/
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	count int64
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
			for i := 0; i < 1000; i++ {
				atomic.AddInt64(&count, 1)
				fmt.Printf("goroutine%d: %d\n", k+1, count)
			}
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
