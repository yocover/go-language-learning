我来创建一个条件变量的示例，模拟一个简单的生产者-消费者模式，展示条件变量的使用：

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Queue 表示一个固定大小的队列
type Queue struct {
	items     []interface{} // 队列中的元素
	maxSize   int          // 队列最大容量
	mu        sync.Mutex   // 互斥锁保护队列
	notEmpty  *sync.Cond   // 队列非空条件
	notFull   *sync.Cond   // 队列未满条件
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

func TestConditionVariable() {
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

func main() {
	fmt.Println("=== 测试条件变量 ===")
	TestConditionVariable()
}
```
让我修复一下 main 函数重复声明的问题：

```go:features/concurrency/condition_test.go
package main

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// Queue 表示一个固定大小的队列
type Queue struct {
	items     []interface{} // 队列中的元素
	maxSize   int          // 队列最大容量
	mu        sync.Mutex   // 互斥锁保护队列
	notEmpty  *sync.Cond   // 队列非空条件
	notFull   *sync.Cond   // 队列未满条件
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
```
让我们运行这个测试：

```bash
cd features/concurrency && go test -v -run TestConditionVariable
```
这个示例展示了条件变量的核心概念和使用方法。让我解释一下关键部分：

1. **条件变量的创建和初始化**：
```go
type Queue struct {
    mu       sync.Mutex   // 互斥锁保护队列
    notEmpty *sync.Cond   // 队列非空条件
    notFull  *sync.Cond   // 队列未满条件
}

func NewQueue(size int) *Queue {
    q := &Queue{...}
    q.notEmpty = sync.NewCond(&q.mu)
    q.notFull = sync.NewCond(&q.mu)
    return q
}
```
- 条件变量总是和互斥锁一起使用
- 每个条件变量都需要一个互斥锁作为参数来创建

2. **生产者的条件等待**：
```go
func (q *Queue) Put(item interface{}) {
    q.mu.Lock()
    defer q.mu.Unlock()

    // 当队列满时，等待队列有空间
    for len(q.items) == q.maxSize {
        fmt.Printf("队列已满，生产者等待...\n")
        q.notFull.Wait()
    }
    
    // 添加元素并通知消费者
    q.items = append(q.items, item)
    q.notEmpty.Signal()
}
```
- 使用 `Wait()` 等待条件满足
- 使用 `Signal()` 通知其他等待的协程

3. **消费者的条件等待**：
```go
func (q *Queue) Get() interface{} {
    q.mu.Lock()
    defer q.mu.Unlock()

    // 当队列为空时，等待新元素
    for len(q.items) == 0 {
        fmt.Printf("队列为空，消费者等待...\n")
        q.notEmpty.Wait()
    }
    
    // 获取元素并通知生产者
    item := q.items[0]
    q.items = q.items[1:]
    q.notFull.Signal()
    return item
}
```

条件变量的重要特点：

1. **避免忙等待**：
   - 不需要持续检查条件
   - 当条件不满足时，协程会被挂起
   - 只有在条件可能满足时才会被唤醒

2. **精确通知**：
   - `Signal()` 只唤醒一个等待的协程
   - `Broadcast()` 可以唤醒所有等待的协程
   - 避免了不必要的唤醒

3. **与互斥锁配合**：
   - 条件变量必须与互斥锁一起使用
   - 在检查条件前必须持有锁
   - `Wait()` 会自动释放锁，并在返回时重新获取锁

使用场景：

1. **生产者-消费者模式**：
   - 队列满/空的条件等待
   - 资源池管理
   - 任务调度系统

2. **多线程同步**：
   - 等待特定条件发生
   - 线程协调和通信
   - 资源状态变化通知

3. **实现复杂的同步逻辑**：
   - 多个条件的组合等待
   - 状态机的状态转换
   - 复杂的任务依赖关系

注意事项：

1. 总是在循环中使用 `Wait()`，因为条件可能会发生虚假唤醒
2. 在调用 `Wait()` 前必须持有互斥锁
3. 在发送信号前也必须持有互斥锁
4. 使用 `defer` 确保锁一定会被释放

这个示例展示了条件变量在实现线程同步时的强大功能，特别是在需要基于特定条件进行等待和通知的场景中。
