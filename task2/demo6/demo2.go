package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var waitGroup sync.WaitGroup

	/* 题目 ：编写一个程序，使用 go 关键字启动两个协程，一个协程打印从1到10的奇数，另一个协程打印从2到10的偶数。
	考察点 ： go 关键字的使用、协程的并发执行。*/
	waitGroup.Add(1)
	go func() {
		for i := 1; i <= 10; i++ {
			if i%2 == 1 {
				fmt.Println("1-10奇数：", i)
			}
		}
		waitGroup.Done()
	}()

	waitGroup.Add(1)
	go func() {
		for i := 2; i <= 10; i++ {
			if i%2 == 0 {
				fmt.Println("2-10偶数：", i)
			}
		}
		waitGroup.Done()
	}()

	fmt.Println("-------------------------------------------------------------------------------------")
	/*题目 ：设计一个任务"调度器"，接收一组任务（可以用函数表示），并使用协程并发执行这些任务，同时统计每个任务的执行时间。
	  考察点 ：协程原理、并发任务调度。 */

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()
	for t := range ticker.C {
		starTime := time.Now()
		waitGroup.Add(1)
		go func() {
			fmt.Println("任务调度时间：", t)
			waitGroup.Done()
		}()
		endTime := time.Now()
		fmt.Println("任务执行用时：", endTime.Sub(starTime).Nanoseconds())
	} 


	waitGroup.Wait()

	fmt.Println("main 结束")
}
