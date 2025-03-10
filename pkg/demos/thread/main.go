package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	// for i := 0; i < 10; i++ {
	// 	wg.Add(1) // 启动一个goroutine 就登记+1
	// 	go hello(i)
	// }

	// wg.Wait() // 等待所有goroutine执行完毕

	go func() {

		i := 0

		for {
			i++

			fmt.Printf("new goroutine: %d\n", i)
			time.Sleep(time.Second)
		}
	}()
	i := 0

	for {
		i++
		fmt.Printf("main goroutine: %d\n", i)
		time.Sleep(time.Second)
		if i == 2 {
			break
		}
	}

}

func hello(i int) {

	defer wg.Done()
	fmt.Println("hello", i)
}
