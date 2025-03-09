package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Queue 表示一个固定大小的队列
type Queue struct {
	items    []interface{} // 队列中的元素
	maxSize  int           // 队列最大容量
	mu       sync.Mutex    // 互斥锁保护队列
	notEmpty *sync.Cond    // 队列非空条件
	notFull  *sync.Cond    // 队列未满条件
}

// NewQueue 创建一个新的队列
func NewQueue(size int) *Queue {
	q := &Queue{
		items:   make([]interface{}, 0, size),
		maxSize: size,
	}
	q.notEmpty = sync.NewCond(&q.mu)
	q.notFull = sync.NewCond(&q.mu)
	return q
}

// Put 向队列中添加元素
func (q *Queue) Put(item interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 当队列满时，等待队列有空间
	for len(q.items) == q.maxSize {
		fmt.Printf("队列已满，生产者等待...\n")
		q.notFull.Wait()
	}

	// 添加元素
	q.items = append(q.items, item)
	fmt.Printf("生产者添加: %v, 队列大小: %d\n", item, len(q.items))

	// 通知等待的消费者
	q.notEmpty.Signal()
}

// Get 从队列中获取元素
func (q *Queue) Get() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	// 当队列为空时，等待新元素
	for len(q.items) == 0 {
		fmt.Printf("队列为空，消费者等待...\n")
		q.notEmpty.Wait()
	}

	// 获取第一个元素
	item := q.items[0]
	q.items = q.items[1:]
	fmt.Printf("消费者获取: %v, 队列大小: %d\n", item, len(q.items))

	// 通知等待的生产者
	q.notFull.Signal()
	return item
}

func TestConditionVariable(t *testing.T) {
	fmt.Println("=== 测试条件变量 ===")

	// 创建一个容量为3的队列
	queue := NewQueue(3)
	var wg sync.WaitGroup

	// 启动生产者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 6; i++ {
			queue.Put(i)
			time.Sleep(time.Millisecond * 100) // 模拟生产过程
		}
	}()

	// 启动消费者
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 1; i <= 6; i++ {
			_ = queue.Get()
			time.Sleep(time.Millisecond * 200) // 模拟消费过程
		}
	}()

	// 等待所有操作完成
	wg.Wait()
	fmt.Println("所有操作已完成")
}
