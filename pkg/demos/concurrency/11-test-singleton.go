package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// 1. 懒汉式（使用 sync.Once）
// 懒汉式：在需要的时候才创建实例
type Singleton struct {
	data string
}

var (
	instance *Singleton
	once     sync.Once
)

// GetInstance 获取单例实例（推荐使用这种方式）
func GetInstance() *Singleton {
	once.Do(func() {
		instance = &Singleton{data: "单例数据"}
	})
	return instance
}

// 2. 双重检查锁定模式
type Singleton2 struct {
	data string
}

var (
	instance2 *Singleton2
	mu        sync.Mutex
)

// GetInstance2 获取单例实例（双重检查锁定）
func GetInstance2() *Singleton2 {
	if instance2 == nil {
		mu.Lock()
		defer mu.Unlock()
		if instance2 == nil {
			instance2 = &Singleton2{data: "单例数据2"}
		}
	}
	return instance2
}

// 3. 原子操作方式
type Singleton3 struct {
	data string
}

var (
	instance3 atomic.Value
	initiated uint32
)

// GetInstance3 获取单例实例（原子操作）
func GetInstance3() *Singleton3 {
	if atomic.LoadUint32(&initiated) == 1 {
		return instance3.Load().(*Singleton3)
	}

	mu.Lock()
	defer mu.Unlock()

	if initiated == 0 {
		s := &Singleton3{data: "单例数据3"}
		instance3.Store(s)
		atomic.StoreUint32(&initiated, 1)
	}

	return instance3.Load().(*Singleton3)
}

// 4. 包级别私有变量（最简单但不是懒加载）
var instance4 = &Singleton{data: "单例数据4"}

func GetInstance4() *Singleton {
	return instance4
}

// 测试单例模式
func testSingleton() {
	fmt.Println("=== 测试单例模式 ===")

	// 1. 测试 sync.Once 方式
	fmt.Println("\n1. 测试 sync.Once 方式:")
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			instance := GetInstance()
			fmt.Printf("协程 %d 获取的实例地址: %p\n", id, instance)
		}(i)
	}
	wg.Wait()

	// 2. 测试双重检查锁定方式
	fmt.Println("\n2. 测试双重检查锁定方式:")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			instance := GetInstance2()
			fmt.Printf("协程 %d 获取的实例地址: %p\n", id, instance)
		}(i)
	}
	wg.Wait()

	// 3. 测试原子操作方式
	fmt.Println("\n3. 测试原子操作方式:")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			instance := GetInstance3()
			fmt.Printf("协程 %d 获取的实例地址: %p\n", id, instance)
		}(i)
	}
	wg.Wait()

	// 4. 测试包级别私有变量方式
	fmt.Println("\n4. 测试包级别私有变量方式:")
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			instance := GetInstance4()
			fmt.Printf("协程 %d 获取的实例地址: %p\n", id, instance)
		}(i)
	}
	wg.Wait()
}

/* 单例模式的几种实现方式比较：

1. sync.Once 方式（推荐）：
   - 优点：
     * 线程安全
     * 惰性初始化
     * 实现简单
     * 性能好
   - 缺点：
     * 无法传递参数

2. 双重检查锁定：
   - 优点：
     * 线程安全
     * 惰性初始化
     * 可以减少锁的使用
   - 缺点：
     * 实现相对复杂
     * 在某些情况下可能有内存可见性问题

3. 原子操作：
   - 优点：
     * 线程安全
     * 性能好
     * 可以避免锁的开销
   - 缺点：
     * 实现较复杂
     * 代码可读性较差

4. 包级别私有变量：
   - 优点：
     * 实现最简单
     * 线程安全（Go的包初始化是线程安全的）
   - 缺点：
     * 非惰性初始化
     * 无法延迟创建实例

最佳实践：
1. 如果不需要延迟初始化，使用包级别私有变量
2. 如果需要延迟初始化，使用 sync.Once
3. 如果需要传递参数，可以考虑使用双重检查锁定或原子操作
4. 在大多数情况下，sync.Once 是最好的选择
*/
