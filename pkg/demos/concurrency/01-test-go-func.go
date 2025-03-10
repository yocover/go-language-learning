package main

import (
	"fmt"
	"sync"
)

// 1. 协程
// 协程是轻量级的线程，由 Go 语言的运行时（runtime）管理。
// 协程的调度由 Go 语言的调度器（scheduler）管理，调度器负责将协程分配到可用的线程上执行。
// 协程的调度是协作式的，即协程需要主动让出 CPU 时间片，否则调度器不会切换到其他协程。
// 协程的调度是抢占式的，即调度器可以随时中断正在执行的协程，将 CPU 时间片分配给其他协程。

// 协程的创建
// 协程的创建非常简单，只需要在函数前面加上 `go` 关键字即可。
func testCreateGoroutine() {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("在协程中执行")
	}()

	fmt.Println("在主协程中执行")
	wg.Wait()
}
