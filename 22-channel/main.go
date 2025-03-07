package main

import (
	"fmt"
	"sync"
	"time"
)

// channel 是管道
// 管道可以用于并发编程，可以让多个协程之间进行通信
// 同时也可以用于并发控制
// go 通过关键字 chan 来代表管道

func main() {
	// chanExample()
	// unbufferedChan()
	// hasbufferChan()
	// chanUnbufferAndBuffer()
	// chanZuseExample()
	// fmt.Println("\n=== 安全的计数器测试 ===")
	// syncMutexTest() // 互斥锁示例

	// fmt.Println("=== 不安全的计数器测试 ===")
	// unsafeCounterTest() // 不安全的计数器示例

	// TestSelect()
	// TestWatiGroup()
	// CounterWaitOut()
	UseContextTest()
}

func CounterWaitOut() {

	mianWait := sync.WaitGroup{}
	childWait := sync.WaitGroup{}
	mianWait.Add(10)

	fmt.Println("开始计数")

	for i := range 10 {
		childWait.Add(1)
		go func() {
			fmt.Println("计数中... ", i)
			time.Sleep(time.Second)
			childWait.Done()
			mianWait.Done()
		}()

		childWait.Wait()
	}

	mianWait.Wait()
	fmt.Println("计数完毕")
}

func TestWatiGroup() {

	var wait sync.WaitGroup
	// 指定子协程的数量
	wait.Add(1)

	go func() {
		fmt.Println(1)
		wait.Done()
	}()

	wait.Wait()

	fmt.Println(2)

}

type UnsafeCounter struct {
	count int
}

func (c *UnsafeCounter) Increment() {
	c.count++                    // 没有锁保护的并发访问
	time.Sleep(time.Millisecond) // 模拟耗时操作
}

func (c *UnsafeCounter) GetCount() int {
	return c.count
}

// 不安全的计数器示例
func unsafeCounterTest() {
	counter := &UnsafeCounter{}

	// 启动100个协程同时操作计数器
	for range 100 {
		go counter.Increment()
	}

	time.Sleep(time.Second) // 等待协程执行完毕
	fmt.Println("最终计数值（不安全）：", counter.GetCount())
}

type Counter struct {
	mu    sync.Mutex // 互斥锁
	count int
}

func (c *Counter) Increment() {
	c.mu.Lock()         // 加锁
	defer c.mu.Unlock() // 函数返回前加锁
	c.count++
	time.Sleep(time.Millisecond) // 模拟耗时操作
}

func (c *Counter) GetCount() int {
	c.mu.Lock() // 解锁
	defer c.mu.Unlock()
	fmt.Println("获取计数值：", c.count)
	return c.count
}

// 互斥锁
func syncMutexTest() {
	counter := &Counter{}

	// 启动多个协程同时操作计数器
	for range 100 {
		go counter.Increment()
	}

	time.Sleep(time.Second) // 等待协程执行完毕
	fmt.Println("最终计数值：", counter.GetCount())
}

func chanZuseExample() {

	ch := make(chan struct{})

	defer close(ch)

	go goFunc()

	go func() {
		fmt.Println("协程开始执行")

		// 写入数据
		ch <- struct{}{}
	}()

	// 阻塞等带读取
	<-ch

	fmt.Println("协程执行完毕")

}

func goFunc() {
	// 协程函数
	fmt.Println("协程执行 ----")
}

// 结合使用示例
// 无缓冲通道和有缓冲通道
func chanUnbufferAndBuffer() {
	chbuffer := make(chan int, 5)

	chWrite := make(chan struct{})
	chRead := make(chan struct{})

	defer func() {

		defer fmt.Println("关闭所有通道")

		close(chbuffer)
		close(chWrite)
		close(chRead)
	}()

	// 负责写的协程
	go func() {
		for i := range 10 {
			chbuffer <- i
			fmt.Println("写入数据：", i)
		}
		chWrite <- struct{}{}
	}()

	// 负责读的协程
	go func() {
		for range 10 {
			// 读取数据都需要花费1毫秒
			time.Sleep(time.Second)
			fmt.Println("读取数据：", <-chbuffer)
		}
		chRead <- struct{}{}
	}()

	fmt.Println("写入完毕---------------", <-chWrite)
	fmt.Println("读取完毕---------------", <-chRead)
}

// 有缓冲通道
// 有缓冲通道，写入数据时，会先把数据放到缓冲区里
// 只有当缓冲区容量满了 才会阻塞的等待协程来读取数据
// 同样的，读取数据时，如果缓冲区里没有数据，则会阻塞等待协程写入数据
func hasbufferChan() {
	//
	ch := make(chan int, 2)
	defer close(ch)

	// 写入数据
	ch <- 123

	value := <-ch // 读取数据
	fmt.Println(value)
}

// 无缓冲通道在发送数据时必须立刻有人接收，否则会一直阻塞
// 无缓冲管道示例
// 无缓冲通道不应该同步使用
// 应该开启一个新的协程来发送数据
func unbufferedChan() {
	ch := make(chan int)

	// 无缓冲管道，不会临时存放任何数据
	// 正应为无缓冲管道无法存放数据
	// 所以写入数据时必须立刻有其它协程来读取数据
	// 否则会阻塞等待

	// 写入数据
	defer close(ch)

	go func() {
		ch <- 42
		fmt.Println("写入数据完成")
	}()
	n := <-ch // 读取数据

	fmt.Println(n)
}

func chanExample() {

	// 1. 创建通道

	// 无缓冲通道
	ch1 := make(chan int)
	// 有缓冲通道，缓冲区大小为2
	ch2 := make(chan string, 2)

	// 2.基础发送和接收
	go func() {
		ch1 <- 42 // 发送数据到通道
	}()

	value := <-ch1 // 从通道接收数据
	fmt.Println("接收到的值：", value)

	// 3.有缓冲通道示例
	ch2 <- "hello" //不会阻塞，因为有缓冲区
	ch2 <- "world" // 不会阻塞，缓冲区还有空间

	fmt.Println(<-ch2) // 输出 "hello"
	fmt.Println(<-ch2) // 输出 "world"

	// 4.通道的关闭
	close(ch2)

	// 5.使用 range 遍历通道

	numbers := make(chan int, 5)
	go func() {
		for i := 1; i <= 5; i++ {
			numbers <- i
		}
		close(numbers)
	}()

	for num := range numbers {
		fmt.Println("收到数字：", num)
	}

	// 6. select 多路复用
	ch3 := make(chan string)
	ch4 := make(chan string)

	go func() {
		time.Sleep(1 * time.Second)
		ch3 <- "来自通道3的数据"
	}()

	go func() {
		time.Sleep(2 * time.Second)
		ch4 <- "来自通道4的数据"
	}()

	for range 2 {
		select {
		case msg := <-ch3:
			fmt.Println("收到来自通道3的数据：", msg)
		case msg := <-ch4:
			fmt.Println("收到来自通道4的数据：", msg)
		}
	}
}
