package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// 演示主协程退出时其他协程的行为
func testMainExit() {
	fmt.Println("=== 测试主协程退出 ===")

	// 启动一个执行长时间任务的协程
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("后台任务正在执行，第 %d 次\n", i)
			time.Sleep(time.Second)
		}
		fmt.Println("后台任务完成！") // 这行代码可能不会被执行
	}()

	// 主协程只等待很短时间就退出
	fmt.Println("主协程：等待 0.5 秒后退出...")
	time.Sleep(time.Second * 2)
	fmt.Println("主协程：即将退出")

	// 注意：此时主协程退出，程序将立即结束
	// 后台协程不会完成其任务
}

// 演示等待协程完成工作 使用 sync.WaitGroup
func testMainExit2() {
	fmt.Println("=== 测试主协程退出 ===")

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for i := 1; i <= 5; i++ {
			fmt.Printf("后台任务正在执行，第 %d 次\n", i)
			time.Sleep(time.Second)
		}
		fmt.Println("后台任务完成！")
	}()

	wg.Wait()
	fmt.Println("主协程：等待协程完成工作")
}

// 演示使用 channel 等待协程完成工作
func testMainExit3() {
	fmt.Println("=== 测试主协程退出 ===")

	ch := make(chan bool) // 创建一个通道 用于协程和主协程之间的通信

	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("后台任务正在执行，第 %d 次\n", i)
			time.Sleep(time.Second)
		}

		// 任务完成后 发送信号表示任务完成
		ch <- true
	}()

	// 阻塞等待协程完成工作
	<-ch
	fmt.Println("主协程：等待协程完成工作")
}

// 演示使用 channel 等待协程完成工作 使用 select
func testMainExit4() {
	fmt.Println("=== 测试主协程退出 ===")

	ch := make(chan bool)

	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Printf("后台任务正在执行，第 %d 次\n", i)
			time.Sleep(time.Second)

			// 如果协程在3秒内完成工作 则发送信号表示任务完成
			// 如果协程在3秒内没有完成工作 则不会发送信号 主协程会一直等待
			if i == 3 {
				// 如果协程在3秒内完成工作 则发送信号表示任务完成
				// 	如果不发送信号 主协程会一直等待 进入select的default分支

				// ch <- true
			}
		}
	}()

	// 使用 select 等待协程完成工作
	// 如果协程在6秒内完成工作 则打印"主协程：等待协程完成工作"
	// 如果协程在6秒内没有完成工作 则打印"主协程：等待协程超时"
	// select 会阻塞等待 直到其中一个 case 满足条件
	// time.After 会返回一个 channel ，当时间到达时 会向该 channel 发送一个信号
	// 所以 如果协程在6秒内完成工作 则<-ch会接收到信号 打印"主协程：等待协程完成工作"
	// 如果协程在6秒内没有完成工作 则<-time.After(time.Second * 6)会接收到信号 打印"主协程：等待协程超时"
	select {
	case <-ch:
		fmt.Println("主协程：等待协程完成工作")
	case <-time.After(time.Second * 6):
		fmt.Println("主协程：等待协程超时")
	}
}

// 使用context等待协程完成工作
func testMainExit5() {
	fmt.Println("=== 测试主协程退出 ===")

	// 创建一个上下文 设置超时时间6秒
	// 如果协程在6秒内完成工作 则ctx.Done()会接收到信号 打印"主协程：等待协程完成工作"
	// 如果协程在6秒内没有完成工作 则ctx.Done()不会接收到信号 主协程会一直等待
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
	defer cancel()

	go func(ctx context.Context) {
		for i := 1; i <= 5; i++ {
			fmt.Printf("后台任务正在执行，第 %d 次\n", i)
			time.Sleep(time.Second)

			// 如果协程在3秒内完成工作 则发送信号表示任务完成
			if i == 3 {
				// cancel()
			}
		}
	}(ctx)

	// 使用 context 等待协程完成工作
	select {
	case <-ctx.Done():
		fmt.Println("主协程：等待协程完成工作")
	}
}
