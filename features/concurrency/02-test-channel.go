package main

import (
	"fmt"
	"time"
)

// 2. 通道
// 通道是用于协程之间通信的机制。
// 通道是类型安全的，即通道只能传递特定类型的数据。
// 通道是线程安全的，即通道的读写操作是原子性的。
// 通道是阻塞的，即通道的读写操作会阻塞当前协程，直到有数据可读或可写。
func testCreateChannel() {
	fmt.Println("\n=== 测试无缓冲通道 ===")
	testChannel()

	fmt.Println("\n=== 测试带缓冲通道 ===")
	testBufferedChannel()

	fmt.Println("\n=== 测试通道读写 ===")
	testChannelRead()

	fmt.Println("\n=== 测试通道关闭 ===")
	testChannelClose()
}

// 通道的创建和基本使用
// 展示无缓冲通道的使用
func testChannel() {
	ch := make(chan int) // 创建无缓冲通道

	// 创建一个协程用于写入数据
	go func() {
		fmt.Println("协程: 准备写入数据")
		ch <- 1
		fmt.Println("协程: 数据写入完成")
	}()

	// 主协程读取数据
	fmt.Println("主协程: 准备读取数据")
	value := <-ch
	fmt.Printf("主协程: 读取到数据: %d\n", value)
}

// 展示带缓冲通道的使用
func testBufferedChannel() {
	ch := make(chan int, 2) // 创建容量为2的带缓冲通道

	// 写入数据不会阻塞，直到缓冲区满
	ch <- 1
	ch <- 2
	fmt.Println("成功写入两个数据到带缓冲通道")

	// 读取数据
	fmt.Printf("从通道读取的第一个数据: %d\n", <-ch)
	fmt.Printf("从通道读取的第二个数据: %d\n", <-ch)
}

// 通道的读写
// 展示不同的读写模式
func testChannelRead() {
	ch := make(chan int, 3)

	// 启动写入协程
	go func() {
		for i := 1; i <= 3; i++ {
			ch <- i
			fmt.Printf("写入数据: %d\n", i)
			time.Sleep(time.Millisecond * 100)
		}
		close(ch) // 写入完成后关闭通道
	}()

	// 使用for range读取通道数据
	for num := range ch {
		fmt.Printf("读取到数据: %d\n", num)
	}
}

// 通道的关闭
// 展示如何安全地关闭通道
func testChannelClose() {
	ch := make(chan int, 3)

	// 写入一些数据
	ch <- 1
	ch <- 2
	ch <- 3

	// 关闭通道
	close(ch)

	// 从已关闭的通道读取数据
	value, ok := <-ch
	if ok {
		fmt.Printf("读取到数据: %d\n", value)
	}

	value, ok = <-ch
	if ok {
		fmt.Printf("读取到数据: %d\n", value)
	}

	// 当通道关闭且没有数据时，ok 将为 false
	value, ok = <-ch
	if !ok {
		fmt.Println("通道已关闭且没有更多数据")
	}
}
