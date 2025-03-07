package main

import (
	"fmt"
	"sync"
)

func TestWaitGroup() {
	fmt.Println("Start TestWaitGroup")

	var wg sync.WaitGroup
	for i := range 10 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			fmt.Println(n)
		}(i)
	}

	wg.Wait() // 等待所有协程执行完毕

	fmt.Println("End TestWaitGroup")
}
