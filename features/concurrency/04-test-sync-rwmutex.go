package main

import (
	"fmt"
	"sync"
	"time"
)

// Cache 结构体表示一个简单的缓存系统
type Cache struct {
	sync.RWMutex // 嵌入读写锁
	data         map[string]string
}

// NewCache 创建一个新的缓存实例
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]string),
	}
}

// Read 读取缓存中的数据
func (c *Cache) Read(key string) (string, bool) {
	c.RLock() // 获取读锁
	defer c.RUnlock()

	// 模拟读取操作的耗时
	time.Sleep(time.Millisecond * 10)

	value, exists := c.data[key]
	return value, exists
}

// Write 写入数据到缓存
func (c *Cache) Write(key, value string) {
	c.Lock() // 获取写锁
	defer c.Unlock()

	// 模拟写入操作的耗时
	time.Sleep(time.Millisecond * 50)

	c.data[key] = value
}

// 测试读写锁
// 读写锁是用于保护共享资源的机制。
// 读写锁是线程安全的，即读写锁的加锁和解锁操作是原子性的。
// 读写锁是阻塞的，即读写锁的加锁操作会阻塞当前协程，直到有锁可加。
// 读写锁是公平的，即读写锁的加锁操作会按照先来先得的原则进行。
// 读写锁是可重入的，即读写锁的加锁操作可以嵌套。

func testReadWriteLock() {
	fmt.Println("=== 测试读写锁 ===")
	cache := NewCache()
	var wg sync.WaitGroup

	// 创建一个通道用于同步所有协程的启动
	ready := make(chan struct{})

	// 启动写入协程（3个写入者）
	fmt.Println("\n准备启动3个写入协程...")
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 等待启动信号
			<-ready

			// 每个写入者写入3次不同的数据
			for j := 0; j < 3; j++ {
				key := fmt.Sprintf("key-%d-%d", id, j)
				value := fmt.Sprintf("value-%d-%d", id, j)

				startTime := time.Now()
				cache.Write(key, value)
				fmt.Printf("[写入协程 %d] 写入 %s=%s (耗时: %v)\n",
					id, key, value, time.Since(startTime))
			}
		}(i)
	}

	// 启动读取协程（10个读取者）
	fmt.Println("准备启动10个读取协程...")
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 等待启动信号
			<-ready

			// 每个读取者随机读取数据
			for j := 0; j < 5; j++ {
				// 随机选择一个key来读取
				writerID := j % 3
				keyID := j % 3
				key := fmt.Sprintf("key-%d-%d", writerID, keyID)

				startTime := time.Now()
				value, exists := cache.Read(key)
				if exists {
					fmt.Printf("[读取协程 %d] 读取 %s=%s (耗时: %v)\n",
						id, key, value, time.Since(startTime))
				} else {
					fmt.Printf("[读取协程 %d] 键 %s 不存在 (耗时: %v)\n",
						id, key, time.Since(startTime))
				}

				// 短暂休眠，模拟实际的读取间隔
				time.Sleep(time.Millisecond * 20)
			}
		}(i)
	}

	// 给协程一些时间来准备
	time.Sleep(time.Millisecond * 100)
	fmt.Println("\n所有协程已准备就绪，开始执行...")

	// 记录开始时间
	startTime := time.Now()

	// 发出启动信号
	close(ready)

	// 等待所有协程完成
	wg.Wait()

	// 打印最终统计信息
	fmt.Printf("\n=== 执行完成 ===\n")
	fmt.Printf("总耗时: %v\n", time.Since(startTime))
	fmt.Printf("缓存中的数据量: %d\n", len(cache.data))

	// 打印缓存中的所有数据
	fmt.Println("\n缓存内容:")
	cache.RLock()
	for k, v := range cache.data {
		fmt.Printf("%s = %s\n", k, v)
	}
	cache.RUnlock()
}
