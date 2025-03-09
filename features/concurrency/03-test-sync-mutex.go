package main

import (
	"fmt"
	"sync"
	"time"
)

// 3. 互斥锁
// 互斥锁是用于保护共享资源的机制。
// 互斥锁是线程安全的，即互斥锁的加锁和解锁操作是原子性的。
// 互斥锁是阻塞的，即互斥锁的加锁操作会阻塞当前协程，直到有锁可加。

// Counter 结构体用于演示互斥锁的使用
type Counter struct {
	mu    sync.Mutex
	count int
}

// Increment 方法安全地增加计数器的值
func (c *Counter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// GetCount 方法安全地获取计数器的值
func (c *Counter) GetCount() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.count
}

func testMutex() {
	fmt.Println("=== 测试互斥锁 - 并发读写示例 ===")
	// 创建一个计数器实例
	counter := &Counter{}

	// 创建一个等待组来等待所有协程完成
	var wg sync.WaitGroup

	// 记录开始时间
	startTime := time.Now()

	// 创建一个通道用于同步所有协程的启动
	ready := make(chan struct{})

	fmt.Println("\n=== 第一阶段：准备协程 ===")

	// 启动100个写入协程
	fmt.Println("准备启动100个写入协程...")
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			// 1. 报告协程已准备就绪
			fmt.Printf("写入协程 %d 已准备就绪，等待启动信号...\n", id)

			defer wg.Done()
			// 2. 等待启动信号
			<-ready

			// 3. 执行实际工作
			time.Sleep(time.Millisecond * time.Duration(1+id%10))
			counter.Increment()
			fmt.Printf("[写入协程 %3d] 增加计数器，当前值: %d (耗时: %v)\n",
				id, counter.GetCount(), time.Since(startTime))
		}(i)
	}

	// 启动5个读取协程
	fmt.Println("\n准备启动5个读取协程...")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			// 1. 报告协程已准备就绪
			fmt.Printf("读取协程 %d 已准备就绪，等待启动信号...\n", id)

			defer wg.Done()
			// 2. 等待启动信号
			<-ready

			// 3. 执行实际工作
			for j := 0; j < 3; j++ {
				time.Sleep(time.Millisecond * 20)
				count := counter.GetCount()
				fmt.Printf("[读取协程 %d] 第 %d 次读取，当前计数: %d (耗时: %v)\n",
					id, j+1, count, time.Since(startTime))
			}
		}(i)
	}

	// 给所有协程一些时间来准备
	time.Sleep(time.Millisecond * 100)

	fmt.Println("\n=== 第二阶段：启动所有协程 ===")
	fmt.Println("所有协程已准备就绪")
	fmt.Println("关闭 ready 通道，发出启动信号...")
	startTime = time.Now() // 重置开始时间
	close(ready)           // 通过关闭通道来发出启动信号

	// 等待所有协程完成
	wg.Wait()

	// 打印最终结果和总耗时
	fmt.Printf("\n=== 第三阶段：任务完成 ===\n")
	fmt.Printf("最终计数器的值: %d\n", counter.GetCount())
	fmt.Printf("总耗时: %v\n", time.Since(startTime))

	// 演示没有使用互斥锁的问题
	time.Sleep(time.Second) // 添加间隔，使输出更清晰
	demonstrateRaceCondition()
}

// demonstrateRaceCondition 演示不使用互斥锁可能导致的问题
func demonstrateRaceCondition() {
	fmt.Println("\n=== 演示不使用互斥锁的问题 ===")

	// 共享变量
	count := 0
	var wg sync.WaitGroup
	startTime := time.Now()

	// 启动多个协程并发访问
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			// 不安全的操作：直接访问共享变量
			current := count
			time.Sleep(time.Microsecond) // 故意增加竞争条件的可能性
			count = current + 1
			fmt.Printf("[不安全操作 %d] 当前值: %d (耗时: %v)\n",
				id, count, time.Since(startTime))
		}(i)
	}

	wg.Wait()
	fmt.Printf("\n不使用互斥锁的结果: %d (预期应该是100)\n", count)
	fmt.Printf("总耗时: %v\n", time.Since(startTime))
}
