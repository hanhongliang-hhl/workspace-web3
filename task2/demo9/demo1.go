package main

import (
	"fmt"
	"sync"
)

func main() {
	fmt.Println("main 开始")
	var gr sync.WaitGroup

	/* 题目 ：编写一个程序，使用通道实现两个协程之间的通信。一个协程生成从1到10的整数，并将这些整数发送到通道中，另一个协程从通道中接收这些整数并打印出来。
	考察点 ：通道的基本使用、协程间通信。*/
	/* chan1 := make(chan int)
	gr.Add(1)
	go func() {
		for i := 1; i <= 10; i++ {
			chan1 <- i
		}
		gr.Done()
	}()

	gr.Add(1)
	go func() {
		for i := 1; i <= 10; i++ {
			fmt.Println(<-chan1)
		}
		gr.Done()
	}()
	*/
	/*题目 ：实现一个带有缓冲的通道，生产者协程向通道中发送100个整数，消费者协程从通道中接收这些整数并打印。
	  考察点 ：通道的缓冲机制。 */
	chan2 := make(chan int, 10)
	gr.Add(1)
	go func() {
		for i := 1; i <= 100; i++ {
			chan2 <- i
		}
		close(chan2)
		gr.Done()
	}()

	gr.Add(1)
	go func() {
		// time.Sleep(time.Second * 10)
		for i := 1; i <= 1020; i++ {
			fmt.Println(<-chan2)
		}
		gr.Done()
	}()

	gr.Wait()
	fmt.Println("main 结束")
}
