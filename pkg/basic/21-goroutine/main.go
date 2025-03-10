package main

import (
	"fmt"
	"time"
)

// 并发

/*
携程 coroutinue
1. 协程是轻量级线程，它与线程相比，更加轻量级，占用更少的系统资源。
2. 协程的调度由GO自身的调度器运行时调度，因此，协程的切换比线程切换要快得多。
*/
func main() {
	// 开启一个协程
	// go fmt.Println("Hello, world! - 1")
	// go sayHello()
	// go func() {
	// 	fmt.Println("Hello, world! - 2")
	// }()

	// 协程是并发的，因此，主线程不会等待协程执行完毕，而是直接退出。
	// 系统创建协程需要时间，而在此之前，主协程早已执行完毕，
	// 一旦主线程退出，其它子协程自然也就退出了
	// 并且协程的执行顺序也是不确定的
	// 就无法进行预判

	fmt.Println("start")
	for i := range 10 {
		go fmt.Println(i)
		time.Sleep(time.Millisecond)
	}
	// 让主协程等一会
	time.Sleep(time.Second * 2)
	fmt.Println("end")

	// TestWaitGroup()

	// fmt.Println("--------------use defer start---------------")
	// UseDefer()
	// fmt.Println("--------------use defer end---------------")

	// processFile()
}

func sayHello() {
	fmt.Println("Say hello!")
}

func processFile() {

	fmt.Println("processFile start")

	defer fmt.Println("processFile end")
	defer fmt.Println("Process sort 03")
	defer fmt.Println("Process sort 02")
	defer fmt.Println("Process sort 01")

	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("捕获到panic: %v\n", err)
		}
	}()

	// 可能发生panic的代码...
	panic("触发panic")
	// fmt.Println("正常执行")
}
