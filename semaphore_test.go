package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// 使用计数器实现一个简单的信号量
type Semaphore struct {
	permits int           // 允许的数量
	channel chan struct{} // 使用通道实现信号量
}

// 创建一个新的信号量
func NewSemaphore(permits int) *Semaphore {
	return &Semaphore{
		permits: permits,
		channel: make(chan struct{}, permits),
	}
}

// 获取许可
func (s *Semaphore) Acquire() {
	s.channel <- struct{}{}
}

// 释放许可
func (s *Semaphore) Release() {
	<-s.channel
}

func TestSimpleSemaphore(t *testing.T) {
	fmt.Println("=== 测试简单信号量 ===")

	// 创建一个最多允许3人同时进入的电梯
	elevator := NewSemaphore(3)
	var wg sync.WaitGroup

	// 模拟5个人要使用电梯
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(personID int) {
			defer wg.Done()

			fmt.Printf("人员 %d 等待进入电梯\n", personID)
			elevator.Acquire() // 进入电梯
			fmt.Printf("人员 %d 进入电梯\n", personID)

			// 模拟在电梯中停留
			time.Sleep(time.Second * 2)

			fmt.Printf("人员 %d 离开电梯\n", personID)
			elevator.Release() // 离开电梯
		}(i)

		// 短暂延迟，使输出更清晰
		time.Sleep(time.Millisecond * 100)
	}

	wg.Wait()
	fmt.Println("所有人都已使用完电梯")
}
