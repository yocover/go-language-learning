package main

import (
	"fmt"
	"sync"
	"time"
)

// 演示死锁的简单示例
func testDeadlock() {
	fmt.Println("=== 测试死锁 ===")

	// 创建两个互斥锁
	lock1 := &sync.Mutex{}
	lock2 := &sync.Mutex{}

	// 创建等待组
	var wg sync.WaitGroup
	wg.Add(2)

	// 协程1：先获取lock1，然后尝试获取lock2
	go func() {
		defer wg.Done()

		fmt.Println("协程1: 尝试获取lock1")
		lock1.Lock()
		fmt.Println("协程1: 已获取lock1")

		// 短暂等待，确保死锁发生
		time.Sleep(time.Millisecond * 100)

		fmt.Println("协程1: 尝试获取lock2") // 这里会发生死锁
		lock2.Lock()
		fmt.Println("协程1: 已获取lock2")

		// 释放锁
		lock2.Unlock()
		lock1.Unlock()
	}()

	// 协程2：先获取lock2，然后尝试获取lock1
	go func() {
		defer wg.Done()

		fmt.Println("协程2: 尝试获取lock2")
		lock2.Lock()
		fmt.Println("协程2: 已获取lock2")

		// 短暂等待，确保死锁发生
		time.Sleep(time.Millisecond * 100)

		fmt.Println("协程2: 尝试获取lock1") // 这里会发生死锁
		lock1.Lock()
		fmt.Println("协程2: 已获取lock1")

		// 释放锁
		lock1.Unlock()
		lock2.Unlock()
	}()

	// 等待协程完成（注意：由于死锁，这里会一直等待）
	fmt.Println("等待协程完成...")
	wg.Wait()
	fmt.Println("程序结束") // 这一行不会被执行到，因为发生了死锁
}
