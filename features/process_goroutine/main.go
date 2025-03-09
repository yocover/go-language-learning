package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// 演示协程的基本使用
func demonstrateGoroutine() {
	fmt.Println("\n=== 协程(Goroutine)演示 ===")
	var wg sync.WaitGroup

	// 启动多个协程
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("协程 %d 正在执行\n", id)
			time.Sleep(time.Millisecond * 500) // 模拟工作
			fmt.Printf("协程 %d 完成\n", id)
		}(i)
	}

	fmt.Println("主协程继续执行其他工作...")
	wg.Wait() // 等待所有协程完成
	fmt.Println("所有协程已完成")
}

// 演示进程的创建和管理
func demonstrateProcess() {
	fmt.Println("\n=== 进程(Process)演示 ===")

	// 获取当前进程信息
	fmt.Printf("当前进程 ID: %d\n", os.Getpid())
	fmt.Printf("父进程 ID: %d\n", os.Getppid())

	// 创建子进程
	cmd := exec.Command("ls", "-l")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("执行命令出错: %v\n", err)
		return
	}
	fmt.Printf("子进程执行 'ls -l' 的输出:\n%s", output)
}

// 展示系统资源使用情况
func showResourceUsage() {
	fmt.Println("\n=== 系统资源使用情况 ===")
	fmt.Printf("CPU 核心数: %d\n", runtime.NumCPU())
	fmt.Printf("当前运行的协程数: %d\n", runtime.NumGoroutine())
}

func main() {
	fmt.Println("进程与协程的区别和使用示例")
	fmt.Println("============================")

	// 1. 展示协程的使用
	demonstrateGoroutine()

	// 2. 展示进程的使用
	demonstrateProcess()

	// 3. 显示资源使用情况
	showResourceUsage()
}
