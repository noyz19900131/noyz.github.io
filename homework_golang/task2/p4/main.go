/*
题目 ：设计一个任务调度器，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
考察点 ：协程原理、并发任务调度。
*/
package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	tasks := []func(){
		func() {
			fmt.Println("Task 1")
		},
		func() {
			fmt.Println("Task 2")
		},
		func() {
			fmt.Println("Task 3")
		},
		func() {
			fmt.Println("Task 4")
		},
		func() {
			fmt.Println("Task 5")
		},
	}

	TaskScheduler(tasks)
	PrintTaskTime(tasks)
	fmt.Println("All tasks completed")
}

/*	TODO: 实现任务调度器	*/
func TaskScheduler(tasks []func()) {
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	defer wg.Wait()

	for _, task := range tasks {
		go task()
	}
	for range tasks {
		time.Sleep(time.Second)
	}
}

/*	TODO: 打印每个任务的执行时间	*/
func PrintTaskTime(tasks []func()) {
	for _, task := range tasks {
		start := time.Now()
		task()
		end := time.Now()
		fmt.Println("Task time:", end.Sub(start))
	}
}
