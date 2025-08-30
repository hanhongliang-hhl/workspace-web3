package main

import (
	"fmt"
	"sync"
)



func main() {
	//题目 ：编写一个程序，使用 sync.Mutex 来保护一个共享的计数器。启动10个协程，每个协程对计数器进行1000次递增操作，最后输出计数器的值。
	var wg sync.WaitGroup
	var mutex sync.Mutex
	
	var count int = 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 1000; j++ {
				fmt.Println("协程", i, "正在执行第", j, "次操作")
				mutex.Lock()
				count++
				mutex.Unlock()
			}
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println(count)
	
}
