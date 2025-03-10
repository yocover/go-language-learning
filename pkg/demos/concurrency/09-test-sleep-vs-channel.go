package main

import (
	"fmt"
	"time"
)

// 演示 Sleep 和 channel 的不同阻塞方式
func testSleepVsChannel() {
	fmt.Println("=== 对比 Sleep 和 Channel 的阻塞方式 ===")

	// 1. 使用 Sleep 的方式
	fmt.Println("\n1. Sleep 方式:")
	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("Sleep方式 - 第 %d 次\n", i)
			// Sleep 不会创建 channel，而是使用运行时的计时器
			time.Sleep(time.Second)
		}
	}()

	// 2. 使用 channel 的方式
	fmt.Println("\n2. Channel 方式:")
	ticker := time.NewTicker(time.Second) // 创建一个定时器，每秒触发一次
	defer ticker.Stop()                   // 确保退出时停止定时器

	go func() {
		for i := 1; i <= 3; i++ {
			// 通过 channel 接收定时器的信号

			// 行代码的意思是 从ticker.C这个channel中接收信号
			// 如果ticker.C这个channel中没有信号 则阻塞等待
			// 如果ticker.C这个channel中有信号 则接收信号 并打印"Channel方式 - 第 %d 次\n"
			<-ticker.C // 这里确实使用了 channel
			fmt.Printf("Channel方式 - 第 %d 次\n", i)
		}
	}()

	// 3. 手动实现的 channel 方式
	fmt.Println("\n3. 手动实现的 Channel 方式:")
	done := make(chan bool)
	go func() {
		for i := 1; i <= 3; i++ {
			fmt.Printf("手动Channel方式 - 第 %d 次\n", i)
			// 创建一个一次性的定时器
			timer := time.NewTimer(time.Second)
			<-timer.C
		}
		done <- true
	}()

	// 等待第三种方式完成
	<-done
	fmt.Println("\n所有示例执行完成")
}

/* Sleep 的工作原理：
1. time.Sleep 是由 Go 运行时直接管理的
2. 它不会创建 channel，而是使用运行时的计时器系统
3. 当调用 Sleep 时，当前协程会被挂起，但不是通过 channel 实现的
4. 运行时会在指定时间后重新唤醒该协程

Channel 方式的特点：
1. time.Ticker/time.Timer 内部会创建 channel
2. 可以通过 channel 接收定时信号
3. 更灵活，可以配合 select 使用
4. 需要手动管理（关闭、清理等）
*/
