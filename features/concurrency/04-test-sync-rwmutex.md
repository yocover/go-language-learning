我来创建一个读写锁的示例，模拟一个简单的缓存系统，展示多个读取者和写入者的并发访问：

```go
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

func main() {
	testReadWriteLock()
}
```
这个示例展示了读写锁（`sync.RWMutex`）的实际应用。让我解释一下关键部分：

1. **缓存结构设计**：
```go
type Cache struct {
    sync.RWMutex // 嵌入读写锁
    data         map[string]string
}
```
- 使用读写锁保护 map 数据
- 读操作使用 `RLock()`/`RUnlock()`
- 写操作使用 `Lock()`/`Unlock()`

2. **读取操作**：
```go
func (c *Cache) Read(key string) (string, bool) {
    c.RLock() // 获取读锁
    defer c.RUnlock()
    
    // 模拟读取操作的耗时
    time.Sleep(time.Millisecond * 10)
    
    value, exists := c.data[key]
    return value, exists
}
```
- 使用读锁允许多个读取者同时访问
- 读操作不会阻塞其他读操作
- 但会被写操作阻塞

3. **写入操作**：
```go
func (c *Cache) Write(key, value string) {
    c.Lock() // 获取写锁
    defer c.Unlock()
    
    // 模拟写入操作的耗时
    time.Sleep(time.Millisecond * 50)
    
    c.data[key] = value
}
```
- 使用写锁确保独占访问
- 写操作会阻塞所有其他读写操作

4. **并发访问模式**：
- 3个写入协程，每个写入3次数据
- 10个读取协程，每个读取5次数据
- 使用 ready 通道同步所有协程的启动

读写锁的优势：
1. **并发读取**：多个读取操作可以同时进行
2. **写入保护**：写入操作时会阻塞所有其他操作，确保数据一致性
3. **性能优化**：相比普通互斥锁，在读多写少的场景下性能更好

让我们运行这个示例：

```bash
cd features/concurrency && go run main.go
```
这个示例展示了读写锁的几个重要特点：

1. **读写分离**：
   - 读操作可以并发执行
   - 写操作需要独占访问
   - 写操作会阻塞所有读操作

2. **性能优化**：
   - 读多写少的场景性能好
   - 多个读取者可以同时访问
   - 写入者获得优先权

3. **使用场景**：
   - 缓存系统
   - 配置管理
   - 数据库连接池
   - 任何读多写少的并发场景

4. **注意事项**：
   - 写锁优先级高于读锁
   - 避免长时间持有写锁
   - 注意死锁风险

这个示例很好地展示了在实际应用中如何使用读写锁来保护共享资源，同时保证了较好的并发性能。
